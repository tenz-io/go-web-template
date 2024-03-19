package repository

import "fmt"

type User interface {
	UserProfile(name string) (profile string)
}

func NewUser() User {
	return &user{}
}

type user struct {
}

func (u *user) UserProfile(name string) (profile string) {
	return fmt.Sprintf("my name is %s", name)
}
