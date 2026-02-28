package repository

import (
	"context"

	"gorm.io/gorm"
)

// AnalyticsRepository 基于 PostgreSQL 提供只读的分析查询。
//
// 性能提示：TodaySales 和 TodayOrderCount 查询按 created_at 过滤。
// 建议在生产环境中为 payments 表添加 (status, created_at) 复合索引，
// 并为 bookings 表添加 created_at 索引以优化查询性能。
type AnalyticsRepository struct{ db *gorm.DB }

// NewAnalyticsRepository 创建统计分析仓储实例。
func NewAnalyticsRepository(db *gorm.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

// startOfTodayExpr 根据数据库类型返回“今日起点”的 SQL 表达式。
func (r *AnalyticsRepository) startOfTodayExpr() string {
	if r.db.Dialector.Name() == "sqlite" {
		return "date('now')"
	}
	return "CURRENT_DATE"
}

// TodaySales 返回今日已支付的总金额（单位：分）。
func (r *AnalyticsRepository) TodaySales(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw(
		"SELECT COALESCE(SUM(amount_cents), 0) FROM payments WHERE status = 'paid' AND created_at >= " + r.startOfTodayExpr(),
	).Scan(&total).Error
	return total, err
}

// WeeklyTrend 返回过去 7 天（含今日）的每日销售总额。
// 使用 PostgreSQL 的 generate_series 填充无销售数据的日期为 0。
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
		// SQLite 使用递归 CTE 生成近 7 天日期序列，保证测试环境可运行。
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

// TodayOrderCount 返回今日创建的预订订单数量。
func (r *AnalyticsRepository) TodayOrderCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Raw(
		"SELECT COUNT(*) FROM bookings WHERE created_at >= " + r.startOfTodayExpr(),
	).Scan(&count).Error
	return count, err
}
