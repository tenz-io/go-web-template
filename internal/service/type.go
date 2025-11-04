package service

// CreateUserParam 创建用户请求
type CreateUserParam struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     int32  `json:"role" binding:"required"`
}

type UpdatePasswordParam struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type VerifyUserParam struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}
