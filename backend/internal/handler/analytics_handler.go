package handler

import (
	"context"
	"net/http"

	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// AnalyticsSummaryService 提供仪表盘分析数据。
type AnalyticsSummaryService interface {
	TodaySales(ctx context.Context) (int64, error)
	WeeklyTrend(ctx context.Context) ([]int64, error)
	TodayOrderCount(ctx context.Context) (int64, error)
}

// AnalyticsHandler 提供仪表盘分析相关的 HTTP 端点。
type AnalyticsHandler struct{ svc AnalyticsSummaryService }

// NewAnalyticsHandler 使用给定的服务创建 AnalyticsHandler 实例。
func NewAnalyticsHandler(svc AnalyticsSummaryService) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

// Summary 处理 GET /admin/analytics/summary 请求。
// 返回今日销售总额、过去7天的趋势以及今日订单数。
func (h *AnalyticsHandler) Summary(c *gin.Context) {
	ctx := c.Request.Context()

	todaySales, err := h.svc.TodaySales(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}

	weeklyTrend, err := h.svc.WeeklyTrend(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}

	todayOrders, err := h.svc.TodayOrderCount(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, gin.H{
		"today_sales":  todaySales,
		"weekly_trend": weeklyTrend,
		"today_orders": todayOrders,
	})
}
