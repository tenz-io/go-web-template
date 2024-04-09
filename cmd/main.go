package main

import (
	"fmt"
	"log"

	"github.com/tenz-io/gokit/app"

	"go-web-template/cmd/webgo"
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
	server := webgo.NewServer()
	cfg := app.Config{
		Name:  "go-web-template",
		Usage: "Go Web Template",
		Conf:  &config.Config{},
		Inits: []app.InitFunc{
			app.WithYamlConfig(),
			app.WithLogger(true),
			app.WithAdminHTTPServer(),
			updateConfByFlags(),
			server.Init(),
		},
		Run: server.Run(),
	}

	app.Run(cfg, flags)
}

func updateConfByFlags() app.InitFunc {
	return func(c *app.Context, confPtr any) (app.CleanFunc, error) {
		var (
			cleanF = func() {}
		)

		conf, ok := confPtr.(*config.Config)
		if !ok {
			return cleanF, fmt.Errorf("invalid config type: %T", confPtr)
		}

		if v, err := c.GetFlags().Bool(app.FlagNameVerbose); err == nil {
			conf.Verbose = v
			log.Printf("verbose: %v\n", v)
		}

		if port, err := c.GetFlags().Int(app.FlagNamePort); err == nil {
			conf.App.Port = fmt.Sprintf("%d", port)
		}

		return func() {
			log.Println("close server on port:", conf.App.Port)
		}, nil
	}
}
