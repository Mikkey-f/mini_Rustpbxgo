package main

import (
	"github.com/gin-gonic/gin"
	"miniRustpbxgo/internal/api"
	"miniRustpbxgo/internal/service"
)

type AppConfig struct {
}

func main() {
	endPoint := "ws://175.27.250.177:8080"
	callType := "webrtc"
	backendForRust := service.NewBackendForRust(endPoint)
	backendForRust.Connect(callType)
	backendForWeb := service.NewBackendForWeb()
	backendForWeb.GoToRustConn = backendForRust.GoToRustConn

	// backendForWeb和backendForRust两者要做好初始化
	api.Routers(gin.Default(), backendForWeb, backendForRust)
	select {}
}
