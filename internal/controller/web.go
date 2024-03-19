package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go-web-template/internal/config"
)

type WebServer struct {
	engine *gin.Engine
	cfg    *config.Config
	api    *Api
}

func NewWebServer(
	cfg *config.Config,
	api *Api,
) *WebServer {
	return &WebServer{
		engine: gin.New(),
		cfg:    cfg,
		api:    api,
	}
}

func (ws *WebServer) Init() error {
	// auth middleware
	//ws.manager.Use(ws.user.applyAuth())

	tmplPattern := strings.Join([]string{
		ws.cfg.App.Web,
		"*.html",
	}, "/")

	static := strings.Join([]string{
		ws.cfg.App.Web,
		"static",
	}, "/")

	ws.engine.LoadHTMLGlob(tmplPattern)
	ws.engine.Static("/static", static)

	baseGrp := ws.engine.Group("")
	ws.registerHomepage(baseGrp)

	apiGroup := ws.engine.Group("api")
	ws.api.register(apiGroup)

	return nil
}

func (ws *WebServer) registerHomepage(rg *gin.RouterGroup) {
	rg.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"name": ws.cfg.App.Name,
		})
	})
}

func (ws *WebServer) Run(errC chan<- error) {
	addr := fmt.Sprintf(":%s", ws.cfg.App.Port)

	log.Println("[Controllers] http server listen on ", addr)

	err := ws.engine.Run(addr)

	if err != nil {
		errC <- err
	}

}
