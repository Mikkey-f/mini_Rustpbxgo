package service

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type BackendForRust struct {
	GoForRustConn *websocket.Conn
	EndPoint      string
}

// CreateBackendForRust 创建go到rust的后端管理者
func CreateBackendForRust(endPoint string) *BackendForRust {
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
		logrus.Fatal(err)
	}
	logrus.Info("Successfully connected")
	backendForRust.GoForRustConn = conn
	return conn
}
