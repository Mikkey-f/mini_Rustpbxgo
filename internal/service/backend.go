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
	Event string `json:"event"`
	Text  string `json:"text"`
	Sdp   string `json:"sdp"`
}

// NewBackendForRust 创建go到rust的后端管理者
func NewBackendForRust(endPoint string) *BackendForRust {
	return &BackendForRust{
		EndPoint: endPoint,
	}
}

// Connect 创建go到rust的ws连接
func (backendForRust *BackendForRust) Connect(callType string) *websocket.Conn {
	url := backendForRust.EndPoint
	url += "/call/" + callType

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logrus.Fatal("goBackend connect rustBackend error", err)
	}
	logrus.Info("Successfully connected")
	backendForRust.GoToRustConn = conn
	return conn
}

// ListenGoToRustWs 监听go与rust的ws连接信息
func (backendForRust *BackendForRust) ListenGoToRustWs(backendForWeb *BackendForWeb) {
	conn := backendForRust.GoToRustConn
	for {
		if conn == nil {
			log.Println("goBackend to rustBackend not connected")
			return
		}
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Listen GoToRust Conn Message err", err)
			continue
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
		}

		backendForWeb.ForwardToWebConn(&event)
	}

}
