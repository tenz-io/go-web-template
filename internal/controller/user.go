package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/constant"
	"go-web-template/internal/controller/request"
	"go-web-template/internal/controller/response"
	"go-web-template/internal/middleware"
	"go-web-template/internal/model"
	"go-web-template/internal/service"
)

type UserServer struct {
	userService service.User
	tokenSvc    service.Token
	jwtManager  *middleware.JWTManager
}

func NewUserServer(userService service.User, tokenService service.Token, jwtManager *middleware.JWTManager) *UserServer {
	return &UserServer{
		userService: userService,
		tokenSvc:    tokenService,
		jwtManager:  jwtManager,
	}
}

// 注册用户路由
func (u *UserServer) RegisterRoutes(r *gin.RouterGroup) {
	// 用户页面路由
	r.GET("/home", u.home)
	r.POST("/change_password", u.changePassword)
	r.GET("/api_tokens", u.listTokens)
	r.POST("/api_tokens", u.createToken)
	r.DELETE("/api_tokens/:id", u.deleteToken)
}

func (u *UserServer) home(c *gin.Context) {
	le := logger.FromContext(c.Request.Context())
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

	userModel, err := u.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		le.WithError(err).Error("failed to get user model")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    500,
				Message: "获取用户信息失败",
			},
		})
		return
	}

	c.HTML(http.StatusOK, "user_home.html", gin.H{
		"name": userModel.Username,
		"role": constant.Role(userModel.Role).String(),
	})
}

func (u *UserServer) changePassword(c *gin.Context) {
	var req request.UserChangePasswordRequest
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

	updateReq := &model.UpdatePasswordRequest{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	if err := u.userService.UpdatePassword(c.Request.Context(), userID, updateReq); err != nil {
		logger.FromContext(c.Request.Context()).WithError(err).Error("failed to change password")
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

// createToken 用户生成 API token
func (u *UserServer) createToken(c *gin.Context) {
	var req request.UserCreateTokenRequest
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
	le.Debug("user create api token called")

	// 从 JWT token 中获取用户ID和角色
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

	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    400,
				Message: "Token 名称不能为空",
			},
		})
		return
	}

	expireSeconds := req.Expire
	if expireSeconds > int((24 * time.Hour * 30).Seconds()) {
		expireSeconds = int((24 * time.Hour * 30).Seconds())
	}

	expDuration := time.Duration(expireSeconds) * time.Second

	// 生成 JWT token
	token, err := u.jwtManager.GenerateTokenWithExpire(userID, userRole, expDuration)
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

	expiresAt := time.Now().Add(expDuration)
	if _, err := u.tokenSvc.CreateToken(c.Request.Context(), userID, name, token, &expiresAt); err != nil {
		le.WithError(err).Error("failed to persist api token")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    500,
				Message: "保存令牌失败",
			},
		})
		return
	}

	le.Info("user token generated successfully")
	c.JSON(http.StatusOK, response.UserCreateTokenResponse{
		BaseResponse: response.BaseResponse{
			Code:    0,
			Message: "令牌生成成功",
		},
		Token: token,
	})
}

func (u *UserServer) listTokens(c *gin.Context) {
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

	tokens, err := u.tokenSvc.ListTokens(c.Request.Context(), userID)
	if err != nil {
		logger.FromContext(c.Request.Context()).WithError(err).Error("failed to list tokens")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    500,
				Message: "获取 Token 列表失败",
			},
		})
		return
	}

	var items []gin.H
	for _, token := range tokens {
		item := gin.H{
			"id":         token.ID,
			"name":       token.Name,
			"created_at": token.CreatedAt,
		}
		if token.ExpiresAt != nil {
			item["expires_at"] = token.ExpiresAt
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		BaseResponse: response.BaseResponse{
			Code:    0,
			Message: "获取 Token 列表成功",
		},
		Data: gin.H{
			"tokens": items,
		},
	})
}

func (u *UserServer) deleteToken(c *gin.Context) {
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

	tokenIDStr := c.Param("id")
	tokenID, err := strconv.ParseInt(tokenIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    400,
				Message: "无效的 Token ID",
			},
		})
		return
	}

	if err := u.tokenSvc.DeleteToken(c.Request.Context(), userID, tokenID); err != nil {
		if errors.Is(err, service.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, response.ErrorResponse{
				BaseResponse: response.BaseResponse{
					Code:    404,
					Message: "Token 不存在",
				},
			})
			return
		}

		logger.FromContext(c.Request.Context()).WithError(err).Error("failed to delete token")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			BaseResponse: response.BaseResponse{
				Code:    500,
				Message: "删除 Token 失败",
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		BaseResponse: response.BaseResponse{
			Code:    0,
			Message: "Token 已删除",
		},
	})
}
