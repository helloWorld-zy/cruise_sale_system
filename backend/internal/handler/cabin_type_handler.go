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

// CabinTypeHandler handles CRUD for cabin types.
type CabinTypeHandler struct {
	svc *service.CabinTypeService
}

// NewCabinTypeHandler creates a CabinTypeHandler with injected service.
func NewCabinTypeHandler(svc *service.CabinTypeService) *CabinTypeHandler {
	return &CabinTypeHandler{svc: svc}
}

// CabinTypeRequest is the common create/update payload.
type CabinTypeRequest struct {
	CruiseID    int64   `json:"cruise_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	EnglishName string  `json:"english_name"`
	Capacity    int     `json:"capacity"`
	Area        float64 `json:"area"`
	Deck        string  `json:"deck"`
	Description string  `json:"description"`
	SortOrder   int     `json:"sort_order"`
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
func (h *CabinTypeHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
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
