package controller

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/tenz-io/gokit/ginext/errcode"
	"github.com/tenz-io/gokit/ginext/metadata"
	"github.com/tenz-io/gokit/logger"

	pbapp "go-web-template/api/http/app"
	"go-web-template/internal/service"
)

var (
	_ pbapp.ApiServerHTTPServer = (*ApiServer)(nil)
)

type ApiServer struct {
	userService service.User
}

func NewApiServer(
	userService service.User,
) *ApiServer {
	return &ApiServer{
		userService: userService,
	}
}

func (as *ApiServer) Hello(ctx context.Context, request *pbapp.HelloRequest) (*pbapp.HelloResponse, error) {
	var (
		le = logger.FromContext(ctx)
	)
	defer func() {
		le.Debug("hello called")
	}()

	user, err := as.userService.GetByName(ctx, request.GetName())
	if err != nil {
		return nil, err
	}
	return &pbapp.HelloResponse{
		Message: user.Profile,
	}, nil
}

func (as *ApiServer) Login(ctx context.Context, request *pbapp.LoginRequest) (*pbapp.LoginResponse, error) {
	var (
		meta, existing = metadata.FromContext(ctx)
		le             = logger.FromContext(ctx).WithFields(logger.Fields{
			"username": request.GetUsername(),
			"meta":     meta,
			"existing": existing,
		})
	)

	defer func() {
		le.Debug("login called")
	}()

	_, err := as.userService.GetByName(ctx, request.GetUsername())
	if err != nil {
		return nil, errcode.NotFound(http.StatusOK, "user not found")
	}

	if request.Username == request.Password {
		return &pbapp.LoginResponse{}, nil
	}

	return nil, errcode.Unauthorized(http.StatusOK, "password incorrect")
}

func (as *ApiServer) GetImage(ctx context.Context, request *pbapp.GetImageRequest) (*pbapp.GetImageResponse, error) {
	var (
		meta, existing = metadata.FromContext(ctx)
		le             = logger.FromContext(ctx).WithFields(logger.Fields{
			"meta":     meta,
			"existing": existing,
			"key":      request.GetKey(),
		})
	)

	le.Debug("get image")
	// TODO
	return nil, errcode.NotFound(http.StatusOK, "image not found")
}

func (as *ApiServer) UploadImage(ctx context.Context, request *pbapp.UploadImageRequest) (*pbapp.UploadImageResponse, error) {
	var (
		meta, existing = metadata.FromContext(ctx)
		le             = logger.FromContext(ctx).WithFields(logger.Fields{
			"meta":      meta,
			"existing":  existing,
			"file_size": len(request.GetFile()),
		})
		key = request.GetKey()
	)

	if key == "" {
		key = fmt.Sprintf("%x", md5.Sum(request.GetFile()))
	}

	le = le.WithFields(logger.Fields{
		"key": key,
	})

	le.Debug("upload image")
	return &pbapp.UploadImageResponse{
		Key: key,
	}, nil
}
