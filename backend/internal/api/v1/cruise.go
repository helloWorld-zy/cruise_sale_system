package v1

import (
	"cruise_booking_system/internal/core"
	"cruise_booking_system/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CruiseHandler struct {
	service *core.CruiseService
}

func NewCruiseHandler(service *core.CruiseService) *CruiseHandler {
	return &CruiseHandler{service: service}
}

func (h *CruiseHandler) List(c *gin.Context) {
	destination := c.Query("destination")
	date := c.Query("date")

	cruises, err := h.service.ListCruises(c.Request.Context(), destination, date)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, cruises)
}

func (h *CruiseHandler) GetDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, 400, "id is required")
		return
	}

	detail, err := h.service.GetCruiseDetail(c.Request.Context(), id)
	if err != nil {
		// Differentiate not found vs internal error if needed
		response.ServerError(c, err)
		return
	}

	response.Success(c, detail)
}
