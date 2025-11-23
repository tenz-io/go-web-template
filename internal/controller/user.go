package controller

import (
	"go-web-template/internal/controller/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/config"
	"go-web-template/internal/constant"
	"go-web-template/internal/controller/request"
	"go-web-template/internal/controller/response"
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
func (uc *UserController) RegisterRoutes(r *gin.RouterGroup) {
	// 页面
	r.GET("/home", uc.home)

	// API
	r.POST("/generate_token", uc.generateToken)
}

func (uc *UserController) home(c *gin.Context) {
	le := logger.FromContext(c.Request.Context())
	userID, _, err := middleware.GetUserInfoFromContext(c)
	if err != nil {
		le.WithError(err).Error("failed to get user info from context")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	userModel, err := uc.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		le.WithError(err).Error("failed to get user model")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	role := constant.Role(userModel.Role)

	c.HTML(http.StatusOK, "home.html", gin.H{
		"appName":     uc.appName,
		"name":        uc.appName,
		"username":    userModel.Username,
		"displayName": userModel.Username,
		"role":        role.String(),
		"isAdmin":     role == constant.RoleAdmin,
	})
}

// generateToken 用户生成 API Token（JWT 格式）
func (uc *UserController) generateToken(c *gin.Context) {
	var (
		le  = logger.FromContext(c.Request.Context())
		req request.UserGenerateTokenRequest
	)
	if err := c.ShouldBindJSON(&req); err != nil {
		le.Warn("invalid request params")
		response.FailWithJson(c, 400, "请求参数错误")
		return
	}

	le.Debug("user generate api token")

	userID, userRole, err := middleware.GetUserInfoFromContext(c)
	if err != nil {
		le.WithError(err).Error("failed to get user info from context")
		response.FailWithJson(c, 401, "未授权访问")
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

	token, err := uc.jwtManager.GenerateTokenWithExpire(userID, userRole, expDuration)
	if err != nil {
		le.WithError(err).Error("failed to generate token")
		response.FailWithJson(c, 500, "生成令牌失败")
		return
	}

	le.Info("user token generated successfully")
	response.OkWithJson(c, response.UserGenerateTokenResponseBody{
		Token: token,
	})
}
