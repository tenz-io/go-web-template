package request

// API 请求结构体

// HelloRequest Hello 接口请求
type HelloRequest struct {
	Name string `form:"name" binding:"required"`
}
