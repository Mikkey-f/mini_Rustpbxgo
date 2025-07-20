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
	api.Routers(gin.Default(), backendForWeb)
	bools := make(chan bool)
	<-bools
}
