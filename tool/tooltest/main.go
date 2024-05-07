package main

import (
	"log"

	"github.com/tenz-io/gokit/cmd"

	"go-web-template/internal/config"
)

var commands = []*cmd.Command{
	getCmd,
	setCmd,
}

var flags = []cmd.Flag{
	&cmd.StringFlag{
		Name:  "foo",
		Value: "bar",
		Usage: "foo bar",
	},
}

func main() {
	app := cmd.App{
		Name:    "tooltest",
		Usage:   "Tool Test",
		ConfPtr: &config.Config{},
		Inits: []cmd.InitFunc{
			cmd.WithYamlConfig(),
		},
	}

	err := cmd.Run(app, flags, commands...)
	if err != nil {
		log.Fatal(err)
	}
}
