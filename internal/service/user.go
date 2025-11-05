package service

import (
	"context"
	"errors"
	"fmt"
	"go-web-template/internal/model/db"
	"go-web-template/internal/repository/dao"
	"go-web-template/internal/util"

	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/config"

	"gorm.io/gorm"
)

// User 用户服务接口
type User interface {
	// 认证相关
	VerifyUser(ctx context.Context, param VerifyUserParam) (*db.User, error)

	// 用户管理
	CreateUser(ctx context.Context, param CreateUserParam) (*db.User, error)
	UpdatePassword(ctx context.Context, userID int64, param UpdatePasswordParam) error
	DeleteUser(ctx context.Context, userID int64) error
	GetUser(ctx context.Context, userID int64) (*db.User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]*db.User, int64, error)
}

func NewUserService(
	cfg *config.Config,
	userDao dao.User,
) User {
	return &user{
		cfg:     cfg,
		userDao: userDao,
	}
}

type user struct {
	cfg     *config.Config
	userDao dao.User
}

// VerifyUser 验证用户凭据
func (u *user) VerifyUser(ctx context.Context, param VerifyUserParam) (*db.User, error) {
	le := logger.FromContext(ctx).WithField("username", param.Username)

	userModel, err := u.userDao.GetByName(ctx, param.Username)
	if err != nil {
		le.Info("user not found")
		return nil, err
	}

	hashPass := util.HashPasswordWithSalt(param.Password, userModel.Salt)

	if userModel.Password != hashPass {
		le.Warn("password not match")
		return nil, fmt.Errorf("invalid password")
	}

	return userModel, nil

}

// CreateUser 创建用户
func (u *user) CreateUser(ctx context.Context, param CreateUserParam) (*db.User, error) {
	le := logger.FromContext(ctx)

	// 检查用户名是否已存在
	_, err := u.userDao.GetByName(ctx, param.Username)
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
	hashedPassword := util.HashPasswordWithSalt(param.Password, salt)

	// 创建用户
	userModel := &db.User{
		Username: param.Username,
		Password: hashedPassword,
		Salt:     salt,
		Role:     param.Role,
		Profile:  "",
	}

	err = u.userDao.Create(ctx, userModel)
	if err != nil {
		le.Error("failed to create user")
		return nil, err
	}

	le.Info("user created successfully")
	return userModel, nil
}

// UpdatePassword 更新用户密码
func (u *user) UpdatePassword(ctx context.Context, userID int64, param UpdatePasswordParam) error {
	le := logger.FromContext(ctx)

	// 获取用户信息
	userModel, err := u.userDao.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// 解码 base64 编码的哈希值
	hashedOldPass := util.HashPasswordWithSalt(param.OldPassword, userModel.Salt)

	// 验证旧密码（使用 HMAC-SHA256 + 盐值）
	if hashedOldPass != userModel.Password {
		le.Debug("old password is incorrect")
		return fmt.Errorf("旧密码不正确")
	}

	hashedNewPass := util.HashPasswordWithSalt(param.NewPassword, userModel.Salt)

	// 更新密码
	err = u.userDao.UpdatePassword(ctx, userID, hashedNewPass)
	if err != nil {
		le.Error("failed to update password")
		return err
	}

	le.Info("password updated successfully")
	return nil
}

// GetUser 获取用户信息
func (u *user) GetUser(ctx context.Context, userID int64) (*db.User, error) {
	return u.userDao.GetByID(ctx, userID)
}

// ListUsers 获取用户列表
func (u *user) ListUsers(ctx context.Context, limit, offset int) ([]*db.User, int64, error) {
	users, err := u.userDao.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := u.userDao.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// DeleteUser 删除用户
func (u *user) DeleteUser(ctx context.Context, userID int64) error {
	le := logger.FromContext(ctx).WithField("userID", userID)

	// 调用 repository 删除用户
	err := u.userDao.Delete(ctx, userID)
	if err != nil {
		le.Error("failed to delete user")
		return err
	}

	le.Info("user deleted successfully")
	return nil
}
