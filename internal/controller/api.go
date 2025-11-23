package controller

import (
	"go-web-template/internal/controller/middleware"
	"go-web-template/internal/repository/dao"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/controller/request"
	"go-web-template/internal/controller/response"
)

type ApiController struct {
	userRepo   dao.User
	jwtManager *middleware.JWTManager
}

func NewApiController(userRepo dao.User, jwtManager *middleware.JWTManager) *ApiController {
	return &ApiController{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// 注册 API 路由
func (ac *ApiController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/hello", ac.Hello)
}

// Hello 接口
func (ac *ApiController) Hello(c *gin.Context) {
	var req request.HelloRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithJson(c, 400, "请求参数错误")
		return
	}

	le := logger.FromContext(c.Request.Context())
	le.Debug("hello called")

	user, err := ac.userRepo.GetByName(c.Request.Context(), req.Name)
	if err != nil {
		response.FailWithJson(c, 500, "获取用户信息失败")
		return
	}

	response.OkWithJson(c, user.Profile)
}
