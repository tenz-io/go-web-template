package main

import (
	"fmt"

	"github.com/tenz-io/gokit/cmd"

	"go-web-template/internal/config"
	"go-web-template/internal/constant"
	"go-web-template/internal/repository/dao"
	"go-web-template/internal/service"
	"go-web-template/internal/setup"
)

var addUserCmd = &cmd.Command{
	Name:  "add-user",
	Usage: "add user",
	Flags: []cmd.Flag{
		&cmd.StringFlag{
			Name:    "user",
			Aliases: []string{"u"},
			Usage:   "username",
		},
		&cmd.StringFlag{
			Name:    "pass",
			Aliases: []string{"p"},
			Usage:   "password",
			Value:   "123456",
		},
		&cmd.IntFlag{
			Name:    "role",
			Aliases: []string{"r"},
			Usage:   "role, 1: admin, 2: user, 4: user_plus, 8: openapi",
			Value:   2,
		},
	},
	Action: addUser,
}

func addUser(c *cmd.Context) error {
	user := c.String("user")
	if user == "" {
		return fmt.Errorf("user is required")
	}
	pass := c.String("pass")
	if pass == "" {
		return fmt.Errorf("pass is required")
	}

	role := constant.Role(c.Int("role"))
	if !role.IsValid() {
		return fmt.Errorf("role is invalid")
	}

	cfg, err := cmd.GetConfig[*config.Config](c)
	if err != nil {
		return fmt.Errorf("get config error: %w", err)
	}

	db, err := setup.ProvideDB(cfg)
	if err != nil {
		return fmt.Errorf("provide db error: %w", err)
	}
	userDao := dao.NewUserDao(db)
	userService := service.NewUserService(cfg, userDao)

	userModel, err := userService.CreateUser(c.Context, service.CreateUserParam{
		Username: user,
		Password: pass,
		Role:     int32(role),
	})
	if err != nil {
		return fmt.Errorf("register user error: %w", err)
	}

	fmt.Printf("register user success, userid: %d\n", userModel.ID)

	return nil
}
