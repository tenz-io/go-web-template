package repository

import (
	"context"
	"fmt"
)

// User is the interface that provides user methods.
//
//go:generate mockery --name User --filename user_mock.go --inpackage
type User interface {
	UserProfile(ctx context.Context, name string) (profile string, err error)
}

func NewUser() User {
	return &user{}
}

type user struct {
}

func (u *user) UserProfile(_ context.Context, name string) (profile string, err error) {
	return fmt.Sprintf("my name is %s", name), nil
}
