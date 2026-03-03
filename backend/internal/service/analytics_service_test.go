package service

import (
	"context"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeAnalyticsRepo struct{}

func (f fakeAnalyticsRepo) TodaySales(_ context.Context) (int64, error) { return 1000, nil }
func (f fakeAnalyticsRepo) WeeklyTrend(_ context.Context) ([]int64, error) {
	return []int64{100, 200, 300, 400, 500, 600, 700}, nil
}
func (f fakeAnalyticsRepo) TodayOrderCount(_ context.Context) (int64, error) { return 5, nil }
func (f fakeAnalyticsRepo) Trend(_ context.Context, _ int) ([]domain.TrendDataItem, error) {
	return nil, nil
}
func (f fakeAnalyticsRepo) CabinHotnessRanking(_ context.Context, _ int) ([]domain.CabinRankingItem, error) {
	return nil, nil
}
func (f fakeAnalyticsRepo) InventoryOverview(_ context.Context) (*domain.InventoryOverviewData, error) {
	return nil, nil
}
func (f fakeAnalyticsRepo) PageViewStats(_ context.Context) ([]domain.PageViewData, error) {
	return nil, nil
}

func TestAnalyticsTodaySales(t *testing.T) {
	svc := NewAnalyticsService(fakeAnalyticsRepo{})
	v, err := svc.TodaySales(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(1000), v)
}

func TestAnalyticsWeeklyTrend(t *testing.T) {
	svc := NewAnalyticsService(fakeAnalyticsRepo{})
	res, err := svc.WeeklyTrend(context.Background())
	require.NoError(t, err)
	assert.Len(t, res, 7)
}

func TestAnalyticsTodayOrderCount(t *testing.T) {
	svc := NewAnalyticsService(fakeAnalyticsRepo{})
	count, err := svc.TodayOrderCount(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
}
