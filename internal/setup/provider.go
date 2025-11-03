package setup

import (
	"go-web-template/internal/config"
	"go-web-template/internal/controller"
	"go-web-template/internal/database"
	"go-web-template/internal/middleware"
	"go-web-template/internal/repository/dao"
	"go-web-template/internal/service"
)

// 简化的依赖注入提供者函数

// 数据库提供者
func NewDatabase(cfg *config.Config) (*database.DB, error) {
	return database.NewDB(&cfg.DB)
}

// JWT 管理器提供者
func NewJWTManager(cfg *config.Config) *middleware.JWTManager {
	return middleware.NewJWTManager(&cfg.JWT)
}

// 仓库层提供者
func NewUserRepository(db *database.DB) dao.User {
	return dao.NewUser(db.GetConn())
}

// 服务层提供者
func NewUserService(cfg *config.Config, userRepo dao.User) service.User {
	return service.NewUser(cfg, userRepo)
}

// 控制器层提供者
func NewApiController(userRepo dao.User, jwtManager *middleware.JWTManager) *controller.ApiServer {
	return controller.NewApiServer(userRepo, jwtManager)
}

func NewAdminController(userService service.User, jwtManager *middleware.JWTManager) *controller.AdminServer {
	return controller.NewAdminServer(userService, jwtManager)
}

func NewAuthController(userRepo dao.User, userService service.User, jwtManager *middleware.JWTManager) *controller.AuthServer {
	return controller.NewAuthServer(userRepo, userService, jwtManager)
}

func NewUserController(userService service.User, jwtManager *middleware.JWTManager) *controller.UserServer {
	return controller.NewUserServer(userService, jwtManager)
}

func NewWebController(cfg *config.Config, apiController *controller.ApiServer, adminController *controller.AdminServer, authController *controller.AuthServer, userController *controller.UserServer, jwtManager *middleware.JWTManager) *controller.WebServer {
	return controller.NewWebServer(cfg, apiController, adminController, authController, userController, jwtManager)
}

// 主控制器提供者
func NewControllers(webController *controller.WebServer) *Controllers {
	return &Controllers{
		webServer: webController,
	}
}
