package repository

import (
	"context"

	"go-web-template/internal/model"
)

// User 用户仓库接口
type User interface {
	// 基础查询
	GetByName(ctx context.Context, name string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)

	// 用户管理
	Create(ctx context.Context, user *model.User) error
	UpdatePassword(ctx context.Context, userID int64, newPassword string) error
	Delete(ctx context.Context, userID int64) error

	// 认证
	VerifyUser(ctx context.Context, username, password string) (*model.User, error)

	// 列表查询
	List(ctx context.Context, limit, offset int) ([]*model.User, error)
	Count(ctx context.Context) (int64, error)
}
