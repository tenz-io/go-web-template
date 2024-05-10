package controller

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc/metadata"

	"github.com/tenz-io/gokit/ginext/errcode"
	"github.com/tenz-io/gokit/logger"

	pbapp "go-web-template/api/http/app"
	"go-web-template/internal/service"
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
		meta, existing = metadata.FromIncomingContext(ctx)
		le             = logger.FromContext(ctx).WithFields(logger.Fields{
			"username":     request.GetUsername(),
			"meta":         meta,
			"existing":     existing,
			"Content-Type": meta.Get("Content-Type"),
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
		meta, existing = metadata.FromIncomingContext(ctx)
		le             = logger.FromContext(ctx).WithFields(logger.Fields{
			"meta":      meta,
			"existing":  existing,
			"image_key": request.GetKey(),
		})
	)

	le.Debug("get image")
	// TODO
	return nil, errcode.NotFound(http.StatusOK, "image not found")
}

func (as *ApiServer) UploadImage(ctx context.Context, request *pbapp.UploadImageRequest) (*pbapp.UploadImageResponse, error) {
	var (
		meta, existing = metadata.FromIncomingContext(ctx)
		le             = logger.FromContext(ctx).WithFields(logger.Fields{
			"meta":      meta,
			"existing":  existing,
			"filename":  request.GetFilename(),
			"file_size": len(request.GetFile()),
		})
	)

	le.Debug("upload image")
	// TODO
	return nil, errcode.InternalServer(http.StatusOK, fmt.Sprintf("upload image error"))
}
