package dao

import (
	"context"

	"go-web-template/internal/model"
)

// APIToken 仓库接口
type APIToken interface {
	Create(ctx context.Context, token *model.APIToken) error
	ListByUser(ctx context.Context, userID int64) ([]*model.APIToken, error)
	DeleteByID(ctx context.Context, userID, tokenID int64) error
	FindByHash(ctx context.Context, hash string) (*model.APIToken, error)
}
