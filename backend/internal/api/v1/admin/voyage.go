package admin

import (
	"cruise_booking_system/internal/core"
	"cruise_booking_system/internal/model"
	"cruise_booking_system/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminVoyageHandler struct {
	voyageSvc    *core.VoyageService
	inventorySvc *core.InventoryService
}

func NewAdminVoyageHandler(v *core.VoyageService, i *core.InventoryService) *AdminVoyageHandler {
	return &AdminVoyageHandler{
		voyageSvc:    v,
		inventorySvc: i,
	}
}

type CreateVoyageRequest struct {
	CruiseID      string    `json:"cruise_id"`
	DepartureDate time.Time `json:"departure_date"`
	ReturnDate    time.Time `json:"return_date"`
}

func (h *AdminVoyageHandler) Create(c *gin.Context) {
	var req CreateVoyageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	cid, _ := uuid.Parse(req.CruiseID)
	voyage := &model.Voyage{
		CruiseID:      cid,
		DepartureDate: req.DepartureDate,
		ReturnDate:    req.ReturnDate,
	}

	if err := h.voyageSvc.CreateVoyage(c.Request.Context(), voyage); err != nil {
		response.ServerError(c, err)
		return
	}
	response.Success(c, voyage.ID)
}
