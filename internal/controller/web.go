package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/config"
	"go-web-template/internal/middleware"
)

type WebServer struct {
	engine     *gin.Engine
	cfg        *config.Config
	api        *ApiServer
	admin      *AdminServer
	auth       *AuthServer
	jwtManager *middleware.JWTManager
}

func NewWebServer(cfg *config.Config, apiServer *ApiServer, adminServer *AdminServer, authServer *AuthServer, jwtManager *middleware.JWTManager) *WebServer {
	if cfg.Verbose {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	ws := &WebServer{
		engine:     gin.New(),
		cfg:        cfg,
		api:        apiServer,
		admin:      adminServer,
		auth:       authServer,
		jwtManager: jwtManager,
	}

	return ws
}

func (ws *WebServer) Init() error {
	// 设置全局中间件
	ws.engine.Use(middleware.Logger())
	ws.engine.Use(gin.Recovery())
	ws.engine.Use(middleware.CORS())

	// 设置静态文件和模板
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

	// 注册路由
	ws.registerRoutes()

	return nil
}

func (ws *WebServer) registerRoutes() {
	// 首页
	ws.engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"name": ws.cfg.App.Name,
		})
	})

	// 统一登录页面
	ws.engine.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"name": ws.cfg.App.Name,
		})
	})

	// 管理后台页面（需要管理员权限）
	ws.engine.GET("/admin/", middleware.AdminAuth(ws.jwtManager), func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_index.html", gin.H{
			"name": ws.cfg.App.Name,
		})
	})

	// 注册统一认证路由
	authGroup := ws.engine.Group("")
	ws.auth.RegisterRoutes(authGroup)

	// 注册 API 路由（可选鉴权）
	apiGroup := ws.engine.Group("")
	apiGroup.Use(middleware.Auth(middleware.AuthConfig{
		Type:     middleware.AuthTypeBearer,
		Required: false,
	}, ws.jwtManager))
	ws.api.RegisterRoutes(apiGroup)

	// 注册管理路由（需要管理员权限）
	adminGroup := ws.engine.Group("")
	adminGroup.Use(middleware.Auth(middleware.AuthConfig{
		Type:     middleware.AuthTypeCookie,
		Required: true,
	}, ws.jwtManager))
	ws.admin.RegisterRoutes(adminGroup)
}

func (ws *WebServer) Run(errC chan<- error) {
	addr := fmt.Sprintf(":%d", ws.cfg.App.Port)
	logger.WithField("addr", addr).Info("HTTP server starting")

	err := ws.engine.Run(addr)
	if err != nil {
		errC <- err
	}
}
