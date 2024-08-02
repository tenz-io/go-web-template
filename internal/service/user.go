package service

import (
	"context"

	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/config"
	"go-web-template/internal/constant"
	"go-web-template/internal/model"
	"go-web-template/internal/repository"
)

// User is the interface that provides user methods.
//
//go:generate mockery --name User --filename user_mock.go --inpackage
type User interface {
	GetByName(ctx context.Context, name string) (model.User, error)
	VerifyAdmin(ctx context.Context, name, pass string) (ok bool, err error)
}

func NewUser(
	cfg *config.Config,
	userRepo repository.User,
) User {
	return &user{
		cfg:     cfg,
		useRepo: userRepo,
	}
}

type user struct {
	cfg     *config.Config
	useRepo repository.User
}

func (u *user) GetByName(ctx context.Context, name string) (model.User, error) {
	userProfile, err := u.useRepo.UserProfile(ctx, name)
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		Userid:   123,
		Role:     constant.RoleAdmin,
		Username: name,
		Profile:  userProfile,
	}, nil
}

func (u *user) VerifyAdmin(ctx context.Context, name, pass string) (ok bool, err error) {
	var (
		le = logger.FromContext(ctx)
	)
	if name == u.cfg.App.AdminUser && pass == u.cfg.App.AdminPass {
		le.Debug("admin verified")
		return true, nil
	}
	le.Debug("admin verify failed")
	return false, nil
}
