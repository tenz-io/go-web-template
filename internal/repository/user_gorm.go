package repository

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"gorm.io/gorm"

	"go-web-template/internal/model"
)

// UserGORM 用户仓库 GORM 实现
type UserGORM struct {
	db *gorm.DB
}

// NewUserGORM 创建用户仓库
func NewUserGORM(db *gorm.DB) User {
	return &UserGORM{db: db}
}

// GetByName 根据用户名获取用户
func (r *UserGORM) GetByName(ctx context.Context, name string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", name).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByID 根据ID获取用户
func (r *UserGORM) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *UserGORM) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// UpdatePassword 更新用户密码
func (r *UserGORM) UpdatePassword(ctx context.Context, userID int64, newPassword string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Update("password", newPassword).Error
}

// VerifyUser 验证用户凭据
func (r *UserGORM) VerifyUser(ctx context.Context, username, password string) (*model.User, error) {
	user, err := r.GetByName(ctx, username)
	if err != nil {
		return nil, err
	}

	// 解码 base64 编码的哈希值
	hashedBytes, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid password format")
	}

	// 使用 HMAC-SHA256 + 盐值验证密码
	h := hmac.New(sha256.New, []byte(user.Salt))
	h.Write([]byte(password))
	expectedHash := h.Sum(nil)

	if !hmac.Equal(hashedBytes, expectedHash) {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}

// List 获取用户列表
func (r *UserGORM) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	var users []*model.User
	err := r.db.WithContext(ctx).
		Select("id, username, role, email, profile, created_at, updated_at").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	return users, err
}

// Count 获取用户总数
func (r *UserGORM) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error
	return count, err
}

// Delete 删除用户
func (r *UserGORM) Delete(ctx context.Context, userID int64) error {
	err := r.db.WithContext(ctx).Delete(&model.User{}, userID).Error
	return err
}

// hashPasswordWithSalt 使用 HMAC-SHA256 + 盐值哈希密码，并进行 base64 编码
func hashPasswordWithSalt(password, salt string) (string, error) {
	// 使用 HMAC-SHA256 生成哈希
	h := hmac.New(sha256.New, []byte(salt))
	h.Write([]byte(password))
	hashedBytes := h.Sum(nil)

	// 对哈希结果进行 base64 编码
	return base64.StdEncoding.EncodeToString(hashedBytes), nil
}
