package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 通用响应结构体

// BaseResponse 基础响应结构
type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// CommonResponse 通用响应
type CommonResponse struct {
	BaseResponse
	Data any `json:"data,omitempty"`
}

// 创建成功响应
func Ok(data any) CommonResponse {
	return CommonResponse{
		BaseResponse: BaseResponse{
			Code:    0,
			Message: "success",
		},
		Data: data,
	}
}

// 创建错误响应
func Fail(code int, message string) CommonResponse {
	return CommonResponse{
		BaseResponse: BaseResponse{
			Code:    code,
			Message: message,
		},
	}
}

// OkWithJson 返回成功响应，自动使用 HTTP 200 状态码
func OkWithJson(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Ok(data))
}

// FailWithJson 返回错误响应，自动从 code 推断 HTTP 状态码
func FailWithJson(c *gin.Context, code int, message string) {
	httpStatus := codeToHTTPStatus(code)
	c.JSON(httpStatus, Fail(code, message))
}

// codeToHTTPStatus 将业务错误码转换为 HTTP 状态码
func codeToHTTPStatus(code int) int {
	// 常见的 HTTP 状态码直接映射
	switch code {
	case 400:
		return http.StatusBadRequest
	case 401:
		return http.StatusUnauthorized
	case 403:
		return http.StatusForbidden
	case 404:
		return http.StatusNotFound
	case 500:
		return http.StatusInternalServerError
	default:
		// 如果 code 在有效的 HTTP 状态码范围内（100-599），直接使用
		if code >= 100 && code <= 599 {
			return code
		}
		// 否则默认返回 500
		return http.StatusInternalServerError
	}
}
