package handler

import (
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct{}

func NewAnalyticsHandler() *AnalyticsHandler { return &AnalyticsHandler{} }

func (h *AnalyticsHandler) Summary(c *gin.Context) {
	response.Success(c, gin.H{"today_sales": 1000, "weekly_trend": []int64{1, 2, 3}})
}
