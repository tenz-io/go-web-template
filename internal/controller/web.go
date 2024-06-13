package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	pbapp "go-web-template/api/http/app"
	"go-web-template/internal/config"
)

type WebServer struct {
	engine *gin.Engine
	cfg    *config.Config
	api    *ApiServer
}

func NewWebServer(
	cfg *config.Config,
	apiServer *ApiServer,
) *WebServer {
	if cfg.Verbose {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	ws := &WebServer{
		engine: gin.New(),
		cfg:    cfg,
		api:    apiServer,
	}

	return ws
}

func (ws *WebServer) Init() error {
	if ws.cfg.Verbose {
		ws.engine.Use(gin.Logger())
	}

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

	pbapp.RegisterApiServerHTTPServer(ws.engine, ws.api)

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
	addr := fmt.Sprintf(":%d", ws.cfg.App.Port)
	log.Println("[Controllers] http server listen on ", addr)
	err := ws.engine.Run(addr)
	if err != nil {
		errC <- err
	}

}
