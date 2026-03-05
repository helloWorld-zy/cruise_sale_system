package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

const defaultUploadMaxBytes int64 = 10 * 1024 * 1024 // 10MB

var allowedImageMIMEs = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
	"image/gif":  {},
}

// UploadHandler 处理文件上传端点。
type UploadHandler struct {
	uploadDir   string
	publicPath  string
	maxFileSize int64
}

// NewUploadHandler 创建文件上传处理器实例。
func NewUploadHandler() *UploadHandler {
	return &UploadHandler{
		uploadDir:   "uploads",
		publicPath:  "/uploads",
		maxFileSize: defaultUploadMaxBytes,
	}
}

// NewUploadHandlerWithConfig 使用自定义存储目录和公开路径创建上传处理器。
func NewUploadHandlerWithConfig(uploadDir, publicPath string, maxFileSize int64) *UploadHandler {
	if strings.TrimSpace(uploadDir) == "" {
		uploadDir = "uploads"
	}
	if strings.TrimSpace(publicPath) == "" {
		publicPath = "/uploads"
	}
	if maxFileSize <= 0 {
		maxFileSize = defaultUploadMaxBytes
	}
	return &UploadHandler{uploadDir: uploadDir, publicPath: publicPath, maxFileSize: maxFileSize}
}

// StorageDir 返回上传文件在本地磁盘的存储目录。
func (h *UploadHandler) StorageDir() string {
	if h == nil || strings.TrimSpace(h.uploadDir) == "" {
		return "uploads"
	}
	return h.uploadDir
}

// PublicBasePath 返回上传文件对外暴露的 URL 前缀。
func (h *UploadHandler) PublicBasePath() string {
	if h == nil || strings.TrimSpace(h.publicPath) == "" {
		return "/uploads"
	}
	return h.publicPath
}

// UploadImage godoc
// @Summary Upload an image
// @Tags Upload
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file"
// @Success 200 {object} gin.H
// @Router /api/v1/admin/upload/image [post]
// UploadImage 上传图片并返回可直接访问的 URL。
func (h *UploadHandler) UploadImage(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "file is required"})
		return
	}

	if fileHeader.Size <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "empty file is not allowed"})
		return
	}
	if fileHeader.Size > h.maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": fmt.Sprintf("file size exceeds limit (%d bytes)", h.maxFileSize)})
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to read upload file"})
		return
	}
	defer src.Close()

	buf := make([]byte, 512)
	n, _ := src.Read(buf)
	contentType := http.DetectContentType(buf[:n])
	if _, ok := allowedImageMIMEs[contentType]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "unsupported image format"})
		return
	}

	if seeker, ok := src.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to process upload file"})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "upload stream is not seekable"})
		return
	}

	if err := os.MkdirAll(h.StorageDir(), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to prepare upload directory"})
		return
	}

	fileName := buildUploadFileName(fileHeader, contentType)
	filePath := filepath.Join(h.StorageDir(), fileName)
	if err := saveUploadedFile(src, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to save upload file"})
		return
	}

	urlPath := strings.TrimRight(h.PublicBasePath(), "/") + "/" + fileName
	fullURL := buildPublicURL(c, urlPath)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data": gin.H{
			"url": fullURL,
		},
	})
}

func saveUploadedFile(src multipart.File, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return err
}

func buildUploadFileName(fileHeader *multipart.FileHeader, contentType string) string {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		ext = extByContentType(contentType)
	}
	return fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), fileHeader.Size, ext)
}

func extByContentType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	default:
		return ""
	}
}

func buildPublicURL(c *gin.Context, urlPath string) string {
	scheme := c.GetHeader("X-Forwarded-Proto")
	if scheme == "" {
		if c.Request != nil && c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	return fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, urlPath)
}
