package main

import (
	"flag"
	"fmt"
	"go-web-template/internal/config"
	"go-web-template/internal/setup"
	"go-web-template/internal/util"
	"log"
)

var flags struct {
	config  string
	verbose bool
}

var (
	controllers *setup.Controllers
)

func init() {
	flag.StringVar(&flags.config, "config", "config/app.yaml", "config file path")
	flag.BoolVar(&flags.verbose, "verbose", true, "verbose")
}

func main() {
	log.Println("app starting")
	flag.Parse()
	log.Printf("flags: %+v\n", flags)

	// 1. init all
	if err := initAll(); err != nil {
		log.Fatalf("init error: %+v\n", err)
	}

	// 2. run
	errC := make(chan error)
	run(errC)

	// 3. wait signal
	util.WaitSignal(errC, func() {
		log.Println("prepare clean...")
		shutdown()
		log.Println("clean resources completed!")
	})

}

func initAll() (err error) {
	conf := config.NewConfig(flags.config)
	err = conf.Init()
	if err != nil {
		return fmt.Errorf("init config error: %w", err)
	}

	// init controllers
	controllers, err = setup.InitializeControllers(conf)
	if err != nil {
		return fmt.Errorf("init controllers error: %w", err)
	}

	if err = controllers.Init(); err != nil {
		return fmt.Errorf("init controllers error: %w", err)
	}

	return nil
}

func run(errC chan<- error) {
	log.Println("start to run")
	go controllers.Run(errC)
}

func shutdown() {
	log.Println("shutdown...")
	controllers.Shutdown()
}
