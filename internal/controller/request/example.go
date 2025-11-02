package request

// 示例：如何添加新的请求结构体

// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Age      int    `json:"age" binding:"min=1,max=120"`
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	ID       int    `json:"id" binding:"required"`
	Username string `json:"username,omitempty"`
	Age      int    `json:"age,omitempty" binding:"omitempty,min=1,max=120"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	Size     int    `form:"size" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Status   string `form:"status"`
	OrderBy  string `form:"order_by"`
	OrderDir string `form:"order_dir" binding:"omitempty,oneof=asc desc"`
}
