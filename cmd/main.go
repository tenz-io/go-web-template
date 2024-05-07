package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tenz-io/gokit/cmd"

	"go-web-template/cmd/webgo"
	"go-web-template/internal/config"
)

var flags = []cmd.Flag{
	&cmd.IntFlag{
		Name:    "port",
		Aliases: []string{"p"},
		Value:   8090,
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
			cmd.WithDotEnvConfig(),
			cmd.WithYamlConfig(),
			cmd.WithLogger(true),
			cmd.WithAdminHTTPServer(),
			updateConfig(),
			server.Init(),
		},
		Run: server.Run(),
	}

	err := cmd.Run(app, flags)
	if err != nil {
		log.Fatal(err)
	}
}

func updateConfig() cmd.InitFunc {
	return func(c *cmd.Context, confPtr any) (cmd.CleanFunc, error) {
		var (
			cleanF = func(_ *cmd.Context) {}
		)

		conf, ok := confPtr.(*config.Config)
		if !ok {
			return cleanF, fmt.Errorf("invalid config type: %T", confPtr)
		}

		conf.Verbose = c.Bool(cmd.FlagNameVerbose)
		if c.IsSet("port") {
			conf.App.Port = c.Int("port")
		}
		conf.DB.Pass = os.Getenv("DB_PASS")

		return func(_ *cmd.Context) {
			log.Println("close server on port:", conf.App.Port)
		}, nil
	}
}
