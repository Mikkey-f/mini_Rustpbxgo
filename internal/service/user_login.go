package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"miniRustpbxgo/internal/dao"
	"net/http"
	"strconv"
	"time"
)

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRsp struct {
	Message   string `json:"message" binding:"required"`
	SessionID string `json:"session_id" binding:"required"`
}

func (app *App) Login(ctx *gin.Context) {
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Error("login req error:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var (
		username = req.Username
		password = req.Password
	)

	userDao := dao.NewUserRepo(app.DB)
	user, err := userDao.GetByUsername(username)
	if err != nil {
		logrus.Error("userDao.GetByUsername Username not found:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "请输入正确的账号ID"})
		return
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password)); err != nil {
		logrus.Error("Login CompareHashAndPassword error:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "密码错误"})
		return
	}
	sessionId, err := app.generateSessionId(context.Background(), strconv.Itoa(int(user.ID)))
	if err != nil {
		logrus.Error("generateSessionId error:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "系统错误，稍后重试"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
		"data": &LoginRsp{
			SessionID: sessionId,
		},
	})
	return
}

func (app *App) generateSessionId(ctx context.Context, userId string) (string, error) {
	sessionId := uuid.New().String()
	// key : session_id:{user_id} val : session_id
	sessionKey := fmt.Sprintf("session_id:%s", userId)
	err := app.Rdb.Set(ctx, sessionKey, sessionId, time.Hour*8).Err()
	if err != nil {
		fmt.Printf("rdb set error = %v \n", err)
		return "", err
	}

	authKey := fmt.Sprintf("session_auth:%s", sessionId)
	err = app.Rdb.Set(ctx, authKey, time.Now().Unix(), time.Hour*8).Err()
	if err != nil {
		fmt.Printf("rdb set error = %v \n", err)
		return "", err
	}
	fmt.Println(sessionKey)
	fmt.Println(authKey)
	return sessionId, nil
}
