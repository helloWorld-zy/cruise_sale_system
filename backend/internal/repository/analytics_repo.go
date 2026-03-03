package repository

import (
	"context"
	"fmt"

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
	if days <= 0 {
		days = 7
	}
	if days > 90 {
		days = 90
	}

	query := fmt.Sprintf(`
		WITH RECURSIVE days(day) AS (
			SELECT date('now', '-%d day')
			UNION ALL
			SELECT date(day, '+1 day') FROM days WHERE day < date('now')
		)
		SELECT
			d.day AS date,
			COALESCE(SUM(CASE WHEN p.status = 'paid' THEN p.amount_cents ELSE 0 END), 0) AS sales,
			COALESCE(COUNT(b.id), 0) AS orders
		FROM days d
		LEFT JOIN payments p ON date(p.created_at) = d.day
		LEFT JOIN bookings b ON date(b.created_at) = d.day
		GROUP BY d.day
		ORDER BY d.day
	`, days-1)

	if r.db.Dialector.Name() != "sqlite" {
		query = fmt.Sprintf(`
			WITH days AS (
				SELECT generate_series(CURRENT_DATE - INTERVAL '%d day', CURRENT_DATE, INTERVAL '1 day')::date AS day
			)
			SELECT
				to_char(d.day, 'YYYY-MM-DD') AS date,
				COALESCE(SUM(CASE WHEN p.status = 'paid' THEN p.amount_cents ELSE 0 END), 0) AS sales,
				COALESCE(COUNT(b.id), 0) AS orders
			FROM days d
			LEFT JOIN payments p ON DATE(p.created_at) = d.day
			LEFT JOIN bookings b ON DATE(b.created_at) = d.day
			GROUP BY d.day
			ORDER BY d.day
		`, days-1)
	}

	rows, err := r.db.WithContext(ctx).Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]domain.TrendDataItem, 0)
	for rows.Next() {
		var item domain.TrendDataItem
		if err := rows.Scan(&item.Date, &item.Sales, &item.Orders); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *AnalyticsRepository) CabinHotnessRanking(ctx context.Context, limit int) ([]domain.CabinRankingItem, error) {
	if limit <= 0 {
		limit = 10
	}
	query := `
		SELECT
			b.cabin_sku_id,
			COALESCE(cs.code, '') AS cabin_name,
			SUM(CASE WHEN b.status IN ('paid', 'confirmed', 'pending_travel', 'traveling', 'completed') THEN 1 ELSE 0 END) AS sold_count,
			COUNT(*) AS view_count
		FROM bookings b
		LEFT JOIN cabin_skus cs ON cs.id = b.cabin_sku_id
		GROUP BY b.cabin_sku_id, cs.code
		ORDER BY sold_count DESC, view_count DESC, b.cabin_sku_id ASC
		LIMIT ?
	`
	rows, err := r.db.WithContext(ctx).Raw(query, limit).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]domain.CabinRankingItem, 0)
	for rows.Next() {
		var item domain.CabinRankingItem
		if err := rows.Scan(&item.CabinSKUID, &item.CabinName, &item.SoldCount, &item.ViewCount); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *AnalyticsRepository) InventoryOverview(ctx context.Context) (*domain.InventoryOverviewData, error) {
	query := `
		SELECT
			COUNT(*) AS total_cabins,
			SUM(CASE WHEN (total - locked - sold) <= alert_threshold AND alert_threshold > 0 THEN 1 ELSE 0 END) AS low_stock_count,
			SUM(CASE WHEN (total - locked - sold) <= 0 THEN 1 ELSE 0 END) AS out_of_stock_count
		FROM cabin_inventories
	`
	data := &domain.InventoryOverviewData{}
	if err := r.db.WithContext(ctx).Raw(query).Scan(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *AnalyticsRepository) PageViewStats(ctx context.Context) ([]domain.PageViewData, error) {
	query := `
		SELECT '/bookings' AS page, COUNT(*) AS views FROM bookings
		UNION ALL
		SELECT '/payments' AS page, COUNT(*) AS views FROM payments
		UNION ALL
		SELECT '/cabins' AS page, COUNT(DISTINCT cabin_sku_id) AS views FROM bookings
	`
	rows, err := r.db.WithContext(ctx).Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]domain.PageViewData, 0)
	for rows.Next() {
		var item domain.PageViewData
		if err := rows.Scan(&item.Page, &item.Views); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}
