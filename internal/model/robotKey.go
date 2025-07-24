package model

import (
	"gorm.io/gorm"
	"time"
)

// RobotKey 机器人密钥配置，与robotKeys表映射
type RobotKey struct {
	ID     uint   `gorm:"column:id;primaryKey"`    // 主键ID
	UserID uint   `gorm:"column:user_id;not null"` // 关联用户ID（外键）
	Name   string `gorm:"column:name;size:100"`    // 密钥名称（可选）

	// 大模型配置
	LLMProvider string `gorm:"column:llm_provider;size:100"` // 大模型提供商
	LLMApiKey   string `gorm:"column:llm_api_key;size:255"`  // 大模型API密钥
	LLMApiUrl   string `gorm:"column:llm_api_url;size:255"`  // 大模型API地址

	// 语音识别配置
	ASRProvider  string `gorm:"column:asr_provider;size:100"`             // 语音识别提供商
	ASRAppID     string `gorm:"column:asr_app_id;size:100"`               // 语音识别App ID
	ASRSecretID  string `gorm:"column:asr_secret_id;size:255"`            // 语音识别Secret ID
	ASRSecretKey string `gorm:"column:asr_secret_key;size:255"`           // 语音识别Secret Key
	ASRLanguage  string `gorm:"column:asr_language;size:20;default:'zh'"` // 语音识别语言，默认中文

	// 语音合成配置
	TTSProvider  string `gorm:"column:tts_provider;size:100"`   // 语音合成提供商
	TTSAppID     string `gorm:"column:tts_app_id;size:100"`     // 语音合成App ID
	TTSSecretID  string `gorm:"column:tts_secret_id;size:255"`  // 语音合成Secret ID
	TTSSecretKey string `gorm:"column:tts_secret_key;size:255"` // 语音合成Secret Key

	// 新增API密钥
	APIKey    string `gorm:"column:api_key;size:255"`    // API密钥
	APISecret string `gorm:"column:api_secret;size:255"` // API密钥Secret

	CreatedAt time.Time      `gorm:"column:created_at;size:255"`       // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at;size:255"`       // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"` // 软删除支持
}

// TableName 自定义表名
func (RobotKey) TableName() string {
	return "robotKeys"
}
