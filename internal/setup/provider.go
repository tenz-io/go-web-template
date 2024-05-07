package setup

import (
	"github.com/google/wire"

	"go-web-template/internal/controller"
	"go-web-template/internal/repository"
	"go-web-template/internal/service"
)

var RepoProviderSet = wire.NewSet(
	repository.NewUser,
)

var ServiceProviderSet = wire.NewSet(
	service.NewUser,
)

var ControllerProviderSet = wire.NewSet(
	controller.NewApiServer,
	controller.NewWebServer,
)
