package v1

import (
	"cruise_booking_system/internal/core"
	"cruise_booking_system/pkg/response"

	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct {
	service *core.RecommendationService
}

func NewRecommendationHandler(service *core.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{service: service}
}

func (h *RecommendationHandler) GetRecommendations(c *gin.Context) {
	userID := c.GetString("userID") // May be empty for guest
	recs, err := h.service.GetRecommendations(c.Request.Context(), userID)
	if err != nil {
		response.ServerError(c, err)
		return
	}
	response.Success(c, recs)
}
