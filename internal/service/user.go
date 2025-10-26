package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/config"
	"go-web-template/internal/constant"
	"go-web-template/internal/model"
	"go-web-template/internal/repository"
)

// User 用户服务接口
type User interface {
	// 认证相关
	VerifyUser(ctx context.Context, username, password string) (*model.User, error)
	VerifyAdmin(ctx context.Context, username, password string) (bool, error)

	// 用户管理
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	UpdatePassword(ctx context.Context, userID int64, req *model.UpdatePasswordRequest) error
	DeleteUser(ctx context.Context, userID int64) error
	GetUser(ctx context.Context, userID int64) (*model.User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]*model.User, int64, error)
}

func NewUser(
	cfg *config.Config,
	userRepo repository.User,
) User {
	return &user{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

type user struct {
	cfg      *config.Config
	userRepo repository.User
}

// VerifyUser 验证用户凭据
func (u *user) VerifyUser(ctx context.Context, username, password string) (*model.User, error) {
	le := logger.FromContext(ctx)

	user, err := u.userRepo.VerifyUser(ctx, username, password)
	if err != nil {
		le.Debug("user verification failed")
		return nil, err
	}

	le.Debug("user verified successfully")
	return user, nil
}

// VerifyAdmin 验证管理员凭据
func (u *user) VerifyAdmin(ctx context.Context, username, password string) (bool, error) {
	le := logger.FromContext(ctx)

	user, err := u.userRepo.VerifyUser(ctx, username, password)
	if err != nil {
		le.Debug("admin verification failed")
		return false, err
	}

	if user.Role != int32(constant.RoleAdmin) {
		le.Debug("user is not admin")
		return false, nil
	}

	le.Debug("admin verified successfully")
	return true, nil
}

// CreateUser 创建用户
func (u *user) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	le := logger.FromContext(ctx)

	// 检查用户名是否已存在
	_, err := u.userRepo.GetByName(ctx, req.Username)
	if err == nil {
		return nil, fmt.Errorf("username already exists")
	}

	// 生成盐值
	salt := generateSalt()

	// 哈希密码（带盐值）
	hashedPassword, err := hashPasswordWithSalt(req.Password, salt)
	if err != nil {
		le.Error("failed to hash password")
		return nil, fmt.Errorf("密码哈希失败")
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Salt:     salt,
		Role:     req.Role,
		Email:    req.Email,
		Profile:  "",
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		le.Error("failed to create user")
		return nil, err
	}

	le.Info("user created successfully")
	return user, nil
}

// UpdatePassword 更新用户密码
func (u *user) UpdatePassword(ctx context.Context, userID int64, req *model.UpdatePasswordRequest) error {
	le := logger.FromContext(ctx)

	// 获取用户信息
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// 解码 base64 编码的哈希值
	hashedBytes, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		return fmt.Errorf("invalid password format")
	}

	// 验证旧密码（使用 HMAC-SHA256 + 盐值）
	h := hmac.New(sha256.New, []byte(user.Salt))
	h.Write([]byte(req.OldPassword))
	expectedHash := h.Sum(nil)

	if !hmac.Equal(hashedBytes, expectedHash) {
		return fmt.Errorf("invalid old password")
	}

	// 哈希新密码（带盐值）
	hashedPassword, err := hashPasswordWithSalt(req.NewPassword, user.Salt)
	if err != nil {
		le.Error("failed to hash new password")
		return fmt.Errorf("密码哈希失败")
	}

	// 更新密码
	err = u.userRepo.UpdatePassword(ctx, userID, hashedPassword)
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

// generateSalt 生成随机盐值
func generateSalt() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
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
