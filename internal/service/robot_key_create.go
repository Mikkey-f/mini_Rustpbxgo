package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"miniRustpbxgo/internal/dao"
	"miniRustpbxgo/internal/model"
	"miniRustpbxgo/internal/utils"
	"net/http"
)

type RobotKeyCreateReq struct {
	UserID       uint   `json:"user_id" binding:"required"`                          // 关联用户ID（必填）
	Name         string `json:"name" binding:"omitempty,max=100"`                    // 密钥名称（可选，最长100字符）
	LLMProvider  string `json:"llm_provider" binding:"omitempty,max=100,required"`   // 大模型提供商（可选，最长100字符）
	LLMApiKey    string `json:"llm_api_key" binding:"omitempty,max=255,required"`    // 大模型API密钥（可选，最长255字符）
	LLMApiUrl    string `json:"llm_api_url" binding:"omitempty,max=255,required"`    // 大模型API地址（可选，最长255字符）
	ASRProvider  string `json:"asr_provider" binding:"omitempty,max=100,required"`   // 语音识别提供商（可选，最长100字符）
	ASRAppID     string `json:"asr_app_id" binding:"omitempty,max=100,required"`     // 语音识别AppID（可选，最长100字符）
	ASRSecretID  string `json:"asr_secret_id" binding:"omitempty,max=255,required"`  // 语音识别SecretID（可选，最长255字符）
	ASRSecretKey string `json:"asr_secret_key" binding:"omitempty,max=255,required"` // 语音识别SecretKey（可选，最长255字符）
	ASRLanguage  string `json:"asr_language" binding:"omitempty,required"`           // 语音识别语言（可选，仅支持指定值）
	TTProvider   string `json:"tts_provider" binding:"omitempty,max=100,required"`   // 语音合成提供商（可选，最长100字符）
	TTSAppID     string `json:"tts_app_id" binding:"omitempty,max=100,required"`     // 语音合成AppID（可选，最长100字符）
	TTSSecretID  string `json:"tts_secret_id" binding:"omitempty,max=255,required"`  // 语音合成SecretID（可选，最长255字符）
	TTSSecretKey string `json:"tts_secret_key" binding:"omitempty,max=255,required"` // 语音合成SecretKey（可选，最长255字符）
}

type RobotKeyCreateRsp struct {
	Name      string `json:"name" binding:"omitempty,max=100"`
	APIKey    string `json:"api_key" binding:"omitempty,max=255"`
	APISecret string `json:"api_secret" binding:"omitempty,max=255"`
}

// CreateRobotKey 生成随机字符串
func (app *App) CreateRobotKey(ctx *gin.Context) {
	var (
		req            RobotKeyCreateReq
		robotApiKey    string
		robotApiSecret string
	)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("RobotKeyCreateReq error:%v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	robotKeyRepo := dao.NewRobotKeyRepo(app.DB)
	robotApiKey, err := utils.GenerateSecureRandomString(25)
	if err != nil {
		logrus.Errorf("GenerateSecureRandomString error:%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	robotApiSecret, err = utils.GenerateSecureRandomString(25)
	if err != nil {
		logrus.Errorf("GenerateSecureRandomString error:%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if _, err := robotKeyRepo.CreateRobotKey(&model.RobotKey{
		UserID:       req.UserID,
		Name:         req.Name,
		LLMProvider:  req.LLMProvider,
		LLMApiKey:    req.LLMApiKey,
		LLMApiUrl:    req.LLMApiUrl,
		ASRProvider:  req.ASRProvider,
		ASRAppID:     req.ASRAppID,
		ASRSecretID:  req.ASRSecretID,
		ASRSecretKey: req.ASRSecretKey,
		ASRLanguage:  req.ASRLanguage,
		TTSProvider:  req.TTProvider,
		TTSAppID:     req.TTSAppID,
		TTSSecretID:  req.TTSSecretID,
		TTSSecretKey: req.TTSSecretKey,
		APIKey:       robotApiKey,
		APISecret:    robotApiSecret,
	}); err != nil {
		logrus.Errorf("CreateRobotKey error:%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200,
		"message": "ok",
		"data": &RobotKeyCreateRsp{
			Name:      req.Name,
			APIKey:    robotApiKey,
			APISecret: robotApiSecret,
		}})
}
