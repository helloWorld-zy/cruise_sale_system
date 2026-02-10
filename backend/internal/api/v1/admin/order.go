package admin

import (
	"cruise_booking_system/internal/core"
	"cruise_booking_system/internal/model"
	"cruise_booking_system/pkg/response"
	"cruise_booking_system/pkg/storage"
	"fmt"
	"net/http"

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

func (h *AdminOrderHandler) UploadDepartureNotice(c *gin.Context) {
	id := c.Param("id")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "file is required")
		return
	}
	defer file.Close()

	// Bucket: departure-notices, Object: orderID.pdf (or timestamped)
	objectName := fmt.Sprintf("%s_%s", id, header.Filename)
	url, err := storage.UploadFile(c.Request.Context(), "departure-notices", objectName, file, header.Size, header.Header.Get("Content-Type"))
	if err != nil {
		response.ServerError(c, err)
		return
	}

	if err := h.service.SetDepartureNotice(c.Request.Context(), id, url); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"url": url})
}
