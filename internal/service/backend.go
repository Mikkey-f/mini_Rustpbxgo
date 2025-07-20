package service

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type BackendForRust struct {
	GoFToRustConn *websocket.Conn
	EndPoint      string
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
	backendForRust.GoFToRustConn = conn
	return conn
}
