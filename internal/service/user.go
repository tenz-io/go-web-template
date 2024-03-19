package service

import (
	"go-web-template/internal/model"
	"go-web-template/internal/repository"
)

type User interface {
	GetByName(name string) (model.User, error)
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

func (u *user) GetByName(name string) (model.User, error) {
	return model.User{
		Username: name,
		Profile:  u.useRepo.UserProfile(name),
	}, nil
}
