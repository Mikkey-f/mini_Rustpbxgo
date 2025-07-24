package utils

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/sirupsen/logrus"
)

// GenerateSecureRandomString 生成指定长度的安全随机字符串
func GenerateSecureRandomString(length int) (string, error) {
	// 计算需要的随机字节数（base64编码后每4字节对应3字符，预留额外字节避免截断）
	b := make([]byte, length)
	_, err := rand.Read(b) // 从加密安全的随机源读取字节
	if err != nil {
		logrus.Errorf("GenerateSecureRandomString error:%v", err)
		return "", err
	}
	// 用base64编码转为字符串，再截断到指定长度
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}
