package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

type TrendData struct {
	Date   string
	Sales  int64
	Orders int64
}

type CabinRankingItem struct {
	CabinSKUID int64
	CabinName  string
	SoldCount  int64
	ViewCount  int64
}

type InventoryOverview struct {
	TotalCabins     int64
	LowStockCount   int64
	OutOfStockCount int64
}

type PageViewStat struct {
	Page  string
	Views int64
}

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

// Trend 返回过去 N 天的趋势数据。
func (s *AnalyticsService) Trend(ctx context.Context, days int) ([]TrendData, error) {
	trends, err := s.repo.Trend(ctx, days)
	if err != nil {
		return nil, err
	}
	result := make([]TrendData, len(trends))
	for i, t := range trends {
		result[i] = TrendData{
			Date:   t.Date,
			Sales:  t.Sales,
			Orders: t.Orders,
		}
	}
	return result, nil
}

// CabinHotnessRanking 返回舱位热度排行。
func (s *AnalyticsService) CabinHotnessRanking(ctx context.Context, limit int) ([]CabinRankingItem, error) {
	items, err := s.repo.CabinHotnessRanking(ctx, limit)
	if err != nil {
		return nil, err
	}
	result := make([]CabinRankingItem, len(items))
	for i, item := range items {
		result[i] = CabinRankingItem{
			CabinSKUID: item.CabinSKUID,
			CabinName:  item.CabinName,
			SoldCount:  item.SoldCount,
			ViewCount:  item.ViewCount,
		}
	}
	return result, nil
}

// InventoryOverview 返回库存概览。
func (s *AnalyticsService) InventoryOverview(ctx context.Context) (*domain.InventoryOverviewData, error) {
	return s.repo.InventoryOverview(ctx)
}

// PageViewStats 返回页面访问量统计。
func (s *AnalyticsService) PageViewStats(ctx context.Context) ([]domain.PageViewData, error) {
	return s.repo.PageViewStats(ctx)
}
