package handler

import (
	"net/http"

	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// ImageHandler 处理图片画廊管理接口。
type ImageHandler struct {
	svc *service.ImageService // 图片服务
}

// NewImageHandler 创建图片处理器。
func NewImageHandler(svc *service.ImageService) *ImageHandler {
	return &ImageHandler{svc: svc}
}

// ImageItemRequest 表示单张图片输入。
type ImageItemRequest struct {
	URL       string `json:"url" binding:"required"` // 图片地址
	SortOrder int    `json:"sort_order"`             // 排序
	IsPrimary bool   `json:"is_primary"`             // 是否主图
}

// SaveImagesRequest 表示保存图片列表请求。
type SaveImagesRequest struct {
	EntityType string             `json:"entity_type" binding:"required"` // 实体类型
	EntityID   int64              `json:"entity_id" binding:"required"`   // 实体 ID
	Images     []ImageItemRequest `json:"images"`                         // 图片列表
}

// Save godoc
// @Summary Save images for entity
// @Tags Image
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body SaveImagesRequest true "Image payload"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/images [post]
// Save 保存实体关联的图片列表。
func (h *ImageHandler) Save(c *gin.Context) {
	var req SaveImagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}

	inputs := make([]service.ImageInput, 0, len(req.Images))
	for _, item := range req.Images {
		inputs = append(inputs, service.ImageInput{
			URL:       item.URL,
			SortOrder: item.SortOrder,
			IsPrimary: item.IsPrimary,
		})
	}

	if err := h.svc.SetImages(c.Request.Context(), req.EntityType, req.EntityID, inputs); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, gin.H{"count": len(inputs)})
}

// List godoc
// @Summary List images by entity
// @Tags Image
// @Security BearerAuth
// @Produce json
// @Param entity_type query string true "Entity type"
// @Param entity_id query int true "Entity ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/images [get]
// List 查询实体关联的图片列表。
func (h *ImageHandler) List(c *gin.Context) {
	entityType := c.Query("entity_type")
	if entityType == "" {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid entity_type")
		return
	}
	entityID := queryInt64(c, "entity_id", 0)
	if entityID <= 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid entity_id")
		return
	}
	items, err := h.svc.ListImages(c.Request.Context(), entityType, entityID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, items)
}
