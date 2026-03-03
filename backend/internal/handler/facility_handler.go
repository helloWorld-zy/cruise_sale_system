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

// FacilityHandler 处理设施的 CRUD 端点。
type FacilityHandler struct {
	svc *service.FacilityService // 设施服务
}

// NewFacilityHandler 创建设施处理器，通过依赖注入传入服务。
func NewFacilityHandler(svc *service.FacilityService) *FacilityHandler {
	return &FacilityHandler{svc: svc}
}

// FacilityRequest 是创建/更新设施的请求体结构。
type FacilityRequest struct {
	CategoryID     int64  `json:"category_id" binding:"required"` // 设施分类 ID（必填）
	CruiseID       int64  `json:"cruise_id" binding:"required"`   // 所属邮轮 ID（必填）
	Name           string `json:"name" binding:"required"`        // 设施名称（必填）
	EnglishName    string `json:"english_name"`                   // 英文名称
	Location       string `json:"location"`                       // 设施位置
	OpenHours      string `json:"open_hours"`                     // 开放时间
	ExtraCharge    bool   `json:"extra_charge"`                   // 是否额外收费
	ChargePriceTip string `json:"charge_price_tip"`               // 收费提示
	TargetAudience string `json:"target_audience"`                // 适合人群
	Description    string `json:"description"`                    // 设施描述
	Status         int16  `json:"status"`                         // 状态
	SortOrder      int    `json:"sort_order"`                     // 排序权重
}

// ListByCruise godoc
// @Summary List facilities by cruise
// @Tags Facility
// @Security BearerAuth
// @Param cruise_id query int true "Cruise ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/facilities [get]
// ListByCruise 查询指定邮轮下的所有设施。
func (h *FacilityHandler) ListByCruise(c *gin.Context) {
	cruiseID, err := strconv.ParseInt(c.Query("cruise_id"), 10, 64)
	if err != nil || cruiseID <= 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid cruise_id")
		return
	}
	categoryID := queryInt64(c, "category_id", 0)
	var items []domain.Facility
	if categoryID > 0 {
		items, err = h.svc.ListByCruiseAndCategory(c.Request.Context(), cruiseID, categoryID)
	} else {
		items, err = h.svc.ListByCruise(c.Request.Context(), cruiseID)
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, items)
}

// Get godoc
// @Summary Get a facility
// @Tags Facility
// @Security BearerAuth
// @Param id path int true "Facility ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/facilities/{id} [get]
// Get 查询指定设施详情。
func (h *FacilityHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	item, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "facility not found")
		return
	}
	response.Success(c, item)
}

// Create godoc
// @Summary Create a facility
// @Tags Facility
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body FacilityRequest true "Facility data"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/facilities [post]
// Create 创建新的设施。
func (h *FacilityHandler) Create(c *gin.Context) {
	var req FacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	f := &domain.Facility{
		CategoryID:     req.CategoryID,
		CruiseID:       req.CruiseID,
		Name:           req.Name,
		EnglishName:    req.EnglishName,
		Location:       req.Location,
		OpenHours:      req.OpenHours,
		ExtraCharge:    req.ExtraCharge,
		ChargePriceTip: req.ChargePriceTip,
		TargetAudience: req.TargetAudience,
		Description:    req.Description,
		Status:         req.Status,
		SortOrder:      req.SortOrder,
	}
	if err := h.svc.Create(c.Request.Context(), f); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, f)
}

// Update godoc
// @Summary Update a facility
// @Tags Facility
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Facility ID"
// @Param body body FacilityRequest true "Facility data"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/facilities/{id} [put]
// Update 更新指定设施。
func (h *FacilityHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	existing, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "facility not found")
		return
	}
	var req FacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	existing.CategoryID = req.CategoryID
	existing.CruiseID = req.CruiseID
	existing.Name = req.Name
	existing.EnglishName = req.EnglishName
	existing.Location = req.Location
	existing.OpenHours = req.OpenHours
	existing.ExtraCharge = req.ExtraCharge
	existing.ChargePriceTip = req.ChargePriceTip
	existing.TargetAudience = req.TargetAudience
	existing.Description = req.Description
	existing.Status = req.Status
	existing.SortOrder = req.SortOrder

	if err := h.svc.Update(c.Request.Context(), existing); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, existing)
}

// Delete godoc
// @Summary Delete a facility
// @Tags Facility
// @Security BearerAuth
// @Param id path int true "Facility ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/facilities/{id} [delete]
// Delete 删除指定的设施。
func (h *FacilityHandler) Delete(c *gin.Context) {
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
