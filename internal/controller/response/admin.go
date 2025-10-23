package response

// Admin 响应结构体

// AdminLoginResponse 管理员登录响应
type AdminLoginResponse struct {
	BaseResponse
}

// AdminAddTokenResponse 生成访问令牌响应
type AdminAddTokenResponse struct {
	BaseResponse
	AccessToken string `json:"access_token,omitempty"`
}

// AdminChangePasswordResponse 管理员修改密码响应
type AdminChangePasswordResponse struct {
	BaseResponse
}
