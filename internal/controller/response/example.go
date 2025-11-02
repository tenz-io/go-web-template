package response

// 示例：如何添加新的响应结构体

// UserResponse 用户信息响应
type UserResponse struct {
	BaseResponse
	Data UserData `json:"data"`
}

// UserData 用户数据
type UserData struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Age      int    `json:"age"`
	Status   string `json:"status"`
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	BaseResponse
	Data  []UserData `json:"data"`
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}

// 创建用户响应
func NewUserResponse(user UserData) UserResponse {
	return UserResponse{
		BaseResponse: BaseResponse{
			Code:    0,
			Message: "success",
		},
		Data: user,
	}
}

// 创建用户列表响应
func NewUserListResponse(users []UserData, total int64, page, size int) UserListResponse {
	return UserListResponse{
		BaseResponse: BaseResponse{
			Code:    0,
			Message: "success",
		},
		Data:  users,
		Total: total,
		Page:  page,
		Size:  size,
	}
}
