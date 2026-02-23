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

// FacilityCategoryHandler 处理设施分类的 CRUD 端点。
type FacilityCategoryHandler struct {
	svc *service.FacilityCategoryService // 设施分类服务
}

// NewFacilityCategoryHandler 创建设施分类处理器，通过依赖注入传入服务。
func NewFacilityCategoryHandler(svc *service.FacilityCategoryService) *FacilityCategoryHandler {
	return &FacilityCategoryHandler{svc: svc}
}

// FacilityCategoryRequest 是创建设施分类的请求体结构。
type FacilityCategoryRequest struct {
	Name      string `json:"name" binding:"required"` // 分类名称（必填）
	Icon      string `json:"icon"`                    // 分类图标
	SortOrder int    `json:"sort_order"`              // 排序权重
}

// List godoc
// @Summary List all facility categories
// @Tags FacilityCategory
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/admin/facility-categories [get]
// List 查询所有设施分类。
func (h *FacilityCategoryHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, items)
}

// Create godoc
// @Summary Create a facility category
// @Tags FacilityCategory
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body FacilityCategoryRequest true "Category data"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/facility-categories [post]
// Create 创建新的设施分类。
func (h *FacilityCategoryHandler) Create(c *gin.Context) {
	var req FacilityCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	cat := &domain.FacilityCategory{
		Name:      req.Name,
		Icon:      req.Icon,
		SortOrder: req.SortOrder,
	}
	if err := h.svc.Create(c.Request.Context(), cat); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, cat)
}

// Delete godoc
// @Summary Delete a facility category
// @Tags FacilityCategory
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/facility-categories/{id} [delete]
// Delete 删除指定的设施分类。
func (h *FacilityCategoryHandler) Delete(c *gin.Context) {
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
