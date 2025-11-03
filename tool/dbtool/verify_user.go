package main

import (
	"fmt"

	"github.com/tenz-io/gokit/cmd"

	"go-web-template/internal/config"
	"go-web-template/internal/repository/dao"
	"go-web-template/internal/service"
	"go-web-template/internal/setup"
)

var verifyUserCmd = &cmd.Command{
	Name:  "verify-user",
	Usage: "verify user",
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
	},
	Action: verifyUser,
}

func verifyUser(c *cmd.Context) error {
	user := c.String("user")
	if user == "" {
		return fmt.Errorf("user is required")
	}
	pass := c.String("pass")
	if pass == "" {
		return fmt.Errorf("pass is required")
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

	_, err = userService.VerifyUser(c.Context, user, pass)
	if err != nil {
		return fmt.Errorf("verify user error: %w", err)
	}

	fmt.Println("verify user success")

	return nil
}
