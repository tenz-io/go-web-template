package controller

import (
	"go-web-template/internal/controller/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/logger"

	"go-web-template/internal/config"
	"go-web-template/internal/constant"
	"go-web-template/internal/controller/request"
	"go-web-template/internal/controller/response"
	"go-web-template/internal/service"
)

type AdminController struct {
	userService service.User
	jwtManager  *middleware.JWTManager
	appName     string
}

func NewAdminController(cfg *config.Config, userService service.User, jwtManager *middleware.JWTManager) *AdminController {
	return &AdminController{
		userService: userService,
		jwtManager:  jwtManager,
		appName:     cfg.App.Name,
	}
}

// 注册管理路由
func (ac *AdminController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/home", ac.home)
	r.GET("/users", ac.GetUsers)
	r.POST("/add_user", ac.AddUser)
	r.DELETE("/delete_user", ac.DeleteUser)
}

func (ac *AdminController) home(c *gin.Context) {
	le := logger.FromContext(c.Request.Context())

	userID, _, err := middleware.GetUserInfoFromContext(c)
	if err != nil {
		le.WithError(err).Warn("failed to get user info from context")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	userModel, err := ac.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		le.WithError(err).Error("failed to get user model")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	role := constant.Role(userModel.Role)

	c.HTML(http.StatusOK, "home.html", gin.H{
		"appName":     ac.appName,
		"name":        ac.appName,
		"username":    userModel.Username,
		"displayName": userModel.Username,
		"role":        role.String(),
		"isAdmin":     role == constant.RoleAdmin,
	})
}

// GetUsers 获取用户列表
func (ac *AdminController) GetUsers(c *gin.Context) {
	le := logger.FromContext(c.Request.Context())
	le.Debug("admin get users called")

	// 获取分页参数
	limit := 100 // 默认限制
	offset := 0

	// 获取用户列表
	users, total, err := ac.userService.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		le.Error("failed to get users")
		response.FailWithJson(c, 500, "获取用户列表失败")
		return
	}

	// 转换用户角色为字符串
	var userList []response.UserItem
	for _, user := range users {
		userList = append(userList, response.UserItem{
			ID:        user.ID,
			Username:  user.Username,
			Role:      constant.Role(user.Role).String(),
			CreatedAt: user.CreatedAt,
		})
	}

	response.OkWithJson(c, response.AdminGetUsersResponseBody{
		Users: userList,
		Total: total,
	})
}

// AddUser 添加用户
func (ac *AdminController) AddUser(c *gin.Context) {
	var req request.AdminAddUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithJson(c, 400, "请求参数错误")
		return
	}

	le := logger.FromContext(c.Request.Context())
	le.Debug("admin add user called")

	// 将角色字符串转换为数字
	role, err := constant.ParseRole(req.Role)
	if err != nil {
		le.Error("invalid role")
		response.FailWithJson(c, 400, "无效的用户角色")
		return
	}

	// 创建用户
	createParam := service.CreateUserParam{
		Username: req.Username,
		Password: req.Password,
		Role:     int32(role),
	}
	user, err := ac.userService.CreateUser(c.Request.Context(), createParam)
	if err != nil {
		le.Error("failed to create user")
		response.FailWithJson(c, 400, "创建用户失败: "+err.Error())
		return
	}

	le.Info("user created successfully")
	response.OkWithJson(c, response.AdminAddUserResponseBody{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		Profile:   user.Profile,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// DeleteUser 删除用户
func (ac *AdminController) DeleteUser(c *gin.Context) {
	var req request.AdminDeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithJson(c, 400, "请求参数错误")
		return
	}

	le := logger.FromContext(c.Request.Context())
	le.Debug("admin delete user called")

	// 删除用户
	err := ac.userService.DeleteUser(c.Request.Context(), req.UserID)
	if err != nil {
		le.Error("failed to delete user")
		response.FailWithJson(c, 400, "删除用户失败: "+err.Error())
		return
	}

	response.OkWithJson(c, nil)
}
