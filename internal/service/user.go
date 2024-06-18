package service

import (
	"context"

	"go-web-template/internal/constant"
	"go-web-template/internal/model"
	"go-web-template/internal/repository"
)

// User is the interface that provides user methods.
//
//go:generate mockery --name User --filename user_mock.go --inpackage
type User interface {
	GetByName(ctx context.Context, name string) (model.User, error)
}

func NewUser(
	userRepo repository.User,
) User {
	return &user{
		useRepo: userRepo,
	}
}

type user struct {
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
