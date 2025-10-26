package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/controller/request"
	"go-web-template/internal/controller/response"
	"go-web-template/internal/middleware"
	"go-web-template/internal/service"
)

type UserServer struct {
	userService service.User
	jwtManager  *middleware.JWTManager
}

func NewUserServer(userService service.User, jwtManager *middleware.JWTManager) *UserServer {
	return &UserServer{
		userService: userService,
		jwtManager:  jwtManager,
	}
}

// 注册用户路由
func (u *UserServer) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/generate_token", u.GenerateToken)
}

// GenerateToken 用户生成API token
func (u *UserServer) GenerateToken(c *gin.Context) {
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
	le.Debug("user generate token called")

	// 从 JWT token 中获取用户ID
	userID, _, err := middleware.GetUserInfoFromContext(c)
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

	// 生成 JWT token
	token, err := u.jwtManager.GenerateToken(userID, 0) // 普通用户角色
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

	le.Info("user token generated successfully")
	c.JSON(http.StatusOK, response.UserGenerateTokenResponse{
		BaseResponse: response.BaseResponse{
			Code:    0,
			Message: "令牌生成成功",
		},
		Token: token,
	})
}
