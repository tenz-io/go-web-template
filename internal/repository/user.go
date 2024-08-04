package repository

import (
	"context"
	"fmt"

	"go-web-template/internal/config"
	"go-web-template/internal/model"
)

// User is the interface that provides user methods.
//
//go:generate mockery --name User --filename user_mock.go --inpackage
type User interface {
	GetByName(ctx context.Context, name string) (model.User, error)
}

func NewUser() User {
	return &user{}
}

type user struct {
	cfg *config.Config
}

func (u *user) GetByName(_ context.Context, name string) (model.User, error) {
	return model.User{
		Userid:   123, // TODO replace with real user id
		Username: name,
		Profile:  fmt.Sprintf("profile of %s", name),
	}, nil
}
