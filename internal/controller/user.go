package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/config"
	"go-web-template/internal/constant"
	"go-web-template/internal/controller/request"
	"go-web-template/internal/controller/response"
	"go-web-template/internal/middleware"
	"go-web-template/internal/service"
)

type UserController struct {
	userService service.User
	jwtManager  *middleware.JWTManager
	appName     string
}

func NewUserController(cfg *config.Config, userService service.User, jwtManager *middleware.JWTManager) *UserController {
	return &UserController{
		userService: userService,
		jwtManager:  jwtManager,
		appName:     cfg.App.Name,
	}
}

// RegisterRoutes 注册用户侧路由
func (u *UserController) RegisterRoutes(r *gin.RouterGroup) {
	// 页面
	r.GET("/home", u.home)

	// API
	r.POST("/generate_token", u.generateToken)
}

func (u *UserController) home(c *gin.Context) {
	le := logger.FromContext(c.Request.Context())
	userID, _, err := middleware.GetUserInfoFromContext(c)
	if err != nil {
		le.WithError(err).Error("failed to get user info from context")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	userModel, err := u.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		le.WithError(err).Error("failed to get user model")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	role := constant.Role(userModel.Role)

	c.HTML(http.StatusOK, "home.html", gin.H{
		"appName":     u.appName,
		"name":        u.appName,
		"username":    userModel.Username,
		"displayName": userModel.Username,
		"role":        role.String(),
		"isAdmin":     role == constant.RoleAdmin,
	})
}

// generateToken 用户生成 API Token（JWT 格式）
func (u *UserController) generateToken(c *gin.Context) {
	var req request.UserGenerateTokenRequest
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
	le.Debug("user generate api token")

	userID, userRole, err := middleware.GetUserInfoFromContext(c)
	if err != nil {
		le.WithError(err).Error("failed to get user info from context")
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    401,
				Message: "未授权访问",
			},
		})
		return
	}

	expireHours := req.ExpireHours
	if expireHours <= 0 {
		expireHours = 1
	}

	if expireHours > 24*365*5 {
		expireHours = 24 * 365 * 5
	}

	expDuration := time.Duration(expireHours) * time.Hour

	token, err := u.jwtManager.GenerateTokenWithExpire(userID, userRole, expDuration)
	if err != nil {
		le.WithError(err).Error("failed to generate token")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    500,
				Message: "生成令牌失败",
			},
		})
		return
	}

	le.Info("user token generated successfully")
	c.JSON(http.StatusOK, response.UserGenerateTokenResponse{
		BaseResponse: response.BaseResponse{
			Code:    0,
			Message: "令牌生成成功",
		},
		Token: token,
	})
}
