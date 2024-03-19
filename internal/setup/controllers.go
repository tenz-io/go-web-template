package setup

import (
	"log"

	"go-web-template/internal/controller"
)

type Controllers struct {
	webServer *controller.WebServer
}

func (c *Controllers) Init() error {
	log.Println("[Controllers] Init...")
	if err := c.webServer.Init(); err != nil {
		return err
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
