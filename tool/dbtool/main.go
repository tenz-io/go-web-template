package main

import (
	"log"
	"time"

	"github.com/tenz-io/gokit/cmd"

	"go-web-template/internal/config"
)

var commands = []*cmd.Command{
	addUserCmd,
	verifyUserCmd,
	listUserCmd,
	changePasswdCmd,
}

var flags []cmd.Flag

func main() {
	app := cmd.App{
		Name:    "dbtool",
		Usage:   "DB tool",
		ConfPtr: &config.Config{},
		Inits: []cmd.InitFunc{
			cmd.WithDotEnvConfig(),
			cmd.WithYamlConfig(),
			cmd.WithUpdateConfigByEnv(),
			cmd.WithLogger(true),
		},
	}

	err := cmd.Run(app, flags, commands...)
	time.Sleep(1 * time.Second)
	if err != nil {
		log.Fatal(err)
	}
}
