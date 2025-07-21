package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"miniRustpbxgo/internal/api"
	"miniRustpbxgo/internal/handler"
	"miniRustpbxgo/internal/model"
	"miniRustpbxgo/internal/service"
)

type AppConfig struct {
}

func main() {
	endPoint := "ws://175.27.250.177:8080"
	callType := "webrtc"

	var openaiKey string = "yours"
	var openaiModel string = "qwen-turbo"
	var openaiEndpoint string = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	var openaiSystemPrompt string = "You are a helpful assistant. Provide concise responses. Use 'hangup' tool when the conversation is complete."

	var asrEndpoint string = "asr.tencentcloudapi.com'"
	var asrAppID string = "yours"
	var asrSecretID string = "yours"
	var asrSecretKey string = "yours"
	var asrModelType string = "16k_zh"
	var asrProvider string = "tencent"

	var ttsProvider string = "tencent"
	var ttsEndpoint string = "tts.tencentcloudapi.com"
	var ttsAppID string = "yours"
	var ttsSecretID string = "yours"
	var ttsSecretKey string = "yours"
	var ttsSpeaker string = "601003"

	asrOption := &model.ASROption{
		Provider:  asrProvider,
		AppID:     asrAppID,
		SecretID:  asrSecretID,
		SecretKey: asrSecretKey,
		Endpoint:  asrEndpoint,
		ModelType: asrModelType,
	}

	ttsOption := &model.TTSOption{
		Provider:  ttsProvider,
		Speaker:   ttsSpeaker,
		AppID:     ttsAppID,
		SecretID:  ttsSecretID,
		SecretKey: ttsSecretKey,
		Endpoint:  ttsEndpoint,
	}

	logger := logrus.New()
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	llmHandler := handler.NewLLMHandler(ctx, openaiKey, openaiEndpoint, openaiSystemPrompt, logger)
	backendForRust := service.NewBackendForRust(endPoint)
	backendForWeb := service.NewBackendForWeb(asrOption, ttsOption, llmHandler, openaiModel)
	go backendForRust.Connect(callType, backendForWeb)
	// backendForWeb和backendForRust两者要做好初始化
	api.Routers(gin.Default(), backendForWeb, backendForRust)

	select {}
}
