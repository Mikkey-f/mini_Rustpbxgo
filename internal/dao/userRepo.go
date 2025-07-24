package dao

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"miniRustpbxgo/internal/model"
)

// UserRepo 定义 User 表的数据访问对象
type UserRepo struct {
	db *gorm.DB // 数据库连接实例，通过外部注入
}

// NewUserRepo 创建 UserRepo 实例
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create 1. 创建用户（Create）
func (r *UserRepo) Create(user *model.User) (error, error) {
	// 基础参数校验（非空判断）
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return errors.New("username, email and password are required"), nil
	}
	// 调用 GORM 创建记录
	return r.db.Create(user).Error, nil
}

// GetByID 2. 根据 ID 查询用户（Read）
func (r *UserRepo) GetByID(id uint) (*model.User, error) {
	var user model.User
	// 按 ID 查询，自动过滤软删除记录
	result := r.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 记录不存在，返回 nil（或自定义错误）
		}
		return nil, result.Error // 其他查询错误
	}
	return &user, nil
}

// GetByUsername 3. 根据用户名查询用户（用于登录等场景）
func (r *UserRepo) GetByUsername(username string) (*model.User, error) {
	var user model.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// List 4. 分页查询用户列表（Read）
func (r *UserRepo) List(page, pageSize int) ([]model.User, int64, error) {
	var (
		users []model.User
		total int64
	)

	// 先查询总数（排除密码字段，提高效率）
	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 计算分页偏移量
	offset := (page - 1) * pageSize
	// 分页查询（指定返回字段，排除密码）
	err := r.db.Select("id, username, email, nickname, status, created_at").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	return users, total, err
}

// Update 5. 更新用户信息（Update）
func (r *UserRepo) Update(id uint, updates map[string]interface{}) error {
	// 禁止更新敏感字段（如密码单独处理）
	if _, ok := updates["password"]; ok {
		return errors.New("password update not allowed here")
	}
	// 基于 ID 部分更新
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

// UpdatePassword 6. 单独更新密码（Update）
func (r *UserRepo) UpdatePassword(id uint, newPassword string) error {
	if newPassword == "" {
		return errors.New("password cannot be empty")
	}
	// 实际业务中需加密密码（此处简化）
	return r.db.Model(&model.User{}).Where("id = ?", id).
		Update("password", newPassword).Error
}

// Delete 7. 软删除用户（Delete）
func (r *UserRepo) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// HardDelete 8. 物理删除用户（谨慎使用）
func (r *UserRepo) HardDelete(id uint) error {
	return r.db.Unscoped().Delete(&model.User{}, id).Error
}

func (r *UserRepo) IsExist(username string, email string) (bool, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		logrus.Error("username isExist = [%v]", err)
		return false, err
	}

	err = r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		logrus.Error("email isExist = [%v]", err)
		return false, err
	}
	return true, nil
}
