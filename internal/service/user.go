package service

import (
	"context"

	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/config"
	"go-web-template/internal/repository"
)

// User is the interface that provides user methods.
//
//go:generate mockery --name User --filename user_mock.go --inpackage
type User interface {
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
