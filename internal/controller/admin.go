package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/controller/request"
	"go-web-template/internal/controller/response"
	"go-web-template/internal/middleware"
	"go-web-template/internal/model"
	"go-web-template/internal/service"
)

type AdminServer struct {
	userService service.User
	jwtManager  *middleware.JWTManager
}

func NewAdminServer(userService service.User, jwtManager *middleware.JWTManager) *AdminServer {
	return &AdminServer{
		userService: userService,
		jwtManager:  jwtManager,
	}
}

// 注册管理路由
func (a *AdminServer) RegisterRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin")
	{
		admin.POST("/login", a.Login)
		admin.POST("/add_token", a.AddToken)
		admin.POST("/change_password", a.ChangePassword)
	}
}

// 管理员登录
func (a *AdminServer) Login(c *gin.Context) {
	var req request.AdminLoginRequest
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

	le.Debug("admin login called")

	// 验证管理员凭据
	ok, err := a.userService.VerifyAdmin(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		le.Error("admin verification error")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    500,
				Message: "验证失败",
			},
		})
		return
	}

	if !ok {
		le.Warn("admin login failed - invalid credentials")
		c.JSON(http.StatusOK, response.AdminLoginResponse{
			BaseResponse: response.BaseResponse{
				Code:    401,
				Message: "用户名或密码错误",
			},
		})
		return
	}

	// 获取管理员用户信息
	adminUser, err := a.userService.VerifyUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		le.Error("failed to get admin user")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    500,
				Message: "获取用户信息失败",
			},
		})
		return
	}

	// 生成 JWT token 并设置为 cookie
	token, err := a.jwtManager.GenerateToken(fmt.Sprintf("%d", adminUser.ID), adminUser.Username, adminUser.Role)
	if err != nil {
		le.Error("failed to generate admin token")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    500,
				Message: "生成令牌失败",
			},
		})
		return
	}

	// 设置 JWT token 为 cookie
	c.SetCookie("jwt_token", token, 900, "/", "", false, true)

	le.Info("admin login successful")
	c.JSON(http.StatusOK, response.AdminLoginResponse{
		BaseResponse: response.BaseResponse{
			Code:    0,
			Message: "登录成功",
		},
	})
}

// 生成访问令牌
func (a *AdminServer) AddToken(c *gin.Context) {
	var req request.AdminAddTokenRequest
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
		"userid": req.UserID,
		"expire": req.Expire,
	})

	le.Debug("add token called")

	// 生成 JWT 访问令牌
	accessToken, err := a.jwtManager.GenerateToken(strconv.Itoa(req.UserID), "user_"+strconv.Itoa(req.UserID), "user")
	if err != nil {
		le.Error("failed to generate access token")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    500,
				Message: "生成访问令牌失败",
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.AdminAddTokenResponse{
		BaseResponse: response.BaseResponse{
			Code:    0,
			Message: "令牌生成成功",
		},
		AccessToken: accessToken,
	})
}

// 管理员修改密码
func (a *AdminServer) ChangePassword(c *gin.Context) {
	var req request.AdminChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    400,
				Message: "请求参数错误",
			},
		})
		return
	}

	le := logger.FromContext(c.Request.Context())
	le.Debug("admin change password called")

	// 从 JWT token 中获取用户ID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    401,
				Message: "未授权访问",
			},
		})
		return
	}

	userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
	if err != nil {
		le.Error("invalid user ID")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    500,
				Message: "用户ID无效",
			},
		})
		return
	}

	// 更新密码
	err = a.userService.UpdatePassword(c.Request.Context(), userID, &model.UpdatePasswordRequest{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		le.Error("failed to update password")
		c.JSON(http.StatusOK, response.AdminChangePasswordResponse{
			BaseResponse: response.BaseResponse{
				Code:    400,
				Message: "密码修改失败: " + err.Error(),
			},
		})
		return
	}

	le.Info("admin password changed successfully")
	c.JSON(http.StatusOK, response.AdminChangePasswordResponse{
		BaseResponse: response.BaseResponse{
			Code:    0,
			Message: "密码修改成功",
		},
	})
}
