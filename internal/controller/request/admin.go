package request

// Admin 请求结构体

// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminAddTokenRequest 生成访问令牌请求
type AdminAddTokenRequest struct {
	UserID int `json:"userid" binding:"required"`
	Expire int `json:"expire" binding:"required"`
}

// AdminChangePasswordRequest 管理员修改密码请求
type AdminChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
