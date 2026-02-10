package admin

import (
	"cruise_booking_system/internal/core"
	"cruise_booking_system/internal/model"
	"cruise_booking_system/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminOrderHandler struct {
	service *core.OrderService
}

func NewAdminOrderHandler(service *core.OrderService) *AdminOrderHandler {
	return &AdminOrderHandler{service: service}
}

func (h *AdminOrderHandler) List(c *gin.Context) {
	status := c.Query("status")
	orders, err := h.service.ListAllOrders(c.Request.Context(), model.OrderStatus(status))
	if err != nil {
		response.ServerError(c, err)
		return
	}
	response.Success(c, orders)
}
