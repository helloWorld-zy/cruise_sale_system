package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the unified API envelope.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success writes HTTP 200 with business code 0.
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: data})
}

// Error writes the given HTTP status code and business error code.
// ME-06 FIX: use proper HTTP status codes instead of always returning 200.
func Error(c *gin.Context, httpStatus int, businessCode int, message string) {
	c.JSON(httpStatus, Response{Code: businessCode, Message: message, Data: nil})
}
