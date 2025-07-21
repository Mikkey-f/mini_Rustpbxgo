package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
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
		if backendForRust.GoToRustConn != nil {
			err := backendForRust.GoToRustConn.Close()
			if err != nil {
				logrus.Error("goBackend connect rustBackend closed error", err)
				return
			}
			logrus.Info("goBackend to rustBackend closed connection")
			backendForRust.GoToRustConn = nil
			backendForWeb.GoToRustConn = nil
		}
	}()
	for {
		conn := backendForRust.GoToRustConn
		if conn == nil {
			logrus.Error("goBackend to rustBackend not connected")
			continue
		}
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			logrus.Error("Listen GoToRust Conn Message err", err)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return
			}
			logrus.Error("Unexpected read error:", err)
			return
		}

		if msgType != websocket.TextMessage {
			logrus.Error("Received non-text message: ", msgType)
			continue
		}
		logrus.Info("Received from rust backend (type %d): %s", msgType, string(msg))

		var event Event

		if err = json.Unmarshal(msg, &event); err != nil {
			logrus.Error("ListenGoToRustWs json.Unmarshal error", err)
			continue
		}

		switch event.Event {
		case "asrFinal":
			logrus.Info("Received asrFinal message: ", event)
			backendForWeb.SolveAsrFinalEvent(&event)
		case "asrDelta":
			logrus.Info("Received asrDelta message: ", event)
		case "error":
			logrus.Error("Received an error message: ", event)
		case "close":
			logrus.Info("Received close message: ", event)
		case "hangup":
			logrus.Info("Received hangup message: ", event)
		case "speaking":
			logrus.Info("Received speaking message: ", event)
		case "silence":
			logrus.Info("Received silence message: ", event)
		case "trackStart":
			logrus.Info("Received trackStart message: ", event)
		case "trackEnd":
			logrus.Info("Received trackStop message: ", event)
		}

		backendForWeb.ForwardToWebConn(&event)
	}
}
