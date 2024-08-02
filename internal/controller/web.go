package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/ginext"

	pbapp "go-web-template/api/http/app"
	"go-web-template/internal/config"
)

type WebServer struct {
	engine *gin.Engine
	cfg    *config.Config
	api    *ApiServer
	admin  *AdminServer
}

func NewWebServer(
	cfg *config.Config,
	apiServer *ApiServer,
	adminServer *AdminServer,
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
		admin:  adminServer,
	}

	return ws
}

func (ws *WebServer) Init() error {
	if ws.cfg.Verbose {
		ws.engine.Use(gin.Logger())
	}

	// set app secret key
	ginext.InitJWT(ws.cfg.App.Secret)

	// register middleware template and static
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

	adminBaseGrp := ws.engine.Group("admin")
	ws.registerAdminPages(adminBaseGrp)

	pbapp.RegisterApiServerHTTPServer(ws.engine, ws.api)
	pbapp.RegisterAdminServerHTTPServer(ws.engine, ws.admin)

	return nil
}

func (ws *WebServer) registerHomepage(rg *gin.RouterGroup) {
	rg.GET("/",
		RedirectHandler("/login"),
		ginext.Authenticate(ginext.RoleUser, ginext.AuthTypeWeb),
		func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"name": ws.cfg.App.Name,
			})
		})
}

func (ws *WebServer) registerAdminPages(rg *gin.RouterGroup) {
	rg.GET("/",
		RedirectHandler("/admin/login"),
		ginext.Authenticate(ginext.RoleAdmin, ginext.AuthTypeWeb),
		func(c *gin.Context) {
			c.HTML(http.StatusOK, "admin_index.html", gin.H{
				"name": ws.cfg.App.Name,
			})
		})

	rg.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_login.html", gin.H{
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

// RedirectHandler handler
func RedirectHandler(location string) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Next()

		// if response status is 401, redirect to login page
		if c.Writer.Status() == http.StatusUnauthorized {
			c.Redirect(http.StatusTemporaryRedirect, location)
			c.Abort()
			return
		}
	}
}
