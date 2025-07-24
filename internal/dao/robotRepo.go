package dao

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"miniRustpbxgo/internal/model"
	"time"
)

type RobotRepo struct {
	db *gorm.DB
}

func NewRobotRepo(db *gorm.DB) *RobotRepo {
	return &RobotRepo{db: db}
}

// CreateRobot 创建新的Robot记录
func (r *RobotRepo) CreateRobot(robot *model.Robot) (*model.Robot, error) {
	now := time.Now()
	robot.CreatedAt = now
	robot.UpdatedAt = now

	result := r.db.Create(robot)
	if result.Error != nil {
		logrus.Error("CreateRobot Failed: ", result.Error)
		return nil, result.Error
	}
	return robot, nil
}

// GetRobotByID 根据ID查询单个Robot记录
func (r *RobotRepo) GetRobotByID(id uint) (*model.Robot, error) {
	var robot model.Robot
	result := r.db.Where("id = ?", id).First(&robot)
	if result.Error != nil {
		logrus.Error("GetRobotByID Failed: ", result.Error)
		return nil, result.Error
	}
	return &robot, nil
}

// GetRobotByUsrID 根据ID查询单个Robot记录
func (r *RobotRepo) GetRobotByUsrID(id uint) (*model.Robot, error) {
	var robot model.Robot
	result := r.db.Where("user_id = ?", id).First(&robot)
	if result.Error != nil {
		logrus.Error("GetRobotByUsrID Failed: ", result.Error)
		return nil, result.Error
	}
	return &robot, nil
}

// ListRobotsByUserID 分页查询指定用户的Robot列表
func (r *RobotRepo) ListRobotsByUserID(userID uint, page, pageSize int) ([]model.Robot, int64, error) {
	var (
		robots []model.Robot
		total  int64
	)

	// 查询总数
	if err := r.db.Model(&model.Robot{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		logrus.Error("ListRobotsByUserID", err)
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	result := r.db.Where("user_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&robots)

	if result.Error != nil {
		logrus.Error("ListRobotsByUserID", result.Error)
		return nil, 0, result.Error
	}
	return robots, total, nil
}

// UpdateRobotPartial 部分更新Robot记录（只更新提供的字段）
func (r *RobotRepo) UpdateRobotPartial(id uint, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	result := r.db.Model(&model.Robot{}).Where("id = ?", id).Updates(updates)
	return result.Error
}

// DeleteRobot 软删除Robot记录
func (r *RobotRepo) DeleteRobot(id uint) error {
	result := r.db.Delete(&model.Robot{}, id)
	return result.Error
}

// UpdateRobot 全量更新Robot记录
func (r *RobotRepo) UpdateRobot(robot *model.Robot) error {
	robot.UpdatedAt = time.Now()
	result := r.db.Save(robot)
	return result.Error
}
