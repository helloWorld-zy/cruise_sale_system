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

// CompanyHandler handles CRUD for cruise companies.
// CR-05: Handler now has proper DI via CompanyService.
type CompanyHandler struct {
	svc *service.CompanyService
}

// NewCompanyHandler creates a CompanyHandler with injected service.
func NewCompanyHandler(svc *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{svc: svc}
}

// CompanyRequest is the create/update payload.
type CompanyRequest struct {
	Name        string `json:"name" binding:"required"`
	EnglishName string `json:"english_name"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
	SortOrder   int    `json:"sort_order"`
}

// List godoc
// @Summary List cruise companies
// @Tags Company
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Search keyword"
// @Param page query int false "Page" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/admin/companies [get]
func (h *CompanyHandler) List(c *gin.Context) {
	keyword := c.Query("keyword")
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "page_size", 10)

	items, total, err := h.svc.List(c.Request.Context(), keyword, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, gin.H{"list": items, "total": total})
}

// Create godoc
// @Summary Create a cruise company
// @Tags Company
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CompanyRequest true "Company data"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/companies [post]
func (h *CompanyHandler) Create(c *gin.Context) {
	var req CompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	company := &domain.CruiseCompany{
		Name:        req.Name,
		EnglishName: req.EnglishName,
		Description: req.Description,
		LogoURL:     req.LogoURL,
		SortOrder:   req.SortOrder,
	}
	if err := h.svc.Create(c.Request.Context(), company); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, company)
}

// Update godoc
// @Summary Update a cruise company
// @Tags Company
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Company ID"
// @Param body body CompanyRequest true "Company data"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/companies/{id} [put]
func (h *CompanyHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	existing, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "company not found")
		return
	}
	var req CompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	existing.Name = req.Name
	existing.EnglishName = req.EnglishName
	existing.Description = req.Description
	existing.LogoURL = req.LogoURL
	existing.SortOrder = req.SortOrder
	if err := h.svc.Update(c.Request.Context(), existing); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, existing)
}

// Delete godoc
// @Summary Delete a cruise company
// @Tags Company
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/companies/{id} [delete]
func (h *CompanyHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if err.Error() == "company has cruises" {
			response.Error(c, http.StatusConflict, errcode.ErrCompanyHasCruises, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, nil)
}
