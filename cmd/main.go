package main

import (
	"fmt"
	"log"

	"github.com/tenz-io/gokit/cmd"

	"go-web-template/cmd/webgo"
	"go-web-template/internal/config"
)

var flags = []cmd.Flag{
	&cmd.StringFlag{
		Name:    "env",
		Aliases: []string{"e"},
		Value:   "test",
		Usage:   "Environment",
	},
	&cmd.IntFlag{
		Name:    "port",
		Aliases: []string{"p"},
		Value:   8080,
		Usage:   "Port",
	},
}

func main() {
	server := webgo.NewServer()
	app := cmd.App{
		Name:    "go-web-template",
		Usage:   "Go Web Template",
		ConfPtr: &config.Config{},
		Inits: []cmd.InitFunc{
			cmd.WithYamlConfig(),
			cmd.WithLogger(true),
			cmd.WithAdminHTTPServer(),
			updateConfByFlags(),
			server.Init(),
		},
		Run: server.Run(),
	}

	cmd.Run(app, flags)
}

func updateConfByFlags() cmd.InitFunc {
	return func(c *cmd.Context, confPtr any) (cmd.CleanFunc, error) {
		var (
			cleanF = func(_ *cmd.Context) {}
		)

		conf, ok := confPtr.(*config.Config)
		if !ok {
			return cleanF, fmt.Errorf("invalid config type: %T", confPtr)
		}

		conf.Verbose = c.Bool(cmd.FlagNameVerbose)
		log.Printf("verbose: %v\n", conf.Verbose)

		port := c.Int("port")
		conf.App.Port = fmt.Sprintf("%d", port)

		return func(_ *cmd.Context) {
			log.Println("close server on port:", conf.App.Port)
		}, nil
	}
}
