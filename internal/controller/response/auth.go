package response

// 统一认证响应结构体

// LoginResponse 统一登录响应（已废弃，使用 CommonResponse）
type LoginResponse struct {
	BaseResponse
	Role     string `json:"role,omitempty"`
	Redirect string `json:"redirect,omitempty"` // 登录后重定向地址
}

// LoginResponseBody 登录成功响应体（用于 OkWithJson 的 data 字段）
type LoginResponseBody struct {
	Role     string `json:"role,omitempty"`
	Redirect string `json:"redirect,omitempty"` // 登录后重定向地址
}
