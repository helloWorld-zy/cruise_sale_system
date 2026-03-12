package handler

import (
	"context"

	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type PortCitySearchService interface {
	Search(ctx context.Context, keyword string) ([]service.PortCityOption, error)
}

type PortCityHandler struct {
	svc PortCitySearchService
}

func NewPortCityHandler(svc PortCitySearchService) *PortCityHandler {
	return &PortCityHandler{svc: svc}
}

func (h *PortCityHandler) Search(c *gin.Context) {
	items, err := h.svc.Search(c.Request.Context(), c.Query("keyword"))
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, items)
}
