package model

import (
	"gorm.io/gorm"
	"time"
)

// Robot 机器人配置，与robots表映射
type Robot struct {
	ID           uint           `gorm:"column:id;primaryKey"` // 主键ID
	Name         string         `gorm:"column:name"`
	UserID       uint           `gorm:"column:user_id;not null"`          // 关联用户ID（外键）
	Speed        float32        `gorm:"column:speed;type:float"`          // 语音语速（可选）
	Volume       int            `gorm:"column:volume;type:int"`           // 语音音量（可选）
	Speaker      string         `gorm:"column:speaker;size:50"`           // 发音人（可选）
	Emotion      string         `gorm:"column:emotion;size:50"`           // 语音情感（可选）
	SystemPrompt string         `gorm:"column:system_prompt;type:text"`   // 系统提示词（可选）
	CreatedAt    time.Time      `gorm:"column:created_at"`                // 创建时间
	UpdatedAt    time.Time      `gorm:"column:updated_at"`                // 更新时间
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"` // 软删除支持
}

// TableName 自定义表名
func (Robot) TableName() string {
	return "robots"
}
