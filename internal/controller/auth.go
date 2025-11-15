package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/constant"
	"go-web-template/internal/controller/request"
	"go-web-template/internal/controller/response"
	"go-web-template/internal/middleware"
	"go-web-template/internal/repository/dao"
	"go-web-template/internal/service"
)

// AuthServer 统一认证服务器
type AuthController struct {
	userDao     dao.User
	userService service.User
	jwtManager  *middleware.JWTManager
}

// NewAuthController 创建统一认证服务器
func NewAuthController(userDao dao.User, userService service.User, jwtManager *middleware.JWTManager) *AuthController {
	return &AuthController{
		userDao:     userDao,
		userService: userService,
		jwtManager:  jwtManager,
	}
}

// RegisterRoutes 注册认证路由
func (ac *AuthController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", ac.Login)
	rg.POST("/logout", ac.Logout)

	protected := rg.Group("/auth")
	protected.Use(middleware.Auth(middleware.AuthConfig{
		Type:     middleware.AuthTypeCookie,
		Required: true,
		Role:     constant.RoleUser,
	}, ac.jwtManager))
	protected.POST("/change_password", ac.ChangePassword)
}

// Login 统一登录接口
func (ac *AuthController) Login(c *gin.Context) {
	var (
		le  = logger.FromContext(c.Request.Context())
		req request.LoginRequest
	)

	// 只支持JSON提交
	if err := c.ShouldBindJSON(&req); err != nil {
		le.WithError(err).Error("invalid request")
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    400,
				Message: "请求参数错误",
			},
		})
		return
	}

	le = logger.FromContext(c.Request.Context()).WithFields(logger.Fields{
		"username": req.Username,
	})

	le.Debug("auth login")

	// 验证用户凭据
	verifyParam := service.VerifyUserParam{
		Username: req.Username,
		Password: req.Password,
	}
	userModel, err := ac.userService.VerifyUser(c.Request.Context(), verifyParam)
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
	token, err := ac.jwtManager.GenerateToken(userModel.ID, userModel.Role)
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

	// 设置Cookie
	c.SetCookie(middleware.JWTTokenCookieName, token, 3600*24, "/", "", false, true)

	// 根据角色重定向到不同页面
	redirect := "/user/home"
	if userModel.Role == int32(constant.RoleAdmin) {
		redirect = "/admin/home"
	}

	le.WithFields(logger.Fields{
		"user_id":  userModel.ID,
		"role":     constant.Role(userModel.Role).String(),
		"redirect": redirect,
	}).Info("user login successful")

	// 返回JSON响应
	c.JSON(http.StatusOK, response.LoginResponse{
		BaseResponse: response.BaseResponse{
			Code:    0,
			Message: "登录成功",
		},
		Role:     constant.Role(userModel.Role).String(),
		Redirect: redirect,
	})
}

// Logout 登出接口
func (ac *AuthController) Logout(c *gin.Context) {
	// 清除 Cookie
	c.SetCookie(middleware.JWTTokenCookieName, "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, response.BaseResponse{
		Code:    0,
		Message: "登出成功",
	})
}

// ChangePassword 修改密码（用户/管理员通用）
func (ac *AuthController) ChangePassword(c *gin.Context) {
	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    400,
				Message: "请求参数错误",
			},
		})
		return
	}

	userID, _, err := middleware.GetUserInfoFromContext(c)
	if err != nil {
		logger.FromContext(c.Request.Context()).WithError(err).Warn("failed to get user info from context")
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    401,
				Message: "未授权访问",
			},
		})
		return
	}

	updateParam := service.UpdatePasswordParam{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	if err := ac.userService.UpdatePassword(c.Request.Context(), userID, updateParam); err != nil {
		logger.FromContext(c.Request.Context()).WithError(err).Error("failed to update password")
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    400,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		BaseResponse: response.BaseResponse{
			Code:    0,
			Message: "密码修改成功",
		},
	})
}
