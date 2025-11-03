package request

// Admin 请求结构体

// AdminChangePasswordRequest 管理员修改密码请求
type AdminChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// AdminAddUserRequest 管理员添加用户请求
type AdminAddUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=user admin"`
}

// AdminDeleteUserRequest 管理员删除用户请求
type AdminDeleteUserRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}
