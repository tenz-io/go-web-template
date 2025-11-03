package controller

import (
	"go-web-template/internal/repository/dao"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/controller/request"
	"go-web-template/internal/controller/response"
	"go-web-template/internal/middleware"
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
func (as *ApiController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/hello", as.Hello)
}

// Hello 接口
func (as *ApiController) Hello(c *gin.Context) {
	var req request.HelloRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    400,
				Message: "请求参数错误",
			},
		})
		return
	}

	le := logger.FromContext(c.Request.Context())
	le.Debug("hello called")

	user, err := as.userRepo.GetByName(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    500,
				Message: "获取用户信息失败",
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.HelloResponse{
		BaseResponse: response.BaseResponse{
			Code:    0,
			Message: user.Profile,
		},
	})
}
