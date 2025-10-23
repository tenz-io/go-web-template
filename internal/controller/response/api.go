package response

// API 响应结构体

// HelloResponse Hello 接口响应
type HelloResponse struct {
	BaseResponse
}

// UploadResponse 文件上传响应
type UploadResponse struct {
	BaseResponse
	Key string `json:"key"`
}
