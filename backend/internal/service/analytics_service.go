package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

// AnalyticsService 通过委托给仓储提供仪表盘分析。
type AnalyticsService struct {
	repo domain.AnalyticsRepository
}

// NewAnalyticsService 创建一个新的 AnalyticsService。
func NewAnalyticsService(repo domain.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

// TodaySales 返回今天的支付总额（单位：分）。
func (s *AnalyticsService) TodaySales(ctx context.Context) (int64, error) {
	return s.repo.TodaySales(ctx)
}

// WeeklyTrend 返回过去 7 天的每日销售总额。
func (s *AnalyticsService) WeeklyTrend(ctx context.Context) ([]int64, error) {
	return s.repo.WeeklyTrend(ctx)
}

// TodayOrderCount 返回今天创建的预订数量。
func (s *AnalyticsService) TodayOrderCount(ctx context.Context) (int64, error) {
	return s.repo.TodayOrderCount(ctx)
}
