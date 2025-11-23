package response

import "time"

// 管理员响应结构体

// UserItem 用户列表项
type UserItem struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// AdminGetUsersResponseBody 获取用户列表响应体
type AdminGetUsersResponseBody struct {
	Users []UserItem `json:"users"`
	Total int64      `json:"total"`
}

// AdminAddUserResponseBody 添加用户响应体
type AdminAddUserResponseBody struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Role      int32     `json:"role"`
	Profile   string    `json:"profile"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
