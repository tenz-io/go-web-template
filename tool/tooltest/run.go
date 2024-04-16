package main

import (
	"github.com/tenz-io/gokit/app"
	"github.com/tenz-io/gokit/logger"
)

func initd() app.InitFunc {
	return func(c *app.Context, confPtr any) (app.CleanFunc, error) {
		logger.FromContext(c.Context).Info("init function called")
		return func() {}, nil
	}
}

func run() app.RunFunc {
	return func(c *app.Context, confPtr any, errC chan<- error) {
		logger.FromContext(c.Context).Info("run function called")
		errC <- nil
	}
}
