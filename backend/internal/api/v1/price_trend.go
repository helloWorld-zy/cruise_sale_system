package v1

import (
	"cruise_booking_system/internal/core"
	"cruise_booking_system/pkg/response"

	"github.com/gin-gonic/gin"
)

type TrendHandler struct {
	service *core.TrendService
}

func NewTrendHandler(service *core.TrendService) *TrendHandler {
	return &TrendHandler{service: service}
}

func (h *TrendHandler) GetTrends(c *gin.Context) {
	voyageID := c.Query("voyage_id")
	trends, err := h.service.GetPriceTrends(c.Request.Context(), voyageID)
	if err != nil {
		response.ServerError(c, err)
		return
	}
	response.Success(c, trends)
}
