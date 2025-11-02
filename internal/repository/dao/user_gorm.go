package dao

import (
	"context"
	"errors"
	"fmt"
	"go-web-template/internal/util"

	"gorm.io/gorm"

	"go-web-template/internal/model"
)

// user 用户仓库 GORM 实现
type user struct {
	db *gorm.DB
}

// NewUser 创建用户仓库
func NewUser(db *gorm.DB) User {
	return &user{db: db}
}

// GetByName 根据用户名获取用户
func (r *user) GetByName(ctx context.Context, name string) (*model.User, error) {
	var userModel model.User
	err := r.db.WithContext(ctx).Where("username = ?", name).First(&userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &userModel, nil
}

// GetByID 根据ID获取用户
func (r *user) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var userModel model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &userModel, nil
}

// Create 创建用户
func (r *user) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// UpdatePassword 更新用户密码
func (r *user) UpdatePassword(ctx context.Context, userID int64, newPassword string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Update("password", newPassword).Error
}

// VerifyUser 验证用户凭据
func (r *user) VerifyUser(ctx context.Context, username, password string) (*model.User, error) {
	userModel, err := r.GetByName(ctx, username)
	if err != nil {
		return nil, err
	}

	hashPass := util.HashPasswordWithSalt(password, userModel.Salt)

	if userModel.Password != hashPass {
		return nil, fmt.Errorf("invalid password")
	}

	return userModel, nil
}

// List 获取用户列表
func (r *user) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
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
func (r *user) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error
	return count, err
}

// Delete 删除用户
func (r *user) Delete(ctx context.Context, userID int64) error {
	err := r.db.WithContext(ctx).Delete(&model.User{}, userID).Error
	return err
}
