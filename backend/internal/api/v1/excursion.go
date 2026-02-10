package v1

import (
	"cruise_booking_system/internal/core"
	"cruise_booking_system/pkg/response"

	"github.com/gin-gonic/gin"
)

type ExcursionHandler struct {
	service *core.ExcursionService
}

func NewExcursionHandler(service *core.ExcursionService) *ExcursionHandler {
	return &ExcursionHandler{service: service}
}

func (h *ExcursionHandler) List(c *gin.Context) {
	cruiseID := c.Query("cruise_id")
	excursions, err := h.service.ListExcursions(c.Request.Context(), cruiseID)
	if err != nil {
		response.ServerError(c, err)
		return
	}
	response.Success(c, excursions)
}
