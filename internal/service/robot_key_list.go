package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"miniRustpbxgo/internal/dao"
	"miniRustpbxgo/internal/model"
	"net/http"
)

type RobotKeyListReq struct {
	UserID uint `json:"user_id" binding:"required"` // 关联用户ID（必填）
}

type RobotKeyListRsp struct {
	RobotKeyList []RobotKeyCreateRsp `json:"robot_key_list"`
	Count        int64               `json:"count"`
}

func (app *App) RobotKeyList(c *gin.Context) {
	var (
		req                RobotKeyListReq
		robotKeyCreateList []RobotKeyCreateRsp
		robotKeyList       []model.RobotKey
		count              int64
	)
	if err := c.BindJSON(&req); err != nil {
		logrus.Error("RobotKeyList bind json failed: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	robotKeyRepo := dao.NewRobotKeyRepo(app.DB)
	robotKeyList, count, err := robotKeyRepo.ListRobotKeysByUserID(req.UserID, 1, 10)
	if err != nil {
		logrus.Error("ListRobotKeysByUserID failed: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i := range robotKeyList {
		robotKeyCreateList = append(robotKeyCreateList, RobotKeyCreateRsp{
			Name:      robotKeyList[i].Name,
			APIKey:    robotKeyList[i].APIKey,
			APISecret: robotKeyList[i].APISecret,
		})
	}
	c.JSON(http.StatusOK, gin.H{"code": 200,
		"message": "ok",
		"data": &RobotKeyListRsp{
			Count:        count,
			RobotKeyList: robotKeyCreateList,
		}})
}
