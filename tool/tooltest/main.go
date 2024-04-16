package main

import (
	"github.com/tenz-io/gokit/app"

	"go-web-template/internal/config"
)

var flags = []app.Flag{
	&app.StringFlag{
		Name:  "env",
		Value: "test",
		Usage: "Environment",
	},
}

func main() {
	cfg := app.Config{
		Name:  "go-web-template",
		Usage: "Go Web Template",
		Conf:  &config.Config{},
		Inits: []app.InitFunc{
			app.WithLogger(true),
			app.WithYamlConfig(),
			initd(),
		},
		Run: run(),
	}

	app.Run(cfg, flags)
}
