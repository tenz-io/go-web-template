package webgo

import (
	"fmt"

	"github.com/tenz-io/gokit/cmd"

	"go-web-template/internal/config"
	"go-web-template/internal/setup"
)

type Server struct {
	controllers *setup.Controllers
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() cmd.InitFunc {
	return func(c *cmd.Context, confPtr any) (cmd.CleanFunc, error) {
		var (
			cleanF = func(_ *cmd.Context) {}
			err    error
		)
		conf, ok := confPtr.(*config.Config)
		if !ok {
			return cleanF, fmt.Errorf("invalid config type: %T", confPtr)
		}

		s.controllers, err = setup.InitializeControllers(conf)
		if err != nil {
			return cleanF, fmt.Errorf("init controllers error: %w", err)
		}

		if err = s.controllers.Init(); err != nil {
			return cleanF, fmt.Errorf("init controllers error: %w", err)
		}

		return func(*cmd.Context) {
			s.controllers.Shutdown()
		}, nil
	}

}

func (s *Server) Run() cmd.RunFunc {
	return func(c *cmd.Context, confPtr any, errC chan<- error) {
		s.controllers.Run(errC)
	}
}
