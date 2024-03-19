package util

import (
	"log"
	"os"
	"os/signal"
)

func WaitSignal(errC <-chan error, hook func()) {
	signC := make(chan os.Signal, 1)
	signal.Notify(signC, os.Interrupt, os.Kill)
	select {
	case <-signC:
		log.Println("received interrupt signal")
		hook()
		os.Exit(0)
	case err := <-errC:
		log.Printf("run error: %+v", err)
		hook()
		os.Exit(1)
	}
}
