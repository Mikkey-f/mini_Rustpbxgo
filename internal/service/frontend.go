package service

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"log"
	"miniRustpbxgo/internal/dao"
	"miniRustpbxgo/internal/handler"
	"miniRustpbxgo/internal/model"
	"net/http"
)

type BackendForWeb struct {
	WebToGoConn  *websocket.Conn
	Upgrader     *websocket.Upgrader
	GoToRustConn *websocket.Conn
	AsrOption    *model.ASROption
	TtsOption    *model.TTSOption
	LLMHandler   *handler.LLMHandler
	DB           *gorm.DB
	Model        string
}

type TtsCommand struct {
	Command     string           `json:"command"`
	Text        string           `json:"text"`
	Speaker     string           `json:"speaker,omitempty"`
	PlayID      string           `json:"playId,omitempty"`
	AutoHangup  bool             `json:"autoHangup,omitempty"`
	Streaming   bool             `json:"streaming,omitempty"`
	EndOfStream bool             `json:"endOfStream,omitempty"`
	Option      *model.TTSOption `json:"option,omitempty"`
}

type WebRTCSetUpReq struct {
	ApiKey    string `json:"api_key" binding:"required"`
	ApiSecret string `json:"api_secret" binding:"required"`
	RobotId   int64  `json:"robot_id" binding:"required"`
}

func NewBackendForWeb(asrOption *model.ASROption, ttsOption *model.TTSOption, llmHandler *handler.LLMHandler, model string) *BackendForWeb {
	return &BackendForWeb{
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// 允许cross跨域
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		TtsOption:  ttsOption,
		AsrOption:  asrOption,
		LLMHandler: llmHandler,
		Model:      model,
	}
}

func NewBackendForWebByNoParam(db *gorm.DB) *BackendForWeb {
	return &BackendForWeb{
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// 允许cross跨域
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		DB: db,
	}
}

// HandleWebRtcSetUp 处理前端与go后端关于文本信息的传递
func (backendForWeb *BackendForWeb) HandleWebRtcSetUp(w http.ResponseWriter, r *http.Request, backendForRust *BackendForRust, ctx *gin.Context) {
	//要做业务上面的逻辑检验
	if backendForWeb.AsrOption == nil || backendForWeb.TtsOption == nil || backendForWeb.LLMHandler == nil {
		logrus.Error("not init backendForWeb")
		return
	}
	if backendForWeb.WebToGoConn != nil {
		if err := backendForWeb.WebToGoConn.Close(); err != nil {
			logrus.Error("frontend to goBackend closed error:", err)
			return
		}
		backendForWeb.WebToGoConn = nil // 主动置空旧连接
	}
	conn, err := backendForWeb.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error("websocket upgrade error: ", err)
		return
	}

	// 初始化backendForWeb内容
	backendForWeb.WebToGoConn = conn

	done := make(chan bool)
	go backendForWeb.GoSendMessageToRust(done)
	logrus.Info("Setting up Frontend to goBackend connection")
	<-done
	if backendForRust.GoToRustConn != nil {
		if err := backendForRust.GoToRustConn.Close(); err != nil {
			logrus.Error("goBackend connect rustBackend closed error", err)
			return
		}
		logrus.Info("goBackend to rustBackend closed connection")
		backendForRust.GoToRustConn = nil
		backendForWeb.GoToRustConn = nil
	}
}

func (backendForWeb *BackendForWeb) GoSendMessageToRust(done chan bool) {
	go func() {
		defer func() {
			close(done)
		}()
		webToConn := backendForWeb.WebToGoConn
		var frontendToGoEvent struct {
			Event     string          `json:"event"`
			Sdp       string          `json:"sdp"`
			Candidate json.RawMessage `json:"candidate"`
			Reason    string          `json:"reason"`
			Initiator string          `json:"initiator"`
		}
		for {
			_, msg, err := webToConn.ReadMessage()
			if err != nil {
				logrus.Error("webToConn.ReadMessage error: ", err)
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					return
				}
				continue
			}
			if err := json.Unmarshal(msg, &frontendToGoEvent); err != nil {
				log.Fatal("WebToGoConn Unmarshal error: ", err)
			}
			if frontendToGoEvent.Event == "candidate" && frontendToGoEvent.Candidate != nil {
				backendForWeb.SolveCandidate(frontendToGoEvent.Candidate)
			} else if frontendToGoEvent.Event == "offer" && frontendToGoEvent.Sdp != "" {
				backendForWeb.SolveOffer(frontendToGoEvent.Sdp)
			} else if frontendToGoEvent.Event == "hangup" {
				backendForWeb.SolveHangup(frontendToGoEvent.Reason)
			}
		}
	}()
}

func (backendForWeb *BackendForWeb) SolveCandidate(rawMessage json.RawMessage) {
	if backendForWeb.GoToRustConn == nil {
		logrus.Error("Backend to goBackend not ready")
		return
	}
	logrus.Info("Received ICE candidate: %s", string(rawMessage))

	goToRustConn := backendForWeb.GoToRustConn

	var candidate struct {
		Candidate     string `json:"candidate"`
		SdpMid        string `json:"sdpMid"`
		SdpMLineIndex int    `json:"sdpMLineIndex"`
	}

	if err := json.Unmarshal(rawMessage, &candidate); err != nil {
		logrus.Error("parse candidate failed:", err)
		return
	}

	candidateCmd := model.CandidateCommand{
		Command:    "candidate",
		Candidates: []string{candidate.Candidate},
	}

	cmdBytes, err := json.Marshal(candidateCmd)
	if err != nil {
		logrus.Info("marshal candidate command failed:", err)
		return
	}
	if err := goToRustConn.WriteMessage(websocket.TextMessage, cmdBytes); err != nil {
		logrus.Info("forward candidate command to rust backend err:", err)
	}
}

func (backendForWeb *BackendForWeb) SolveHangup(reason string) {
	hangupCommand := model.HangupCommand{
		Command: "hangup",
		Reason:  reason,
	}
	logrus.Println("hangup command:", hangupCommand)
	if err := backendForWeb.GoToRustConn.WriteJSON(hangupCommand); err != nil {
		log.Println("forward hangup command to rust backend err:", err)
		return
	}
}

func (backendForWeb *BackendForWeb) SolveOffer(sdp string) {
	if backendForWeb.GoToRustConn == nil {
		logrus.Println("Backend to goBackend not ready")
		return
	}
	log.Printf("Received ICE offer: %s", sdp)
	goToRustConn := backendForWeb.GoToRustConn
	inviteCmd := model.InviteCommand{
		Command: "invite",
		Option: model.CallOption{
			Offer:  sdp,
			Caller: "frontend",
			Callee: "rust",
			ASR:    backendForWeb.AsrOption,
			TTS:    backendForWeb.TtsOption,
		},
	}
	cmdBytes, err := json.Marshal(inviteCmd)
	if err != nil {
		log.Println("marshal invite command failed:", err)
		return
	}
	if err := goToRustConn.WriteMessage(websocket.TextMessage, cmdBytes); err != nil {
		log.Println("forward invite command to rust backend err:", err)
	}
}

func (backendForWeb *BackendForWeb) ForwardToWebConn(event *Event) {
	conn := backendForWeb.WebToGoConn
	if conn == nil {
		logrus.Error("goBackend to rustBackend not connected")
		return
	}
	marshal, err := json.Marshal(event)
	if err != nil {
		logrus.Error("ForwardToWebConn json.Marshal error", err)
		return
	}
	if err = conn.WriteMessage(websocket.TextMessage, marshal); err != nil {
		logrus.Error("ForwardToWebConn conn.WriteMessage error", err)
		return
	}
}

func (backendForWeb *BackendForWeb) SolveAsrFinalEvent(event *Event) {
	if event.Text == "" {
		return
	}
	var rep Event
	response, err := backendForWeb.LLMHandler.QueryStream(backendForWeb.Model, event.Text, func(segment string, playID string, autoHangup bool) error {
		if len(segment) == 0 {
			return nil
		}
		logrus.WithFields(logrus.Fields{
			"segment":    segment,
			"playID":     playID,
			"autoHangup": autoHangup,
		}).Info("Sending TTS segment")
		return backendForWeb.SendTTSCommandForRustBackend(segment, playID, autoHangup, nil)
	})
	if err != nil {
		logrus.Error("SolveAsrFinalEvent response error:", err)
		return
	}
	rep.Text = response
	rep.Event = "LLMResult"
	repStr, err := json.Marshal(rep)
	if err != nil {
		logrus.Println("SolveAsrFinalEvent json.Marshal error:", err)
		return
	}
	if err := backendForWeb.WebToGoConn.WriteMessage(websocket.TextMessage, repStr); err != nil {
		logrus.Println("SolveAsrFinalEvent response the LLS Message error: ", err)
	}
}

func (backendForWeb *BackendForWeb) SendTTSCommandForRustBackend(text string, playId string, autoHangup bool, option *model.TTSOption) error {
	ttsCommand := &TtsCommand{
		Command:     "tts",
		Text:        text,
		Speaker:     "",
		PlayID:      playId,
		AutoHangup:  autoHangup,
		Streaming:   false,
		EndOfStream: true,
		Option:      option,
	}
	logrus.Println("send ttsCommand to rust backend", ttsCommand)
	return backendForWeb.GoToRustConn.WriteJSON(ttsCommand)
}

func (backendForWeb *BackendForWeb) FrontendInit(ctx *gin.Context) {
	var webRTCSetUpReq WebRTCSetUpReq
	if err := ctx.ShouldBindJSON(&webRTCSetUpReq); err != nil {
		logrus.Errorf("WebRTCSetUpReq error:%v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo := dao.NewRobotKeyRepo(backendForWeb.DB)
	key, err := repo.GetRobotKeyByAPIKey(webRTCSetUpReq.ApiKey)
	if err != nil {
		logrus.Errorf("FrontendInit GetRobotKeyByAPIKey error:%v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if key == nil {
		logrus.Errorf("FrontendInit GetRobotKeyByAPIKey error:key not found")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "key not found"})
		return
	}
	robotRepo := dao.NewRobotRepo(backendForWeb.DB)
	robot, err := robotRepo.GetRobotByID(uint(webRTCSetUpReq.RobotId))
	if err != nil {
		logrus.Errorf("FrontendInit GetRobotByID error:%v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if robot == nil {
		logrus.Errorf("FrontendInit GetRobotByID error:robot not found")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "robot not found"})
		return
	}
	asrOption := &model.ASROption{
		Provider:  "tencent",
		AppID:     key.ASRAppID,
		SecretID:  key.ASRSecretID,
		SecretKey: key.ASRSecretKey,
		ModelType: key.ASRLanguage,
		Endpoint:  key.ASRProvider,
	}
	ttsOption := &model.TTSOption{
		Provider:  "tencent",
		Speaker:   robot.Speaker,
		AppID:     key.TTSAppID,
		SecretID:  key.TTSSecretID,
		SecretKey: key.TTSSecretKey,
		Endpoint:  key.TTSProvider,
		Speed:     robot.Speed,
		Volume:    int32(robot.Volume),
		Emotion:   robot.Emotion,
	}
	logger := logrus.New()
	c := context.Background()
	llmHandler := handler.NewLLMHandler(c, key.LLMApiKey, key.LLMApiUrl, robot.SystemPrompt, logger)
	backendForWeb.LLMHandler = llmHandler
	backendForWeb.AsrOption = asrOption
	backendForWeb.TtsOption = ttsOption
	backendForWeb.Model = "qwen-turbo"
	ctx.JSON(http.StatusOK, gin.H{"code": 200,
		"message": "初始化成功",
	})
}
