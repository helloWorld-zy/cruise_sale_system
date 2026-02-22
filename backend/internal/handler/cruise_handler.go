package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// CruiseHandler handles CRUD for cruises.
// CR-05: Handler now has proper DI via CruiseService.
type CruiseHandler struct {
	svc *service.CruiseService
}

// NewCruiseHandler creates a CruiseHandler with injected service.
func NewCruiseHandler(svc *service.CruiseService) *CruiseHandler {
	return &CruiseHandler{svc: svc}
}

// CruiseRequest is the create/update payload.
type CruiseRequest struct {
	CompanyID         int64   `json:"company_id" binding:"required"`
	Name              string  `json:"name" binding:"required"`
	EnglishName       string  `json:"english_name"`
	BuildYear         int     `json:"build_year"`
	Tonnage           float64 `json:"tonnage"`
	PassengerCapacity int     `json:"passenger_capacity"`
	RoomCount         int     `json:"room_count"`
	Description       string  `json:"description"`
	SortOrder         int     `json:"sort_order"`
}

// List godoc
// @Summary List cruises
// @Tags Cruise
// @Security BearerAuth
// @Produce json
// @Param company_id query int false "Filter by company"
// @Param page query int false "Page" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/admin/cruises [get]
func (h *CruiseHandler) List(c *gin.Context) {
	companyID := queryInt64(c, "company_id", 0)
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "page_size", 10)

	items, total, err := h.svc.List(c.Request.Context(), companyID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, gin.H{"list": items, "total": total})
}

// Create godoc
// @Summary Create cruise
// @Tags Cruise
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CruiseRequest true "Cruise data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/admin/cruises [post]
func (h *CruiseHandler) Create(c *gin.Context) {
	var req CruiseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	cruise := &domain.Cruise{
		CompanyID:         req.CompanyID,
		Name:              req.Name,
		EnglishName:       req.EnglishName,
		BuildYear:         req.BuildYear,
		Tonnage:           req.Tonnage,
		PassengerCapacity: req.PassengerCapacity,
		RoomCount:         req.RoomCount,
		Description:       req.Description,
		SortOrder:         req.SortOrder,
	}
	if err := h.svc.Create(c.Request.Context(), cruise); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrBadRequest, err.Error())
		return
	}
	response.Success(c, cruise)
}

// Update godoc
// @Summary Update cruise
// @Tags Cruise
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Cruise ID"
// @Param body body CruiseRequest true "Cruise data"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/cruises/{id} [put]
func (h *CruiseHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	existing, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "cruise not found")
		return
	}
	var req CruiseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	existing.Name = req.Name
	existing.EnglishName = req.EnglishName
	existing.BuildYear = req.BuildYear
	existing.Tonnage = req.Tonnage
	existing.PassengerCapacity = req.PassengerCapacity
	existing.RoomCount = req.RoomCount
	existing.Description = req.Description
	existing.SortOrder = req.SortOrder

	if err := h.svc.Update(c.Request.Context(), existing); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, existing)
}

// Delete godoc
// @Summary Delete cruise
// @Tags Cruise
// @Security BearerAuth
// @Param id path int true "Cruise ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/cruises/{id} [delete]
func (h *CruiseHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, service.ErrCruiseHasCabins) {
			response.Error(c, http.StatusConflict, errcode.ErrCruiseHasCabins, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, nil)
}
