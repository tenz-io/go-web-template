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
		wire.Struct(new(Controllers), "*"),
		ComponentProviderSet,
		RepoProviderSet,
		ServiceProviderSet,
		ControllerProviderSet,
		JobProviderSet,
	)

	return nil, nil
}
