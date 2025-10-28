package response

// 统一认证响应结构体

// LoginResponse 统一登录响应
type LoginResponse struct {
	BaseResponse
	Role     string `json:"role,omitempty"`
	Redirect string `json:"redirect,omitempty"` // 登录后重定向地址
}
