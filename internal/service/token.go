package service

import (
	"context"
	"errors"
	"time"

	"go-web-template/internal/model"
	"go-web-template/internal/repository/dao"
	"go-web-template/internal/util"

	"gorm.io/gorm"
)

var ErrTokenNotFound = errors.New("token not found")

// Token API Token 服务接口
type Token interface {
	CreateToken(ctx context.Context, userID int64, name string, token string, expiresAt *time.Time) (*model.APIToken, error)
	ListTokens(ctx context.Context, userID int64) ([]*model.APIToken, error)
	DeleteToken(ctx context.Context, userID, tokenID int64) error
	ValidateToken(ctx context.Context, token string) (*model.APIToken, error)
}

type tokenService struct {
	repo dao.APIToken
}

// NewToken 创建 Token 服务
func NewToken(repo dao.APIToken) Token {
	return &tokenService{repo: repo}
}

func (s *tokenService) CreateToken(ctx context.Context, userID int64, name string, token string, expiresAt *time.Time) (*model.APIToken, error) {
	if name == "" {
		return nil, errors.New("token name required")
	}

	record := &model.APIToken{
		UserID:    userID,
		Name:      name,
		TokenHash: util.HashToken(token),
		ExpiresAt: expiresAt,
	}

	if err := s.repo.Create(ctx, record); err != nil {
		return nil, err
	}

	return record, nil
}

func (s *tokenService) ListTokens(ctx context.Context, userID int64) ([]*model.APIToken, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *tokenService) DeleteToken(ctx context.Context, userID, tokenID int64) error {
	if err := s.repo.DeleteByID(ctx, userID, tokenID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTokenNotFound
		}
		return err
	}
	return nil
}

func (s *tokenService) ValidateToken(ctx context.Context, token string) (*model.APIToken, error) {
	if token == "" {
		return nil, errors.New("token is empty")
	}

	record, err := s.repo.FindByHash(ctx, util.HashToken(token))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	if record.ExpiresAt != nil && time.Now().After(*record.ExpiresAt) {
		return nil, errors.New("token expired")
	}

	return record, nil
}
