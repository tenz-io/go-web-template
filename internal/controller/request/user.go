package request

// 普通用户请求结构体

// UserChangePasswordRequest 用户修改密码请求
type UserChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UserGenerateTokenRequest 用户生成API token请求
type UserGenerateTokenRequest struct {
	Expire int `json:"expire" binding:"required,min=1"` // 过期时间（秒），最大24小时
}
