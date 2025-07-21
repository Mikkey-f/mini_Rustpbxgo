package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
)

type BackendForRust struct {
	GoToRustConn *websocket.Conn
	EndPoint     string
}

type Event struct {
	Event     string `json:"event"`
	Text      string `json:"text"`
	Sdp       string `json:"sdp"`
	Reason    string `json:"reason"`
	Initiator string `json:"initiator"`
	Error     string `json:"error"`
}

// NewBackendForRust 创建go到rust的后端管理者
func NewBackendForRust(endPoint string) *BackendForRust {
	return &BackendForRust{
		EndPoint: endPoint,
	}
}

// Connect 创建go到rust的ws连接
func (backendForRust *BackendForRust) Connect(callType string, backendForWeb *BackendForWeb) {
	for {
		if backendForRust.GoToRustConn == nil {
			url := backendForRust.EndPoint
			url += "/call/" + callType

			conn, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				logrus.Error("goBackend connect rustBackend error", err)
				continue
			}
			logrus.Info("goBackend to rustBackend successfully connected")
			backendForRust.GoToRustConn = conn
			backendForWeb.GoToRustConn = conn
			// 监听机制
			go backendForRust.ListenGoToRustWs(backendForWeb)
		}
	}
}

// ListenGoToRustWs 监听go与rust的ws连接信息
func (backendForRust *BackendForRust) ListenGoToRustWs(backendForWeb *BackendForWeb) {
	defer func() {
		err := backendForRust.GoToRustConn.Close()
		if err != nil {
			logrus.Error("goBackend connect rustBackend closed error", err)
			return
		}
		logrus.Info("goBackend to rustBackend closed connection")
		backendForRust.GoToRustConn = nil
		backendForWeb.GoToRustConn = nil
	}()
	for {
		conn := backendForRust.GoToRustConn
		if conn == nil {
			log.Println("goBackend to rustBackend not connected")
			continue
		}
		msgType, msg, err := conn.ReadMessage()
		logrus.Debug("goBackend received from rustBackend:", msgType, msg)
		if err != nil {
			log.Println("Listen GoToRust Conn Message err", err)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return
			}
			logrus.Println("Unexpected read error:", err)
			return
		}

		if msgType != websocket.TextMessage {
			log.Println("Received non-text message: ", msgType)
			continue
		}
		log.Printf("Received from rust backend (type %d): %s", msgType, string(msg))

		var event Event

		if err = json.Unmarshal(msg, &event); err != nil {
			logrus.Println("ListenGoToRustWs json.Unmarshal error", err)
			continue
		}

		switch event.Event {
		case "asrFinal":
			log.Println("Received asrFinal message: ", event)
			backendForWeb.SolveAsrFinalEvent(&event)
		case "asrDelta":
			log.Println("Received asrDelta message: ", event)
		case "error":
			logrus.Println("Received an error message: ", event)
		case "close":
			logrus.Println("Received close message: ", event)
		case "hangup":
			logrus.Println("Received hangup message: ", event)
		case "speaking":
			logrus.Println("Received speaking message: ", event)
		case "silence":
			logrus.Println("Received silence message: ", event)
		case "trackStart":
			logrus.Println("Received trackStart message: ", event)
		case "trackEnd":
			logrus.Println("Received trackStop message: ", event)
		}

		backendForWeb.ForwardToWebConn(&event)
	}
}
