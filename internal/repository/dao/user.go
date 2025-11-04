package dao

import (
	"context"
	"errors"
	"go-web-template/internal/model/db"

	"gorm.io/gorm"
)

// User 用户仓库接口
type User interface {
	// 基础查询
	GetByName(ctx context.Context, name string) (*db.User, error)
	GetByID(ctx context.Context, id int64) (*db.User, error)

	// 用户管理
	Create(ctx context.Context, user *db.User) error
	UpdatePassword(ctx context.Context, userID int64, newPassword string) error
	Delete(ctx context.Context, userID int64) error

	// 列表查询
	List(ctx context.Context, limit, offset int) ([]*db.User, error)
	Count(ctx context.Context) (int64, error)
}

// user 用户仓库 GORM 实现
type user struct {
	db *gorm.DB
}

// NewUserDao 创建用户仓库
func NewUserDao(db *gorm.DB) User {
	return &user{db: db}
}

// GetByName 根据用户名获取用户
func (r *user) GetByName(ctx context.Context, name string) (*db.User, error) {
	var userModel db.User
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
func (r *user) GetByID(ctx context.Context, id int64) (*db.User, error) {
	var userModel db.User
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
func (r *user) Create(ctx context.Context, user *db.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// UpdatePassword 更新用户密码
func (r *user) UpdatePassword(ctx context.Context, userID int64, newPassword string) error {
	return r.db.WithContext(ctx).Model(&db.User{}).
		Where("id = ?", userID).
		Update("password", newPassword).Error
}

// List 获取用户列表
func (r *user) List(ctx context.Context, limit, offset int) ([]*db.User, error) {
	var users []*db.User
	err := r.db.WithContext(ctx).
		Select("id, username, role, profile, created_at, updated_at").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	return users, err
}

// Count 获取用户总数
func (r *user) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&db.User{}).Count(&count).Error
	return count, err
}

// Delete 删除用户
func (r *user) Delete(ctx context.Context, userID int64) error {
	err := r.db.WithContext(ctx).Delete(&db.User{}, userID).Error
	return err
}
