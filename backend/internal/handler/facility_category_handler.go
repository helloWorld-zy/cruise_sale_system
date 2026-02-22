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

// FacilityCategoryHandler handles CRUD for facility categories.
type FacilityCategoryHandler struct {
	svc *service.FacilityCategoryService
}

// NewFacilityCategoryHandler creates a handler with injected service.
func NewFacilityCategoryHandler(svc *service.FacilityCategoryService) *FacilityCategoryHandler {
	return &FacilityCategoryHandler{svc: svc}
}

// FacilityCategoryRequest is the create payload.
type FacilityCategoryRequest struct {
	Name      string `json:"name" binding:"required"`
	Icon      string `json:"icon"`
	SortOrder int    `json:"sort_order"`
}

// List godoc
// @Summary List all facility categories
// @Tags FacilityCategory
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/admin/facility-categories [get]
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
