//go:build wireinject

package setup

import (
	"github.com/google/wire"

	"go-web-template/internal/config"
)

func InitializeControllers(
	cfg *config.Config,
) (*Controllers, error) {
	wire.Build(
		// 数据库
		NewDatabase,
		// JWT 管理器
		NewJWTManager,
		// 仓库层
		NewUserRepository,
		NewTokenRepository,
		// 服务层
		NewUserService,
		NewTokenService,
		// 控制器层
		NewApiController,
		NewAdminController,
		NewAuthController,
		NewUserController,
		NewWebController,
		// 主控制器
		NewControllers,
	)

	return nil, nil
}
