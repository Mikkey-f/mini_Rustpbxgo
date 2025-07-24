package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"miniRustpbxgo/internal/dao"
	"miniRustpbxgo/internal/model"
	"net/http"
)

type RobotListReq struct {
	UserID uint `json:"user_id" binding:"required"` // 关联用户ID（必填）
}

type RobotListRsp struct {
	RobotList []RobotCreateRsp `json:"robot_key_list"`
	Count     int64            `json:"count"`
}

func (app *App) RobotList(c *gin.Context) {
	var (
		req             RobotListReq
		robotCreateList []RobotCreateRsp
		robotList       []model.Robot
		count           int64
	)
	if err := c.BindJSON(&req); err != nil {
		logrus.Error("RobotList bind json failed: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	robotRepo := dao.NewRobotRepo(app.DB)
	robotList, count, err := robotRepo.ListRobotsByUserID(req.UserID, 1, 10)
	if err != nil {
		logrus.Error("ListRobotsByUserID failed: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i := range robotList {
		robotCreateList = append(robotCreateList, RobotCreateRsp{
			Name: robotList[i].Name,
			Id:   robotList[i].ID,
		})
	}
	c.JSON(http.StatusOK, gin.H{"code": 200,
		"message": "ok",
		"data": &RobotListRsp{
			Count:     count,
			RobotList: robotCreateList,
		}})
}
