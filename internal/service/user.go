package service

import (
	"context"

	"go-web-template/internal/model"
	"go-web-template/internal/repository"
)

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
		Username: name,
		Profile:  userProfile,
	}, nil
}
