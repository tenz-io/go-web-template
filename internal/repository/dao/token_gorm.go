package dao

import (
	"context"

	"go-web-template/internal/model"

	"gorm.io/gorm"
)

type apiTokenRepo struct {
	db *gorm.DB
}

// NewAPIToken 创建 API Token 仓库
func NewAPIToken(db *gorm.DB) APIToken {
	return &apiTokenRepo{db: db}
}

func (r *apiTokenRepo) Create(ctx context.Context, token *model.APIToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *apiTokenRepo) ListByUser(ctx context.Context, userID int64) ([]*model.APIToken, error) {
	var tokens []*model.APIToken
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&tokens).Error
	return tokens, err
}

func (r *apiTokenRepo) DeleteByID(ctx context.Context, userID, tokenID int64) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", tokenID, userID).
		Delete(&model.APIToken{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *apiTokenRepo) FindByHash(ctx context.Context, hash string) (*model.APIToken, error) {
	var token model.APIToken
	err := r.db.WithContext(ctx).Where("token_hash = ?", hash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}
