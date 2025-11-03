package main

import (
	"fmt"
	"go-web-template/internal/model"

	"github.com/tenz-io/gokit/cmd"

	"go-web-template/internal/config"
	"go-web-template/internal/repository/dao"
	"go-web-template/internal/service"
	"go-web-template/internal/setup"
)

var changePasswdCmd = &cmd.Command{
	Name:  "change-passwd",
	Usage: "change password",
	Flags: []cmd.Flag{
		&cmd.StringFlag{
			Name:    "user",
			Aliases: []string{"u"},
			Usage:   "username",
			Value:   "",
		},
		&cmd.StringFlag{
			Name:    "old-passwd",
			Aliases: []string{"o"},
			Usage:   "old password",
			Value:   "",
		},
		&cmd.StringFlag{
			Name:    "new-passwd",
			Aliases: []string{"n"},
			Usage:   "new password",
			Value:   "",
		},
	},
	Action: changePasswd,
}

func changePasswd(c *cmd.Context) error {
	username := c.String("user")
	oldPasswd := c.String("old-passwd")
	newPasswd := c.String("new-passwd")

	if username == "" {
		return fmt.Errorf("user name is required")
	}
	if oldPasswd == "" {
		return fmt.Errorf("old password is required")
	}
	if newPasswd == "" {
		return fmt.Errorf("new password is required")
	}

	if oldPasswd == newPasswd {
		return fmt.Errorf("new password is the same as old password")
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

	userModel, err := userDao.GetByName(c.Context, username)
	if err != nil {
		return fmt.Errorf("get user error: %w", err)
	}

	err = userService.UpdatePassword(c.Context, userModel.ID, &model.UpdatePasswordRequest{
		OldPassword: oldPasswd,
		NewPassword: newPasswd,
	})
	if err != nil {
		return fmt.Errorf("update password error: %w", err)

	}

	fmt.Println("change password success")
	return nil
}
