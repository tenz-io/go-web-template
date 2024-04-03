package webgo

import (
	"fmt"

	"github.com/tenz-io/gokit/app"

	"go-web-template/internal/config"
	"go-web-template/internal/setup"
)

type Server struct {
	controllers *setup.Controllers
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init(c *app.Context, confPtr any) (app.CleanFunc, error) {
	var (
		cleanF = func() {}
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

	return func() {
		s.controllers.Shutdown()
	}, nil
}

func (s *Server) Run(c *app.Context, confPtr any, errC chan<- error) {
	s.controllers.Run(errC)
}
