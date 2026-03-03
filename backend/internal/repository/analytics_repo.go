package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

type AnalyticsRepository struct{ db *gorm.DB }

func NewAnalyticsRepository(db *gorm.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) startOfTodayExpr() string {
	if r.db.Dialector.Name() == "sqlite" {
		return "date('now')"
	}
	return "CURRENT_DATE"
}

func (r *AnalyticsRepository) TodaySales(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw(
		"SELECT COALESCE(SUM(amount_cents), 0) FROM payments WHERE status = 'paid' AND created_at >= " + r.startOfTodayExpr(),
	).Scan(&total).Error
	return total, err
}

func (r *AnalyticsRepository) WeeklyTrend(ctx context.Context) ([]int64, error) {
	query := `
		SELECT COALESCE(SUM(p.amount_cents), 0)
		FROM generate_series(
			CURRENT_DATE - INTERVAL '6 days',
			CURRENT_DATE,
			'1 day'::interval
		) AS d(day)
		LEFT JOIN payments p
			ON DATE(p.created_at) = d.day AND p.status = 'paid'
		GROUP BY d.day
		ORDER BY d.day
	`
	if r.db.Dialector.Name() == "sqlite" {
		query = `
			WITH RECURSIVE days(day) AS (
				SELECT date('now', '-6 day')
				UNION ALL
				SELECT date(day, '+1 day')
				FROM days
				WHERE day < date('now')
			)
			SELECT COALESCE(SUM(p.amount_cents), 0)
			FROM days d
			LEFT JOIN payments p
				ON date(p.created_at) = d.day AND p.status = 'paid'
			GROUP BY d.day
			ORDER BY d.day
		`
	}

	rows, err := r.db.WithContext(ctx).Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trend []int64
	for rows.Next() {
		var v int64
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		trend = append(trend, v)
	}
	return trend, rows.Err()
}

func (r *AnalyticsRepository) TodayOrderCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Raw(
		"SELECT COUNT(*) FROM bookings WHERE created_at >= " + r.startOfTodayExpr(),
	).Scan(&count).Error
	return count, err
}

func (r *AnalyticsRepository) Trend(ctx context.Context, days int) ([]domain.TrendDataItem, error) {
	return []domain.TrendDataItem{}, nil
}

func (r *AnalyticsRepository) CabinHotnessRanking(ctx context.Context, limit int) ([]domain.CabinRankingItem, error) {
	return []domain.CabinRankingItem{}, nil
}

func (r *AnalyticsRepository) InventoryOverview(ctx context.Context) (*domain.InventoryOverviewData, error) {
	return &domain.InventoryOverviewData{}, nil
}

func (r *AnalyticsRepository) PageViewStats(ctx context.Context) ([]domain.PageViewData, error) {
	return []domain.PageViewData{}, nil
}
