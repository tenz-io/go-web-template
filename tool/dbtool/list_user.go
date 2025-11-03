package main

import (
	"encoding/json"
	"fmt"
	"go-web-template/internal/service"

	"github.com/tenz-io/gokit/cmd"

	"go-web-template/internal/config"
	"go-web-template/internal/repository/dao"
	"go-web-template/internal/setup"
)

var listUserCmd = &cmd.Command{
	Name:  "list-user",
	Usage: "list users",
	Flags: []cmd.Flag{
		&cmd.IntFlag{
			Name:    "cursor",
			Aliases: []string{"c"},
			Usage:   "cursor",
			Value:   0,
		},
		&cmd.IntFlag{
			Name:    "limit",
			Aliases: []string{"l"},
			Usage:   "limit",
			Value:   10,
		},
	},
	Action: listUser,
}

func listUser(c *cmd.Context) error {
	cursor := c.Int("cursor")
	limit := c.Int("limit")

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

	userList, total, err := userService.ListUsers(c.Context, cursor, limit)
	if err != nil {
		return fmt.Errorf("list user error: %w", err)
	}

	fmt.Printf("total: %d\n", total)
	for i, userModel := range userList {
		j, _ := json.Marshal(userModel)
		fmt.Printf("%d: %s\n", i, j)
	}

	return nil
}
