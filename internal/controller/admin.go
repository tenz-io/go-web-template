package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/constant"
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
		admin.GET("/users", a.GetUsers)
		admin.POST("/add_user", a.AddUser)
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
	token, err := a.jwtManager.GenerateToken(adminUser.ID, adminUser.Role)
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
	c.SetCookie(middleware.JWTTokenCookieName, token, 900, "/", "", false, true)

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
	accessToken, err := a.jwtManager.GenerateToken(int64(req.UserID), int32(constant.RoleUser))
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
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    401,
				Message: "未授权访问",
			},
		})
		return
	}

	// 处理不同类型的 user_id
	var userID int64
	var err error
	switch v := userIDInterface.(type) {
	case int64:
		userID = v
	default:
		le.Error("invalid user ID type")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    500,
				Message: "用户ID类型无效",
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

// GetUsers 获取用户列表
func (a *AdminServer) GetUsers(c *gin.Context) {
	le := logger.FromContext(c.Request.Context())
	le.Debug("admin get users called")

	// 获取分页参数
	limit := 100 // 默认限制
	offset := 0

	// 获取用户列表
	users, total, err := a.userService.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		le.Error("failed to get users")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    500,
				Message: "获取用户列表失败",
			},
		})
		return
	}

	// 转换用户角色为字符串
	var userList []gin.H
	for _, user := range users {
		userList = append(userList, gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"role":       constant.Role(user.Role).String(),
			"email":      user.Email,
			"created_at": user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		BaseResponse: response.BaseResponse{
			Code:    0,
			Message: "获取用户列表成功",
		},
		Data: gin.H{
			"users": userList,
			"total": total,
		},
	})
}

// AddUser 添加用户
func (a *AdminServer) AddUser(c *gin.Context) {
	var req request.AdminAddUserRequest
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
	le.Debug("admin add user called")

	// 将角色字符串转换为数字
	role, err := constant.ParseRole(req.Role)
	if err != nil {
		le.Error("invalid role")
		c.JSON(http.StatusOK, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    400,
				Message: "无效的用户角色",
			},
		})
		return
	}

	// 创建用户
	user, err := a.userService.CreateUser(c.Request.Context(), &model.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
		Role:     int32(role),
		Email:    req.Email,
	})
	if err != nil {
		le.Error("failed to create user")
		c.JSON(http.StatusOK, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    400,
				Message: "创建用户失败: " + err.Error(),
			},
		})
		return
	}

	le.Info("user created successfully")
	c.JSON(http.StatusOK, response.SuccessResponse{
		BaseResponse: response.BaseResponse{
			Code:    0,
			Message: "用户创建成功",
		},
		Data: user,
	})
}
