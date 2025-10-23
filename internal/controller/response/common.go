package response

// 通用响应结构体

// BaseResponse 基础响应结构
type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	BaseResponse
	Data interface{} `json:"data,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	BaseResponse
}

// PageResponse 分页响应
type PageResponse struct {
	BaseResponse
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// 创建成功响应
func NewSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{
		BaseResponse: BaseResponse{
			Code:    0,
			Message: "success",
		},
		Data: data,
	}
}

// 创建错误响应
func NewErrorResponse(code int, message string) ErrorResponse {
	return ErrorResponse{
		BaseResponse: BaseResponse{
			Code:    code,
			Message: message,
		},
	}
}

// 创建分页响应
func NewPageResponse(data interface{}, total int64, page, size int) PageResponse {
	return PageResponse{
		BaseResponse: BaseResponse{
			Code:    0,
			Message: "success",
		},
		Data:  data,
		Total: total,
		Page:  page,
		Size:  size,
	}
}
