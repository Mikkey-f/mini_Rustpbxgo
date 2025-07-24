package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"miniRustpbxgo/internal/dao"
	"miniRustpbxgo/internal/model"
	"net/http"
)

// RobotCreateReq 接收前端创建Robot的请求体
type RobotCreateReq struct {
	UserID       uint    `json:"user_id" binding:"required"` // 关联用户ID（必传）
	Name         string  `json:"name" binding:"required"`
	Speed        float32 `json:"speed" binding:"omitempty,min=0.5,max=2.0,required"` // 语音语速（可选，范围0.5-2.0）
	Volume       int     `json:"volume" binding:"omitempty,min=0,max=10,required"`   // 语音音量（可选，范围0-10）
	Speaker      string  `json:"speaker" binding:"omitempty,max=50,required"`        // 发音人（可选，最长50字符）
	Emotion      string  `json:"emotion" binding:"omitempty"`                        // 语音情感（可选，仅支持指定值）
	SystemPrompt string  `json:"system_prompt" binding:"omitempty,required"`         // 系统提示词（可选，无长度限制）
}

type RobotCreateRsp struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

func (app *App) CreateRobot(ctx *gin.Context) {
	var (
		req RobotCreateReq
	)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("CreateRobotReq error:%v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	robotRepo := dao.NewRobotRepo(app.DB)
	if _, err := robotRepo.CreateRobot(&model.Robot{
		UserID:       req.UserID,
		Name:         req.Name,
		Speed:        req.Speed,
		Volume:       req.Volume,
		Speaker:      req.Speaker,
		Emotion:      req.Emotion,
		SystemPrompt: req.SystemPrompt,
	}); err != nil {
		logrus.Errorf("CreateRobot error:%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200,
		"message": "ok",
		"data": &RobotCreateRsp{
			Name: req.Name,
		}})

}
