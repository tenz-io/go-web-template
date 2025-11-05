package request

// 普通用户请求结构体

// UserGenerateTokenRequest 用户生成API token请求
type UserGenerateTokenRequest struct {
	ExpireHours int `json:"expire_hours" binding:"required,min=1"` // 过期时间（小时）
}
