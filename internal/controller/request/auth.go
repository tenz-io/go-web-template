package request

// 统一认证请求结构体

// LoginRequest 统一登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
