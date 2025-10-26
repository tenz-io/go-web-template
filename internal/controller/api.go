package controller

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/controller/request"
	"go-web-template/internal/controller/response"
	"go-web-template/internal/middleware"
	"go-web-template/internal/repository"
)

type ApiServer struct {
	userRepo   repository.User
	jwtManager *middleware.JWTManager
}

func NewApiServer(userRepo repository.User, jwtManager *middleware.JWTManager) *ApiServer {
	return &ApiServer{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// 注册 API 路由
func (as *ApiServer) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/login", as.Login)
	r.GET("/hello", as.Hello)
	r.GET("/image/:key", as.GetImage)
	r.POST("/upload", as.UploadImage)
}

// 用户登录
func (as *ApiServer) Login(c *gin.Context) {
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

	le.Debug("api login")

	// 验证用户凭据
	user, err := as.userRepo.VerifyUser(c.Request.Context(), req.Username, req.Password)
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
	token, err := as.jwtManager.GenerateToken(user.ID, user.Role)
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

	c.JSON(http.StatusOK, response.LoginResponse{
		BaseResponse: response.BaseResponse{
			Code:    0,
			Message: "登录成功",
		},
		Token: token,
	})
}

// Hello 接口
func (as *ApiServer) Hello(c *gin.Context) {
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

// 获取图片
func (as *ApiServer) GetImage(c *gin.Context) {
	key := c.Param("key")
	le := logger.FromContext(c.Request.Context()).WithFields(logger.Fields{
		"key": key,
	})

	le.Debug("get image")
	// TODO: 实现图片获取逻辑
	c.JSON(http.StatusNotFound, response.ErrorResponse{
		BaseResponse: response.BaseResponse{
			Code:    404,
			Message: "图片不存在",
		},
	})
}

// 上传图片
func (as *ApiServer) UploadImage(c *gin.Context) {
	key := c.PostForm("key")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    400,
				Message: "文件上传失败",
			},
		})
		return
	}

	le := logger.FromContext(c.Request.Context()).WithFields(logger.Fields{
		"file_size": file.Size,
	})

	if key == "" {
		// 生成文件内容的 MD5 作为 key
		fileContent, _ := file.Open()
		defer fileContent.Close()
		hash := md5.New()
		hash.Write([]byte(file.Filename + strconv.FormatInt(file.Size, 10)))
		key = fmt.Sprintf("%x", hash.Sum(nil))
	}

	le = le.WithFields(logger.Fields{
		"key": key,
	})

	le.Debug("upload image")
	c.JSON(http.StatusOK, response.UploadResponse{
		BaseResponse: response.BaseResponse{
			Code:    0,
			Message: "上传成功",
		},
		Key: key,
	})
}
