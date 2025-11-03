package setup

import (
	"go-web-template/internal/config"
	"go-web-template/internal/controller"
	"go-web-template/internal/database"
	"go-web-template/internal/job"
	"go-web-template/internal/middleware"
	repodao "go-web-template/internal/repository/dao"
	"go-web-template/internal/service"

	"github.com/google/wire"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// 简化的依赖注入提供者函数

var RepoProviderSet = wire.NewSet(
	repodao.NewUserDao,
)

var ServiceProviderSet = wire.NewSet(
	service.NewUserService,
)

var ControllerProviderSet = wire.NewSet(
	controller.NewAuthController,
	controller.NewUserController,
	controller.NewAdminController,
	controller.NewApiController,
	controller.NewWebServer,
)

var JobProviderSet = wire.NewSet(
	job.NewHealthReporter,
	job.NewManager,
)

var ComponentProviderSet = wire.NewSet(
	ProvideDB,
	ProvideJWTManager,
	ProvideCron,
)

func ProvideDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := database.NewDB(&cfg.DB)
	if err != nil {
		return nil, err
	}
	return db.GetConn(), nil
}

func ProvideJWTManager(cfg *config.Config) (*middleware.JWTManager, error) {
	jwt := middleware.NewJWTManager(&cfg.JWT)
	return jwt, nil
}

func ProvideCron() *cron.Cron {
	return cron.New()
}
