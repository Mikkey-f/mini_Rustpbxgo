package utils

import "fmt"

func GetAuthKey(sessionId string) string {
	authKey := fmt.Sprintf("session_auth:%s", sessionId)
	return authKey
}

func GetSessionKey(sessionId string) string {
	sessionKey := fmt.Sprintf("session_auth:%s", sessionId)
	return sessionKey
}
