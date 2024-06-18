package model

import "go-web-template/internal/constant"

type User struct {
	Userid   int64
	Role     constant.Role
	Username string
	Profile  string
}
