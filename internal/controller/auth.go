package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/constant"
	"go-web-template/internal/controller/request"
	"go-web-template/internal/controller/response"
	"go-web-template/internal/middleware"
	"go-web-template/internal/repository"
)

// AuthServer 统一认证服务器
type AuthServer struct {
	userRepo   repository.User
	jwtManager *middleware.JWTManager
}

// NewAuthServer 创建统一认证服务器
func NewAuthServer(userRepo repository.User, jwtManager *middleware.JWTManager) *AuthServer {
	return &AuthServer{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// RegisterRoutes 注册认证路由
func (a *AuthServer) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", a.Login)
		auth.POST("/logout", a.Logout)
	}
}

// Login 统一登录接口
func (a *AuthServer) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    400,
				Message: "请求参数错误",
			},
		})
		return
	}

	le := logger.FromContext(c.Request.Context()).WithFields(logger.Fields{
		"username": req.Username,
	})

	le.Debug("auth login")

	// 验证用户凭据
	user, err := a.userRepo.VerifyUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		le.Warn("user login failed")
		c.JSON(http.StatusOK, response.LoginResponse{
			BaseResponse: response.BaseResponse{
				Code:    401,
				Message: "用户名或密码错误",
			},
		})
		return
	}

	// 生成 JWT token
	token, err := a.jwtManager.GenerateToken(fmt.Sprintf("%d", user.ID), user.Username, user.Role)
	if err != nil {
		le.Error("failed to generate token")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    500,
				Message: "生成令牌失败",
			},
		})
		return
	}

	// 根据用户角色设置重定向地址
	var redirect string
	if user.Role == string(constant.RoleAdmin) {
		redirect = "/admin/"
	} else {
		redirect = "/"
	}

	// 设置 Cookie（用于管理后台）
	if user.Role == string(constant.RoleAdmin) {
		c.SetCookie("jwt_token", token, 3600*24, "/", "", false, true)
	}

	le.Info("user login successful")
	c.JSON(http.StatusOK, response.LoginResponse{
		BaseResponse: response.BaseResponse{
			Code:    0,
			Message: "登录成功",
		},
		Token:    token,
		Role:     user.Role,
		Redirect: redirect,
	})
}

// Logout 登出接口
func (a *AuthServer) Logout(c *gin.Context) {
	// 清除 Cookie
	c.SetCookie("jwt_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, response.BaseResponse{
		Code:    0,
		Message: "登出成功",
	})
}
