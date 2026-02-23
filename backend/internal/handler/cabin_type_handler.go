package handler

import (
	"net/http"
	"strconv"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// CabinTypeHandler 处理舱房类型的 CRUD 端点。
type CabinTypeHandler struct {
	svc *service.CabinTypeService // 舱房类型服务
}

// NewCabinTypeHandler 创建舱房类型处理器，通过依赖注入传入服务。
func NewCabinTypeHandler(svc *service.CabinTypeService) *CabinTypeHandler {
	return &CabinTypeHandler{svc: svc}
}

// CabinTypeRequest 是创建/更新舱房类型的请求体结构。
type CabinTypeRequest struct {
	CruiseID    int64   `json:"cruise_id" binding:"required"` // 所属邮轮 ID
	Name        string  `json:"name" binding:"required"`      // 舱房类型名称
	EnglishName string  `json:"english_name"`                 // 英文名称
	Capacity    int     `json:"capacity"`                     // 容纳人数
	Area        float64 `json:"area"`                         // 面积（平方米）
	Deck        string  `json:"deck"`                         // 所在甲板
	Description string  `json:"description"`                  // 描述
	SortOrder   int     `json:"sort_order"`                   // 排序权重
}

// List godoc
// @Summary List cabin types by cruise
// @Tags CabinType
// @Security BearerAuth
// @Produce json
// @Param cruise_id query int true "Cruise ID"
// @Param page query int false "Page" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/admin/cabin-types [get]
// List 查询指定邮轮下的舱房类型列表（分页）。
func (h *CabinTypeHandler) List(c *gin.Context) {
	cruiseID, err := strconv.ParseInt(c.Query("cruise_id"), 10, 64)
	if err != nil || cruiseID <= 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid cruise_id")
		return
	}
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "page_size", 10)

	items, total, err := h.svc.List(c.Request.Context(), cruiseID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, gin.H{"list": items, "total": total})
}

// Create godoc
// @Summary Create a cabin type
// @Tags CabinType
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CabinTypeRequest true "Cabin type data"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/cabin-types [post]
// Create 创建新的舱房类型。
func (h *CabinTypeHandler) Create(c *gin.Context) {
	var req CabinTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	ct := &domain.CabinType{
		CruiseID:    req.CruiseID,
		Name:        req.Name,
		EnglishName: req.EnglishName,
		Capacity:    req.Capacity,
		Area:        req.Area,
		Deck:        req.Deck,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}
	if err := h.svc.Create(c.Request.Context(), ct); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, ct)
}

// Update godoc
// @Summary Update a cabin type
// @Tags CabinType
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "CabinType ID"
// @Param body body CabinTypeRequest true "Cabin type data"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/cabin-types/{id} [put]
// Update 更新指定的舱房类型。
func (h *CabinTypeHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	// 先查询现有数据
	existing, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "cabin type not found")
		return
	}
	var req CabinTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	// 更新字段
	existing.Name = req.Name
	existing.EnglishName = req.EnglishName
	existing.Capacity = req.Capacity
	existing.Area = req.Area
	existing.Deck = req.Deck
	existing.Description = req.Description
	existing.SortOrder = req.SortOrder

	if err := h.svc.Update(c.Request.Context(), existing); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, existing)
}

// Delete godoc
// @Summary Delete a cabin type
// @Tags CabinType
// @Security BearerAuth
// @Param id path int true "CabinType ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/cabin-types/{id} [delete]
// Delete 删除指定的舱房类型。
func (h *CabinTypeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, nil)
}
