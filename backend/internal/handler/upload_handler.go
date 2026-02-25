package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// queryInt 从查询参数中读取整数值，若参数缺失或无效则返回默认值。
func queryInt(c *gin.Context, key string, defaultVal int) int {
	s := c.Query(key)
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil || v <= 0 {
		return defaultVal
	}
	return v
}

// queryInt64 从查询参数中读取 int64 值，若参数缺失或无效则返回默认值。
func queryInt64(c *gin.Context, key string, defaultVal int64) int64 {
	s := c.Query(key)
	if s == "" {
		return defaultVal
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil || v <= 0 {
		return defaultVal
	}
	return v
}

// UploadHandler 处理文件上传端点（当前为桩实现，完整的 MinIO 集成将在 Sprint 2 完成）。
type UploadHandler struct{}

// NewUploadHandler 创建文件上传处理器实例。
func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

// UploadImage godoc
// @Summary Upload an image
// @Tags Upload
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/upload/image [post]
// UploadImage 上传图片到对象存储（当前为桩实现）。
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// 桩实现：完整的 MinIO 集成将在 Sprint 2 中完成
	c.JSON(http.StatusNotImplemented, gin.H{"code": 501, "message": "upload not yet implemented"})
}
