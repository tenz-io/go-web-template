package controller

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"time"

	"github.com/tenz-io/gokit/ginext"
	"github.com/tenz-io/gokit/ginext/errcode"
	"github.com/tenz-io/gokit/ginext/metadata"
	"github.com/tenz-io/gokit/logger"

	pbapp "go-web-template/api/http/app"
	"go-web-template/internal/service"
)

var (
	_ pbapp.AdminServerHTTPServer = (*AdminServer)(nil)
)

type AdminServer struct {
	userService service.User
}

func NewAdminServer(
	userService service.User,
) *AdminServer {
	return &AdminServer{
		userService: userService,
	}
}

func (a *AdminServer) Login(ctx context.Context, req *pbapp.AdminLoginRequest) (*pbapp.AdminLoginResponse, error) {
	var (
		meta = metadata.SafeFromContext(ctx)
		le   = logger.FromContext(ctx).WithFields(logger.Fields{
			"username": req.GetUsername(),
			"meta":     meta,
		})
	)

	defer func() {
		le.Debug("login called")
	}()

	ok, _ := a.userService.VerifyAdmin(ctx, req.GetUsername(), req.GetPassword())
	if !ok {
		return nil, errcode.Unauthorized(http.StatusOK, "auth failed")
	}

	expiredAt := time.Now().Add(15 * time.Minute)
	accessToken, err := ginext.GenerateToken(1, ginext.RoleAdmin, ginext.TokenTypeAccess, expiredAt)
	if err != nil {
		return nil, errcode.InternalServer(http.StatusInternalServerError, "failed to generate token")
	}

	// set access token to cookie
	meta.Ctx.SetCookie(ginext.CookieTokenName, accessToken, 900, "/", "", false, true)

	return &pbapp.AdminLoginResponse{}, nil
}

func (a *AdminServer) AddToken(ctx context.Context, request *pbapp.AdminAddTokenRequest) (*pbapp.AdminAddTokenResponse, error) {
	// generate access token
	expiresAt := time.Now().Add(time.Duration(request.GetExpire()) * time.Hour * 24)
	accessToken, err := ginext.GenerateToken(request.GetUserid(), ginext.RoleUser, ginext.TokenTypeAccess, expiresAt)
	if err != nil {
		return nil, errcode.InternalServer(http.StatusInternalServerError, "failed to generate token")
	}

	return &pbapp.AdminAddTokenResponse{
		AccessToken: accessToken,
	}, nil

}

func (a *AdminServer) UploadImage(ctx context.Context, request *pbapp.AdminUploadImageRequest) (*pbapp.AdminUploadImageResponse, error) {
	var (
		meta = metadata.SafeFromContext(ctx)
		le   = logger.FromContext(ctx).WithFields(logger.Fields{
			"meta":      meta,
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
	return &pbapp.AdminUploadImageResponse{
		Key: key,
	}, nil
}
