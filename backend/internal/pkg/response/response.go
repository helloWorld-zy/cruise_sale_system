package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 是统一的 API 响应信封结构。
// 所有 API 响应都通过此结构返回，确保前端解析格式一致。
type Response struct {
	Code    int         `json:"code"`    // 业务状态码（0 表示成功）
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
}

// Success 返回 HTTP 200 成功响应，业务状态码为 0。
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: data})
}

// Error 返回错误响应，使用指定的 HTTP 状态码和业务错误码。
// ME-06 修复：使用正确的 HTTP 状态码，而非始终返回 200。
func Error(c *gin.Context, httpStatus int, businessCode int, message string) {
	c.JSON(httpStatus, Response{Code: businessCode, Message: message, Data: nil})
}

// InternalError 返回 HTTP 500 服务器内部错误响应。
func InternalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, Response{Code: 500, Message: err.Error(), Data: nil})
}
