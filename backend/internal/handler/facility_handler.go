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

// FacilityHandler handles CRUD for facilities.
type FacilityHandler struct {
	svc *service.FacilityService
}

// NewFacilityHandler creates a handler with injected service.
func NewFacilityHandler(svc *service.FacilityService) *FacilityHandler {
	return &FacilityHandler{svc: svc}
}

// FacilityRequest is the create/update payload.
type FacilityRequest struct {
	CategoryID  int64  `json:"category_id" binding:"required"`
	CruiseID    int64  `json:"cruise_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	EnglishName string `json:"english_name"`
	Location    string `json:"location"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// ListByCruise godoc
// @Summary List facilities by cruise
// @Tags Facility
// @Security BearerAuth
// @Param cruise_id query int true "Cruise ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/facilities [get]
func (h *FacilityHandler) ListByCruise(c *gin.Context) {
	cruiseID, err := strconv.ParseInt(c.Query("cruise_id"), 10, 64)
	if err != nil || cruiseID <= 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid cruise_id")
		return
	}
	items, err := h.svc.ListByCruise(c.Request.Context(), cruiseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, items)
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
func (h *FacilityHandler) Create(c *gin.Context) {
	var req FacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	f := &domain.Facility{
		CategoryID:  req.CategoryID,
		CruiseID:    req.CruiseID,
		Name:        req.Name,
		EnglishName: req.EnglishName,
		Location:    req.Location,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}
	if err := h.svc.Create(c.Request.Context(), f); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, f)
}

// Delete godoc
// @Summary Delete a facility
// @Tags Facility
// @Security BearerAuth
// @Param id path int true "Facility ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/facilities/{id} [delete]
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
