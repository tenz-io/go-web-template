package service

import (
	"context"
	"errors"
	"fmt"
	"go-web-template/internal/repository/dao"
	"go-web-template/internal/util"

	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/config"
	"go-web-template/internal/model"

	"gorm.io/gorm"
)

// User 用户服务接口
type User interface {
	// 认证相关
	VerifyUser(ctx context.Context, username, password string) (*model.User, error)

	// 用户管理
	Register(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	UpdatePassword(ctx context.Context, userID int64, req *model.UpdatePasswordRequest) error
	DeleteUser(ctx context.Context, userID int64) error
	GetUser(ctx context.Context, userID int64) (*model.User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]*model.User, int64, error)
}

func NewUserService(
	cfg *config.Config,
	userRepo dao.User,
) User {
	return &user{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

type user struct {
	cfg      *config.Config
	userRepo dao.User
}

// VerifyUser 验证用户凭据
func (u *user) VerifyUser(ctx context.Context, username, password string) (*model.User, error) {
	le := logger.FromContext(ctx)

	userModel, err := u.userRepo.VerifyUser(ctx, username, password)
	if err != nil {
		le.Debug("user verification failed")
		return nil, err
	}

	le.Debug("user verified successfully")
	return userModel, nil
}

// Register 创建用户
func (u *user) Register(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	le := logger.FromContext(ctx)

	// 检查用户名是否已存在
	_, err := u.userRepo.GetByName(ctx, req.Username)
	if err == nil {
		return nil, fmt.Errorf("username already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		le.WithError(err).Error("failed to check user existence")
		return nil, err
	}

	// 生成盐值
	salt := util.GenerateSalt(16)

	// 哈希密码（带盐值）
	hashedPassword := util.HashPasswordWithSalt(req.Password, salt)

	// 创建用户
	userModel := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Salt:     salt,
		Role:     req.Role,
		Profile:  "",
	}

	err = u.userRepo.Create(ctx, userModel)
	if err != nil {
		le.Error("failed to create user")
		return nil, err
	}

	le.Info("user created successfully")
	return userModel, nil
}

// UpdatePassword 更新用户密码
func (u *user) UpdatePassword(ctx context.Context, userID int64, req *model.UpdatePasswordRequest) error {
	le := logger.FromContext(ctx)

	// 获取用户信息
	userModel, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// 解码 base64 编码的哈希值
	hashedOldPass := util.HashPasswordWithSalt(req.OldPassword, userModel.Salt)

	// 验证旧密码（使用 HMAC-SHA256 + 盐值）
	if hashedOldPass != userModel.Password {
		le.Debug("old password is incorrect")
		return fmt.Errorf("旧密码不正确")
	}

	hashedNewPass := util.HashPasswordWithSalt(req.NewPassword, userModel.Salt)

	// 更新密码
	err = u.userRepo.UpdatePassword(ctx, userID, hashedNewPass)
	if err != nil {
		le.Error("failed to update password")
		return err
	}

	le.Info("password updated successfully")
	return nil
}

// GetUser 获取用户信息
func (u *user) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	return u.userRepo.GetByID(ctx, userID)
}

// ListUsers 获取用户列表
func (u *user) ListUsers(ctx context.Context, limit, offset int) ([]*model.User, int64, error) {
	users, err := u.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := u.userRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// DeleteUser 删除用户
func (u *user) DeleteUser(ctx context.Context, userID int64) error {
	le := logger.FromContext(ctx).WithField("userID", userID)

	// 调用 repository 删除用户
	err := u.userRepo.Delete(ctx, userID)
	if err != nil {
		le.Error("failed to delete user")
		return err
	}

	le.Info("user deleted successfully")
	return nil
}
