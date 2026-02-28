package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type fakeAnalyticsSvc struct {
	sales  int64
	trend  []int64
	orders int64
	err    error
}

func (f *fakeAnalyticsSvc) TodaySales(_ context.Context) (int64, error)    { return f.sales, f.err }
func (f *fakeAnalyticsSvc) WeeklyTrend(_ context.Context) ([]int64, error) { return f.trend, f.err }
func (f *fakeAnalyticsSvc) TodayOrderCount(_ context.Context) (int64, error) {
	return f.orders, f.err
}

// fakeAnalyticsSvcSelective 用于独立控制每个统计方法的错误返回。
type fakeAnalyticsSvcSelective struct {
	sales     int64
	trend     []int64
	orders    int64
	salesErr  error
	trendErr  error
	ordersErr error
}

func (f *fakeAnalyticsSvcSelective) TodaySales(_ context.Context) (int64, error) {
	return f.sales, f.salesErr
}
func (f *fakeAnalyticsSvcSelective) WeeklyTrend(_ context.Context) ([]int64, error) {
	return f.trend, f.trendErr
}
func (f *fakeAnalyticsSvcSelective) TodayOrderCount(_ context.Context) (int64, error) {
	return f.orders, f.ordersErr
}

// TestAnalyticsHandler_Summary_OK 测试获取分析摘要成功
func TestAnalyticsHandler_Summary_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &fakeAnalyticsSvc{
		sales:  50000,
		trend:  []int64{1000, 2000, 3000, 4000, 5000, 6000, 7000},
		orders: 12,
	}
	h := NewAnalyticsHandler(svc)
	r := gin.New()
	r.GET("/analytics/summary", h.Summary)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/analytics/summary", nil))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "today_sales")
	assert.Contains(t, w.Body.String(), "weekly_trend")
	assert.Contains(t, w.Body.String(), "today_orders")
}

// TestAnalyticsHandler_Summary_SalesError 测试获取分析摘要时发生销售额查询错误
func TestAnalyticsHandler_Summary_SalesError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &fakeAnalyticsSvc{err: errors.New("db error")}
	h := NewAnalyticsHandler(svc)
	r := gin.New()
	r.GET("/analytics/summary", h.Summary)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/analytics/summary", nil))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestAnalyticsHandler_Summary_WeeklyTrendError 测试周趋势查询失败分支。
func TestAnalyticsHandler_Summary_WeeklyTrendError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &fakeAnalyticsSvcSelective{
		sales:    50000,
		trendErr: errors.New("trend fail"),
	}
	h := NewAnalyticsHandler(svc)
	r := gin.New()
	r.GET("/analytics/summary", h.Summary)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/analytics/summary", nil))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestAnalyticsHandler_Summary_OrderCountError 测试今日订单数查询失败分支。
func TestAnalyticsHandler_Summary_OrderCountError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &fakeAnalyticsSvcSelective{
		sales:     50000,
		trend:     []int64{1, 2, 3, 4, 5, 6, 7},
		ordersErr: errors.New("orders fail"),
	}
	h := NewAnalyticsHandler(svc)
	r := gin.New()
	r.GET("/analytics/summary", h.Summary)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/analytics/summary", nil))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
