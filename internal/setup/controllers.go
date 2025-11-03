package setup

import (
	"fmt"
	"go-web-template/internal/job"
	"log"

	"go-web-template/internal/controller"
)

type Controllers struct {
	webServer *controller.WebServer
	jobMgr    *job.Manager
}

func (c *Controllers) Init() error {
	log.Println("[Controllers] Init...")
	if err := c.webServer.Init(); err != nil {
		return err
	}

	if err := c.jobMgr.Init(); err != nil {
		return fmt.Errorf("job manager init error: %w", err)
	}

	return nil
}

func (c *Controllers) Run(errC chan<- error) {
	log.Println("[Controllers] Run...")
	c.webServer.Run(errC)
}

func (c *Controllers) Shutdown() {
	log.Println("[Controllers] Shutdown...")
}
