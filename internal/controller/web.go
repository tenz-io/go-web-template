package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/ginterceptor"

	"go-web-template/internal/config"
)

type WebServer struct {
	engine      *gin.Engine
	interceptor ginterceptor.Interceptor
	cfg         *config.Config
	api         *Api
}

func NewWebServer(
	cfg *config.Config,
	api *Api,
) *WebServer {
	if cfg.Verbose {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := &WebServer{
		engine: gin.New(),
		interceptor: ginterceptor.NewInterceptorWithOpts(
			ginterceptor.WithTracking(true),
			ginterceptor.WithTraffic(true),
			ginterceptor.WithMetrics(false),
			ginterceptor.WithTimeout(0),
		),
		cfg: cfg,
		api: api,
	}

	return engine
}

func (ws *WebServer) Init() error {
	ws.interceptor.Apply(ws.engine)

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
