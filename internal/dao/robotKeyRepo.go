package dao

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"miniRustpbxgo/internal/model"
	"time"
)

type RobotKeyRepo struct {
	db *gorm.DB
}

func NewRobotKeyRepo(db *gorm.DB) *RobotKeyRepo {
	return &RobotKeyRepo{db: db}
}

// CreateRobotKey 创建机器人密钥
func (r *RobotKeyRepo) CreateRobotKey(robotKey *model.RobotKey) (*model.RobotKey, error) {
	// 设置创建时间和更新时间（GORM 也会自动填充，这里手动设置更明确）
	now := time.Now()
	robotKey.CreatedAt = now
	robotKey.UpdatedAt = now

	// 执行插入操作
	result := r.db.Create(robotKey)
	if result.Error != nil {
		logrus.Error("CreateRobotKey failed: ", result.Error)
		return nil, result.Error
	}
	return robotKey, nil
}

// GetRobotKeyByID 根据 ID 查询单条 RobotKey 记录
func (r *RobotKeyRepo) GetRobotKeyByID(id uint) (*model.RobotKey, error) {
	var robotKey model.RobotKey
	result := r.db.Where("id = ?", id).First(&robotKey)
	if result.Error != nil {
		logrus.Error("GetRobotKeyByID failed: ", result.Error)
		return nil, result.Error
	}
	return &robotKey, nil
}

// GetRobotKeyByAPIKey 根据 APIKey 查询记录（用于验证密钥有效性）
func (r *RobotKeyRepo) GetRobotKeyByAPIKey(apiKey string) (*model.RobotKey, error) {
	var robotKey model.RobotKey
	result := r.db.Where("api_key = ?", apiKey).First(&robotKey)
	if result.Error != nil {
		logrus.Error("GetRobotKeyByAPIKey failed: ", result.Error)
		return nil, result.Error
	}
	return &robotKey, nil
}

// ListRobotKeysByUserID 根据用户 ID 查询其所有 RobotKey 记录
func (r *RobotKeyRepo) ListRobotKeysByUserID(userID uint, page, pageSize int) ([]model.RobotKey, int64, error) {
	var (
		robotKeys []model.RobotKey
		total     int64
	)

	// 先查询总条数
	if err := r.db.Model(&model.RobotKey{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 再查询分页数据（按创建时间倒序，最新的在前）
	offset := (page - 1) * pageSize
	result := r.db.Where("user_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&robotKeys)

	if result.Error != nil {
		logrus.Error("ListRobotKeysByUserID failed: ", result.Error)
		return nil, 0, result.Error
	}
	return robotKeys, total, nil
}

// UpdateRobotKey 全量更新 RobotKey 记录（需传入完整结构体）
func (r *RobotKeyRepo) UpdateRobotKey(robotKey *model.RobotKey) error {
	// 更新时间戳
	robotKey.UpdatedAt = time.Now()

	// 全量更新（会覆盖所有字段，慎用）
	result := r.db.Save(robotKey)
	if result.Error != nil {
		logrus.Error("UpdateRobotKey failed: ", result.Error)
	}
	return result.Error
}

// UpdateRobotKeyPartial 部分更新 RobotKey 记录（只更新指定字段）
func (r *RobotKeyRepo) UpdateRobotKeyPartial(id uint, updates map[string]interface{}) error {
	// 强制更新时间戳
	updates["updated_at"] = time.Now()

	// 部分更新（只更新 map 中指定的字段）
	result := r.db.Model(&model.RobotKey{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		logrus.Error("UpdateRobotKeyPartial failed: ", result.Error)
	}
	return result.Error
}

// DeleteRobotKey 软删除 RobotKey 记录（不会从数据库中真正删除，而是更新 deleted_at 字段）
func (r *RobotKeyRepo) DeleteRobotKey(id uint) error {
	result := r.db.Delete(&model.RobotKey{}, id)
	if result.Error != nil {
		logrus.Error("DeleteRobotKey failed: ", result.Error)
	}
	return result.Error
}

// HardDeleteRobotKey 物理删除 RobotKey 记录（谨慎使用！会直接从数据库中删除数据）
func (r *RobotKeyRepo) HardDeleteRobotKey(id uint) error {
	result := r.db.Unscoped().Delete(&model.RobotKey{}, id)
	if result.Error != nil {
		logrus.Error("HardDeleteRobotKey failed: ", result.Error)
	}
	return result.Error
}
