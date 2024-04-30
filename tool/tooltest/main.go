package main

import (
	"github.com/tenz-io/gokit/cmd"

	"go-web-template/internal/config"
)

var commands = []*cmd.Command{
	getCmd,
	setCmd,
}

var flags = []cmd.Flag{
	&cmd.StringFlag{
		Name:  "env",
		Value: "test",
		Usage: "Environment",
	},
}

func main() {
	app := cmd.App{
		Name:    "go-web-template",
		Usage:   "Go Web Template",
		ConfPtr: &config.Config{},
		Inits: []cmd.InitFunc{
			cmd.WithLogger(true),
			cmd.WithYamlConfig(),
			initd(),
		},
		Run: run(),
	}

	cmd.Run(app, flags, commands...)
}
