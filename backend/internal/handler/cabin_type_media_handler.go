package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// CabinTypeMediaHandler 处理舱型媒体资源端点。
type CabinTypeMediaHandler struct {
	svc         *service.CabinTypeMediaService
	uploadDir   string
	publicPath  string
	maxFileSize int64
}

func NewCabinTypeMediaHandler(svc *service.CabinTypeMediaService, uploadDir, publicPath string, maxFileSize int64) *CabinTypeMediaHandler {
	if strings.TrimSpace(uploadDir) == "" {
		uploadDir = "uploads"
	}
	if strings.TrimSpace(publicPath) == "" {
		publicPath = "/uploads"
	}
	if maxFileSize <= 0 {
		maxFileSize = defaultUploadMaxBytes
	}
	return &CabinTypeMediaHandler{svc: svc, uploadDir: uploadDir, publicPath: publicPath, maxFileSize: maxFileSize}
}

type CabinTypeMediaRequest struct {
	MediaType string `json:"media_type" binding:"required"`
	URL       string `json:"url" binding:"required"`
	Title     string `json:"title"`
	SortOrder int    `json:"sort_order"`
	IsPrimary bool   `json:"is_primary"`
}

func (h *CabinTypeMediaHandler) List(c *gin.Context) {
	cabinTypeID, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}
	items, err := h.svc.ListByCabinType(c.Request.Context(), cabinTypeID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	if mediaType := strings.TrimSpace(c.Query("media_type")); mediaType != "" {
		filtered := make([]domain.CabinTypeMedia, 0, len(items))
		for _, item := range items {
			if item.MediaType == mediaType {
				filtered = append(filtered, item)
			}
		}
		response.Success(c, filtered)
		return
	}
	response.Success(c, items)
}

func (h *CabinTypeMediaHandler) Create(c *gin.Context) {
	cabinTypeID, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}
	var req CabinTypeMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	if !isValidCabinMediaType(req.MediaType) {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid media_type")
		return
	}
	item := &domain.CabinTypeMedia{
		CabinTypeID: cabinTypeID,
		MediaType:   req.MediaType,
		URL:         req.URL,
		Title:       req.Title,
		SortOrder:   req.SortOrder,
		IsPrimary:   req.IsPrimary,
	}
	if err := h.svc.Create(c.Request.Context(), item); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, item)
}

func (h *CabinTypeMediaHandler) Update(c *gin.Context) {
	cabinTypeID, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}
	mediaID, ok := parsePositiveID(c, "mediaId")
	if !ok {
		return
	}
	item, err := h.svc.GetByID(c.Request.Context(), mediaID)
	if err != nil || item.CabinTypeID != cabinTypeID {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "cabin type media not found")
		return
	}
	var req CabinTypeMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	if !isValidCabinMediaType(req.MediaType) {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid media_type")
		return
	}
	item.MediaType = req.MediaType
	item.URL = req.URL
	item.Title = req.Title
	item.SortOrder = req.SortOrder
	item.IsPrimary = req.IsPrimary
	if err := h.svc.Update(c.Request.Context(), item); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, item)
}

func (h *CabinTypeMediaHandler) Delete(c *gin.Context) {
	mediaID, ok := parsePositiveID(c, "mediaId")
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), mediaID); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *CabinTypeMediaHandler) Upload(c *gin.Context) {
	cabinTypeID, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}
	mediaType := strings.TrimSpace(c.PostForm("media_type"))
	if !isValidCabinMediaType(mediaType) {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid media_type")
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "file is required")
		return
	}
	if fileHeader.Size <= 0 || fileHeader.Size > h.maxFileSize {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid file size")
		return
	}
	src, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, "failed to read upload file")
		return
	}
	defer src.Close()

	buf := make([]byte, 512)
	n, _ := src.Read(buf)
	contentType := http.DetectContentType(buf[:n])
	if _, ok := allowedImageMIMEs[contentType]; !ok {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "unsupported image format")
		return
	}
	if seeker, ok := src.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}
	if err := os.MkdirAll(h.uploadDir, 0755); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, "failed to prepare upload directory")
		return
	}
	fileName := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), fileHeader.Size, strings.ToLower(filepath.Ext(fileHeader.Filename)))
	if filepath.Ext(fileName) == "" {
		fileName += extByContentType(contentType)
	}
	filePath := filepath.Join(h.uploadDir, fileName)
	out, err := os.Create(filePath)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, "failed to save upload file")
		return
	}
	defer out.Close()
	if _, err := io.Copy(out, src); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, "failed to save upload file")
		return
	}
	urlPath := strings.TrimRight(h.publicPath, "/") + "/" + fileName
	fullURL := buildPublicURL(c, urlPath)
	sortOrder, _ := strconv.Atoi(c.DefaultPostForm("sort_order", "0"))
	isPrimary := c.DefaultPostForm("is_primary", "false") == "true"
	item := &domain.CabinTypeMedia{
		CabinTypeID: cabinTypeID,
		MediaType:   mediaType,
		URL:         fullURL,
		Title:       strings.TrimSpace(c.PostForm("title")),
		SortOrder:   sortOrder,
		IsPrimary:   isPrimary,
	}
	if err := h.svc.Create(c.Request.Context(), item); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, item)
}

func parsePositiveID(c *gin.Context, key string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(key), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return 0, false
	}
	return id, true
}

func isValidCabinMediaType(mediaType string) bool {
	return mediaType == domain.CabinTypeMediaTypeImage || mediaType == domain.CabinTypeMediaTypeFloorPlan
}
