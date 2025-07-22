package model

import (
	"gorm.io/gorm"
	"time"
)

// User 对应数据库中的 users 表，使用 GORM 标签配置
type User struct {
	ID        uint           `gorm:"column:id;primaryKey"`                       // 主键
	Username  string         `gorm:"column:username;size:50;not null;unique"`    // 用户名，唯一非空
	Password  string         `gorm:"column:password;size:255;not null" json:"-"` // 密码，序列化忽略
	Email     string         `gorm:"column:email;size:100;not null;unique"`      // 邮箱，唯一非空
	Phone     *string        `gorm:"column:phone;size:20;unique"`                // 手机号，可选
	Nickname  *string        `gorm:"column:nickname;size:50"`                    // 昵称，可选
	Status    int8           `gorm:"column:status;default:1"`                    // 状态，默认正常
	CreatedAt time.Time      `gorm:"column:created_at"`                          // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at"`                          // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`           // 软删除支持
}

// TableName 自定义表名
func (User) TableName() string {
	return "users"
}
