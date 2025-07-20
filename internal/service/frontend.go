package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"miniRustpbxgo/internal/model"
	"net/http"
)

type BackendForWeb struct {
	WebToGoConn  *websocket.Conn
	Upgrader     *websocket.Upgrader
	GoToRustConn *websocket.Conn
}

func NewBackendForWeb() *BackendForWeb {
	return &BackendForWeb{
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// 允许cross跨域
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

// HandleWebRtcSetUp 处理前端与go后端关于文本信息的传递
func (backendForWeb *BackendForWeb) HandleWebRtcSetUp(w http.ResponseWriter, r *http.Request) {
	conn, err := backendForWeb.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("websocket upgrade error: ", err)
	}
	backendForWeb.WebToGoConn = conn
	defer func() {
		if err := backendForWeb.WebToGoConn.Close(); err != nil {
			log.Fatal("Frontend to goBackend close error: ", err)
			return
		}
	}()

	go backendForWeb.GoSendMessageToRust()
	done := make(chan bool)
	log.Println("Setting up Frontend to goBackend text connection")

	<-done
}

func (backendForWeb *BackendForWeb) GoSendMessageToRust() {
	webToConn := backendForWeb.WebToGoConn
	var frontendToGoEvent struct {
		Event     string          `json:"event"`
		Sdp       string          `json:"sdp"`
		Candidate json.RawMessage `json:"candidate"`
	}
	go func() {
		for {
			msgType, msg, err := webToConn.ReadMessage()
			if err != nil {
				log.Print("WebToGoConn ReadMessage error: ", err)
				continue
			}
			if msgType == websocket.BinaryMessage {
				// 音频信息
			} else {
				if err := json.Unmarshal(msg, &frontendToGoEvent); err != nil {
					log.Fatal("WebToGoConn Unmarshal error: ", err)
				}
				if frontendToGoEvent.Event == "candidate" && frontendToGoEvent.Candidate != nil {
					backendForWeb.SolveCandidate(frontendToGoEvent.Candidate)
				} else if frontendToGoEvent.Event == "offer" && frontendToGoEvent.Sdp != "" {
					backendForWeb.SolveOffer(frontendToGoEvent.Sdp)
				}
			}
		}
	}()
}

func (backendForWeb *BackendForWeb) SolveCandidate(rawMessage json.RawMessage) {
	log.Printf("Received ICE candidate: %s", string(rawMessage))
	goToRustConn := backendForWeb.GoToRustConn

	var candidate struct {
		Candidate     string `json:"candidate"`
		SdpMid        string `json:"sdpMid"`
		SdpMLineIndex int    `json:"sdpMLineIndex"`
	}

	if err := json.Unmarshal(rawMessage, &candidate); err != nil {
		log.Println("parse candidate failed:", err)
		return
	}

	candidateCmd := model.CandidateCommand{
		Command:    "candidate",
		Candidates: []string{candidate.Candidate},
	}

	cmdBytes, err := json.Marshal(candidateCmd)
	if err != nil {
		log.Println("marshal candidate command failed:", err)
		return
	}
	if err := goToRustConn.WriteMessage(websocket.TextMessage, cmdBytes); err != nil {
		log.Println("forward candidate command to rust backend err:", err)
	}
}

func (backendForWeb *BackendForWeb) SolveOffer(sdp string) {
	log.Printf("Received ICE offer: %s", sdp)
	goToRustConn := backendForWeb.GoToRustConn
	inviteCmd := model.InviteCommand{
		Command: "invite",
		Option: model.CallOption{
			Offer:  sdp,
			Caller: "frontend",
			Callee: "rust",
			//ASR:    s.asrOption,
			//TTS:    s.ttsOption,
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
		log.Println("goBackend to rustBackend not connected")
		return
	}
	marshal, err := json.Marshal(event)
	if err != nil {
		logrus.Println("ForwardToWebConn json.Marshal error", err)
		return
	}
	if err = conn.WriteMessage(websocket.TextMessage, marshal); err != nil {
		logrus.Println("ForwardToWebConn conn.WriteMessage error", err)
		return
	}
}
