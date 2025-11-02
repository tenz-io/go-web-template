package controller

import (
	"fmt"
	"go-web-template/internal/constant"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/config"
	"go-web-template/internal/middleware"
	"go-web-template/internal/service"
)

type WebServer struct {
	engine     *gin.Engine
	cfg        *config.Config
	api        *ApiServer
	admin      *AdminServer
	auth       *AuthServer
	user       *UserServer
	jwtManager *middleware.JWTManager
	tokenSvc   service.Token
}

func NewWebServer(cfg *config.Config, apiServer *ApiServer, adminServer *AdminServer, authServer *AuthServer, userServer *UserServer, jwtManager *middleware.JWTManager, tokenSvc service.Token) *WebServer {
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
		user:       userServer,
		jwtManager: jwtManager,
		tokenSvc:   tokenSvc,
	}

	return ws
}

func (ws *WebServer) Init() error {
	// 设置全局中间件
	ws.engine.Use(middleware.Logger())
	ws.engine.Use(gin.Recovery())
	ws.engine.Use(middleware.CORS())

	ws.engine.LoadHTMLGlob(ws.cfg.App.Web + "/*.html")
	ws.engine.Static("/static", ws.cfg.App.Web+"/static")

	// 注册路由
	ws.registerRoutes()

	return nil
}

func (ws *WebServer) registerRoutes() {
	// 注册统一认证路由（最先注册，避免冲突）
	authGroup := ws.engine.Group("")
	ws.auth.RegisterRoutes(authGroup)

	// 注册 API 路由（可选鉴权）
	apiGroup := ws.engine.Group("api")
	apiGroup.Use(middleware.Auth(middleware.AuthConfig{
		Type:     middleware.AuthTypeBearer,
		Required: true,
		Role:     constant.RoleUser,
	}, ws.jwtManager, ws.tokenSvc))
	ws.api.RegisterRoutes(apiGroup)

	// 注册管理路由（需要管理员权限）
	adminGroup := ws.engine.Group("admin")
	adminGroup.Use(middleware.Auth(middleware.AuthConfig{
		Type:     middleware.AuthTypeCookie,
		Required: true,
		Role:     constant.RoleAdmin,
	}, ws.jwtManager, ws.tokenSvc))
	ws.admin.RegisterRoutes(adminGroup)

	// 注册用户路由（需要用户权限）
	userGroup := ws.engine.Group("user")
	userGroup.Use(middleware.Auth(middleware.AuthConfig{
		Type:     middleware.AuthTypeCookie,
		Required: true,
		Role:     constant.RoleUser,
	}, ws.jwtManager, ws.tokenSvc))
	ws.user.RegisterRoutes(userGroup)

	// 页面路由（最后注册，避免被API路由覆盖）
	// 首页 - 不需要认证
	ws.engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// 登录页面 - 不需要认证
	ws.engine.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
}

func (ws *WebServer) Run(errC chan<- error) {
	addr := fmt.Sprintf(":%d", ws.cfg.App.Port)
	logger.WithField("addr", addr).Info("HTTP server starting")

	err := ws.engine.Run(addr)
	if err != nil {
		errC <- err
	}
}
