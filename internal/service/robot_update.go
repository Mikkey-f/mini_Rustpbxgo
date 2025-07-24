package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"miniRustpbxgo/internal/dao"
	"miniRustpbxgo/internal/model"
	"net/http"
	"time"
)

type RobotUpdateReq struct {
	Id           uint    `json:"id" binding:"required"`
	UserID       uint    `json:"user_id" binding:"required"`
	Name         string  `json:"name" binding:"omitempty"`
	Speed        float32 `json:"speed" binding:"omitempty,min=0.5,max=2.0"` // 语音语速（可选，范围0.5-2.0）
	Volume       int     `json:"volume" binding:"omitempty,min=0,max=10"`   // 语音音量（可选，范围0-10）
	Speaker      string  `json:"speaker" binding:"omitempty,max=50"`        // 发音人（可选，最长50字符）
	Emotion      string  `json:"emotion" binding:"omitempty"`               // 语音情感（可选，仅支持指定值）
	SystemPrompt string  `json:"system_prompt" binding:"omitempty"`         // 系统提示词（可选）
}

func (app *App) UpdateRobot(ctx *gin.Context) {
	var (
		req RobotUpdateReq
	)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("UpdateRobot error:%v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	robotRepo := dao.NewRobotRepo(app.DB)
	if err := robotRepo.UpdateRobot(&model.Robot{
		ID:           req.Id,
		UserID:       req.UserID,
		Name:         req.Name,
		Speed:        req.Speed,
		Volume:       req.Volume,
		Speaker:      req.Speaker,
		Emotion:      req.Emotion,
		SystemPrompt: req.SystemPrompt,
		CreatedAt:    time.Now(),
	}); err != nil {
		logrus.Errorf("UpdateRobot error:%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200,
		"message": "ok",
	})

}
