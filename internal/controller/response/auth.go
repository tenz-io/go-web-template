package response

// 统一认证响应结构体

// LoginResponseBody 登录成功响应体（用于 OkWithJson 的 data 字段）
type LoginResponseBody struct {
	Role     string `json:"role,omitempty"`
	Redirect string `json:"redirect,omitempty"` // 登录后重定向地址
}
