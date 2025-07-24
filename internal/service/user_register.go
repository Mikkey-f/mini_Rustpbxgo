package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"miniRustpbxgo/internal/dao"
	"miniRustpbxgo/internal/model"
	"net/http"
)

type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=50" comment:"用户名"` // 用户名（必填，3-50字符）
	Password string `json:"password" binding:"required,min=6,max=100" comment:"密码"` // 密码（必填，6-100字符）
	Email    string `json:"email" binding:"required,email" comment:"邮箱"`            // 邮箱（必填，需符合邮箱格式）
	Phone    string `json:"phone" binding:"omitempty,len=11" comment:"手机号（可选，11位）"` // 手机号（可选，11位数字）
	Nickname string `json:"nickname" binding:"omitempty,max=50" comment:"昵称（可选）"`   // 昵称（可选，最多50字符）
}

type RegisterRsp struct {
	Message string `json:"message" binding:"required"`
}

func (app *App) Register(ctx *gin.Context) {
	var req RegisterReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//密码加密
	hashedPassword, err := encryptPassword(req.Password)
	if err != nil {
		logrus.Errorf("Register encryptPassword error:%v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logrus.Info("hashedPassword:", hashedPassword)
	//账号校验
	userDao := dao.NewUserRepo(app.DB)
	exist, err := userDao.IsExist(req.Username, req.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exist {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "账号已经存在"})
		return
	}

	//账号信息持久化
	if _, err := userDao.Create(&model.User{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Email:    req.Email,
		Status:   1,
	}); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
		"data": &RegisterRsp{
			Message: fmt.Sprintf("注册成功"),
		},
	})
}

func encryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error("bcrypt generate from password error = %v", err)
		return "", err
	}
	return string(hashedPassword), nil
}
