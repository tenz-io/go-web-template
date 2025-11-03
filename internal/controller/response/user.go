package response

// 普通用户响应结构体

// UserGenerateTokenResponse 用户生成API token响应
type UserGenerateTokenResponse struct {
	BaseResponse
	Token string `json:"token,omitempty"`
}
