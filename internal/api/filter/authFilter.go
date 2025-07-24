package filter

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"miniRustpbxgo/internal/utils"
	"net/http"
)

const SessionKey = "session_id"

type SessionAuth struct {
	rdb *redis.Client
}

func NewSessionAuth() *SessionAuth {
	s := &SessionAuth{}
	connRdb(s)
	return s
}

func (s *SessionAuth) Auth(ctx *gin.Context) {
	sessionID := ctx.GetHeader(SessionKey)
	if sessionID == "" {
		logrus.Error("session_id is empty")
		ctx.AbortWithStatusJSON(http.StatusForbidden, "sessionId is null")
		return
	}
	authKey := utils.GetAuthKey(sessionID)
	loginTime, err := s.rdb.Get(ctx, authKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logrus.Errorf("Get auth key %s error: %v", authKey, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "session auth error")
		return
	}
	if loginTime == "" {
		logrus.Error("session auth key not found")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "session auth fail")
		return
	}
	// ctx.Next() 只应该在中间件中使用
	ctx.Next()
}

func connRdb(s *SessionAuth) {
	// redis-cli
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	s.rdb = rdb
}
