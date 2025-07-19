package main

import "miniRustpbxgo/internal/service"

type AppConfig struct {
}

func main() {
	endPoint := "ws://175.27.250.177:8080"
	callType := "webrtc"
	backendForRust := service.CreateBackendForRust(endPoint)
	backendForRust.Connect(callType)
}
