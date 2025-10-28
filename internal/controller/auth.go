package controller

import (
	"go-web-template/internal/repository/dao"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/constant"
	"go-web-template/internal/controller/request"
	"go-web-template/internal/controller/response"
	"go-web-template/internal/middleware"
)

// AuthServer 统一认证服务器
type AuthServer struct {
	userDao    dao.User
	jwtManager *middleware.JWTManager
}

// NewAuthServer 创建统一认证服务器
func NewAuthServer(userDao dao.User, jwtManager *middleware.JWTManager) *AuthServer {
	return &AuthServer{
		userDao:    userDao,
		jwtManager: jwtManager,
	}
}

// RegisterRoutes 注册认证路由
func (a *AuthServer) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", a.Login)
	rg.POST("/logout", a.Logout)
}

// Login 统一登录接口
func (a *AuthServer) Login(c *gin.Context) {
	var req request.LoginRequest

	// 只支持JSON提交
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
	user, err := a.userDao.VerifyUser(c.Request.Context(), req.Username, req.Password)
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
	token, err := a.jwtManager.GenerateToken(user.ID, user.Role)
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
	if user.Role == int32(constant.RoleAdmin) {
		redirect = "/admin/home"
	}

	le.WithFields(logger.Fields{
		"user_id":  user.ID,
		"role":     constant.Role(user.Role).String(),
		"redirect": redirect,
	}).Info("user login successful")

	// 返回JSON响应
	c.JSON(http.StatusOK, response.LoginResponse{
		BaseResponse: response.BaseResponse{
			Code:    0,
			Message: "登录成功",
		},
		Role:     constant.Role(user.Role).String(),
		Redirect: redirect,
	})
}

// Logout 登出接口
func (a *AuthServer) Logout(c *gin.Context) {
	// 清除 Cookie
	c.SetCookie(middleware.JWTTokenCookieName, "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, response.BaseResponse{
		Code:    0,
		Message: "登出成功",
	})
}
