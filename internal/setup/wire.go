//go:build wireinject

package setup

import (
	"github.com/google/wire"

	"go-web-template/internal/config"
)

func InitializeControllers(
	_ *config.Config,
) (*Controllers, error) {
	wire.Build(
		wire.Struct(new(Controllers), "*"),
		RepoProviderSet,
		ServiceProviderSet,
		ControllerProviderSet,
	)

	return nil, nil

}
