package service

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type BackendForWeb struct {
	WebToGoConn *websocket.Conn
	Upgrader    *websocket.Upgrader
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

// HandleSetUpText 处理前端与go后端关于文本信息的传递
func (backendForWeb *BackendForWeb) HandleSetUpText(w http.ResponseWriter, r *http.Request) {
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
	log.Println("Setting up Frontend to goBackend text connection")

}
