package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// queryInt reads a query param as int with a default fallback.
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

// queryInt64 reads a query param as int64 with a default fallback.
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

// paramInt64 reads a path param as int64.
func paramInt64(c *gin.Context, key string) (int64, error) {
	return strconv.ParseInt(c.Param(key), 10, 64)
}

// UploadHandler handles file uploads to MinIO (stub — full impl in Sprint 2).
type UploadHandler struct{}

// NewUploadHandler creates an UploadHandler.
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
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// Stub: full MinIO integration will be implemented in Sprint 2.
	c.JSON(http.StatusNotImplemented, gin.H{"code": 501, "message": "upload not yet implemented"})
}
