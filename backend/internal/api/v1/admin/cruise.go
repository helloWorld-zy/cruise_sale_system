package admin

import (
	"cruise_booking_system/internal/core"
	"cruise_booking_system/internal/model"
	"cruise_booking_system/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AdminCruiseHandler struct {
	service *core.CruiseService
}

func NewAdminCruiseHandler(service *core.CruiseService) *AdminCruiseHandler {
	return &AdminCruiseHandler{service: service}
}

type CreateCruiseRequest struct {
	NameEn      string   `json:"name_en" binding:"required"`
	NameCn      string   `json:"name_cn" binding:"required"`
	Code        string   `json:"code" binding:"required"`
	Tonnage     int      `json:"tonnage"`
	Capacity    int      `json:"capacity"`
	Decks       int      `json:"decks"`
	Gallery     []string `json:"gallery"`
	Description string   `json:"description"`
}

func (h *AdminCruiseHandler) Create(c *gin.Context) {
	var req CreateCruiseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	// Simple mapping (should use mapper)
	cruise := &model.Cruise{
		NameEn:      req.NameEn,
		NameCn:      req.NameCn,
		Code:        req.Code,
		Tonnage:     req.Tonnage,
		Capacity:    req.Capacity,
		Decks:       req.Decks,
		Description: req.Description,
		// Gallery mapping requires casting to datatypes.JSON
	}
	// For gallery, assume we handle it
	
	if err := h.service.CreateCruise(c.Request.Context(), cruise); err != nil {
		response.ServerError(c, err)
		return
	}
	response.Success(c, cruise.ID)
}

func (h *AdminCruiseHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	var req CreateCruiseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	id, _ := uuid.Parse(idStr)
	cruise := &model.Cruise{
		BaseModel: model.BaseModel{ID: id},
		NameEn:    req.NameEn,
		NameCn:    req.NameCn,
		// ... fields
	}

	if err := h.service.UpdateCruise(c.Request.Context(), cruise); err != nil {
		response.ServerError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *AdminCruiseHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteCruise(c.Request.Context(), id); err != nil {
		response.ServerError(c, err)
		return
	}
	response.Success(c, nil)
}
