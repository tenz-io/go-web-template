package main

import (
	"github.com/tenz-io/gokit/cmd"
	"github.com/tenz-io/gokit/logger"
)

func initd() cmd.InitFunc {
	return func(c *cmd.Context, confPtr any) (cmd.CleanFunc, error) {
		logger.FromContext(c.Context).Info("init function called")
		return func(_ *cmd.Context) {}, nil
	}
}

func run() cmd.RunFunc {
	return func(c *cmd.Context, confPtr any, errC chan<- error) {
		logger.FromContext(c.Context).Info("run function called")
		errC <- nil
	}
}
