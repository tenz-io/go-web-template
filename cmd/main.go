package main

import (
	"fmt"
	syslog "log"
	"os"

	"github.com/tenz-io/gokit/app"

	"go-web-template/internal/config"
	"go-web-template/internal/setup"
)

var flags = []app.Flag{
	&app.StringFlag{
		Name:  "env",
		Value: "test",
		Usage: "Environment",
	},
	&app.StringFlag{
		Name:  "port",
		Value: "8080",
		Usage: "http port",
	},
}

var (
	controllers *setup.Controllers
)

func main() {
	cfg := app.Config{
		Name:  "go-web-template",
		Usage: "Go Web Template",
		Conf:  &config.Config{},
		Inits: []app.InitFunc{
			app.InitYamlConfig,
			app.InitLogger,
			app.InitDefaultHandler,
			app.InitAdminHTTPServer,
			updateConfByFlags,
			initControllers,
		},
		Run: run,
	}

	app.Run(cfg, flags)
}

func run(c *app.Context, confPtr any, errC chan<- error) {
	controllers.Run(errC)
}

func updateConfByFlags(c *app.Context, confPtr any) (app.CleanFunc, error) {
	var (
		cleanF  = func() {}
		rawPort string
		protSrc string
	)

	conf, ok := confPtr.(*config.Config)
	if !ok {
		return cleanF, fmt.Errorf("invalid config type: %T", confPtr)
	}

	if v, err := c.GetFlags().Bool("verbose"); err == nil {
		conf.Verbose = v
		syslog.Printf("verbose: %v\n", v)
	}

	// get port from env variable
	if envPort := os.Getenv("PORT"); envPort != "" {
		rawPort = envPort
		protSrc = "env variable"
	}

	if argPort, err := c.GetFlags().String("port"); err == nil && argPort != "" {
		rawPort = argPort
		protSrc = "command line argument"
	}

	if rawPort != "" {
		conf.App.Port = rawPort
		syslog.Printf("source: (%s), port: %s\n", protSrc, rawPort)
	}

	return func() {
		syslog.Println("close port:", conf.App.Port)
	}, nil
}

func initControllers(c *app.Context, confPtr any) (app.CleanFunc, error) {
	var (
		cleanF = func() {}
		err    error
	)
	conf, ok := confPtr.(*config.Config)
	if !ok {
		return cleanF, fmt.Errorf("invalid config type: %T", confPtr)
	}

	controllers, err = setup.InitializeControllers(conf)
	if err != nil {
		return cleanF, fmt.Errorf("init controllers error: %w", err)
	}

	if err = controllers.Init(); err != nil {
		return cleanF, fmt.Errorf("init controllers error: %w", err)
	}

	return func() {
		controllers.Shutdown()
	}, nil
}
