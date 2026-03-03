package repository

import (
	"context"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newAnalyticsTestRepo 创建 SQLite 内存库并返回分析仓储实例。
func newAnalyticsTestRepo(t *testing.T) *AnalyticsRepository {
	t.Helper()
	// 使用独立内存库，避免并发测试之间相互污染。
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	// 仅迁移测试所需的最小模型。
	require.NoError(t, db.AutoMigrate(&domain.Payment{}, &domain.Booking{}))
	return NewAnalyticsRepository(db)
}

func TestAnalyticsRepository_TodaySales_NoPayments(t *testing.T) {
	repo := newAnalyticsTestRepo(t)
	total, err := repo.TodaySales(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
}

func TestAnalyticsRepository_TodaySales_WithPaidPayments(t *testing.T) {
	repo := newAnalyticsTestRepo(t)
	db := repo.db

	// 插入今日已支付与非已支付数据，验证仅统计 paid。
	require.NoError(t, db.Exec("INSERT INTO payments (order_id, provider, trade_no, amount_cents, status, created_at, updated_at) VALUES (1, 'wechat', 'T1', 1000, 'paid', datetime('now'), datetime('now'))").Error)
	require.NoError(t, db.Exec("INSERT INTO payments (order_id, provider, trade_no, amount_cents, status, created_at, updated_at) VALUES (2, 'wechat', 'T2', 2000, 'paid', datetime('now'), datetime('now'))").Error)
	require.NoError(t, db.Exec("INSERT INTO payments (order_id, provider, trade_no, amount_cents, status, created_at, updated_at) VALUES (3, 'wechat', 'T3', 5000, 'failed', datetime('now'), datetime('now'))").Error)

	total, err := repo.TodaySales(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(3000), total)
}

func TestAnalyticsRepository_TodayOrderCount_NoBookings(t *testing.T) {
	repo := newAnalyticsTestRepo(t)
	count, err := repo.TodayOrderCount(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestAnalyticsRepository_TodayOrderCount_WithBookings(t *testing.T) {
	repo := newAnalyticsTestRepo(t)
	db := repo.db

	// 仅统计今天创建的订单。
	require.NoError(t, db.Exec("INSERT INTO bookings (user_id, voyage_id, cabin_sku_id, status, total_cents, created_at, updated_at) VALUES (1, 11, 101, 'created', 8000, datetime('now'), datetime('now'))").Error)
	require.NoError(t, db.Exec("INSERT INTO bookings (user_id, voyage_id, cabin_sku_id, status, total_cents, created_at, updated_at) VALUES (2, 12, 102, 'created', 9000, datetime('now'), datetime('now'))").Error)
	require.NoError(t, db.Exec("INSERT INTO bookings (user_id, voyage_id, cabin_sku_id, status, total_cents, created_at, updated_at) VALUES (3, 13, 103, 'created', 10000, datetime('now', '-1 day'), datetime('now', '-1 day'))").Error)

	count, err := repo.TodayOrderCount(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestAnalyticsRepository_WeeklyTrend(t *testing.T) {
	repo := newAnalyticsTestRepo(t)
	db := repo.db

	// 构造近 7 天中的有效支付与无效支付数据。
	require.NoError(t, db.Exec("INSERT INTO payments (order_id, provider, trade_no, amount_cents, status, created_at, updated_at) VALUES (1, 'wechat', 'W1', 200, 'paid', datetime('now', '-3 day'), datetime('now', '-3 day'))").Error)
	require.NoError(t, db.Exec("INSERT INTO payments (order_id, provider, trade_no, amount_cents, status, created_at, updated_at) VALUES (2, 'wechat', 'W2', 250, 'paid', datetime('now'), datetime('now'))").Error)
	require.NoError(t, db.Exec("INSERT INTO payments (order_id, provider, trade_no, amount_cents, status, created_at, updated_at) VALUES (3, 'wechat', 'W3', 999, 'failed', datetime('now'), datetime('now'))").Error)

	trend, err := repo.WeeklyTrend(context.Background())
	require.NoError(t, err)
	assert.Len(t, trend, 7)

	var sum int64
	for _, v := range trend {
		sum += v
	}
	assert.Equal(t, int64(450), sum)
}

func TestAnalyticsRepository_Trend(t *testing.T) {
	repo := newAnalyticsTestRepo(t)
	db := repo.db

	require.NoError(t, db.Exec("INSERT INTO payments (order_id, provider, trade_no, amount_cents, status, created_at, updated_at) VALUES (10, 'wechat', 'TR1', 1200, 'paid', datetime('now'), datetime('now'))").Error)
	require.NoError(t, db.Exec("INSERT INTO bookings (user_id, voyage_id, cabin_sku_id, status, total_cents, created_at, updated_at) VALUES (1, 11, 201, 'created', 1200, datetime('now'), datetime('now'))").Error)

	items, err := repo.Trend(context.Background(), 7)
	require.NoError(t, err)
	assert.Len(t, items, 7)

	var sales int64
	var orders int64
	for _, item := range items {
		sales += item.Sales
		orders += item.Orders
	}
	assert.Equal(t, int64(1200), sales)
	assert.Equal(t, int64(1), orders)
}

func TestAnalyticsRepository_CabinHotnessRanking(t *testing.T) {
	repo := newAnalyticsTestRepo(t)
	db := repo.db

	require.NoError(t, db.AutoMigrate(&domain.CabinSKU{}))
	require.NoError(t, db.Exec("INSERT INTO cabin_skus (id, code, voyage_id, status, created_at, updated_at) VALUES (301, 'C301', 11, 1, datetime('now'), datetime('now'))").Error)
	require.NoError(t, db.Exec("INSERT INTO cabin_skus (id, code, voyage_id, status, created_at, updated_at) VALUES (302, 'C302', 11, 1, datetime('now'), datetime('now'))").Error)
	require.NoError(t, db.Exec("INSERT INTO bookings (user_id, voyage_id, cabin_sku_id, status, total_cents, created_at, updated_at) VALUES (1, 11, 301, 'paid', 1000, datetime('now'), datetime('now'))").Error)
	require.NoError(t, db.Exec("INSERT INTO bookings (user_id, voyage_id, cabin_sku_id, status, total_cents, created_at, updated_at) VALUES (2, 11, 301, 'confirmed', 1000, datetime('now'), datetime('now'))").Error)
	require.NoError(t, db.Exec("INSERT INTO bookings (user_id, voyage_id, cabin_sku_id, status, total_cents, created_at, updated_at) VALUES (3, 11, 302, 'created', 1000, datetime('now'), datetime('now'))").Error)

	ranking, err := repo.CabinHotnessRanking(context.Background(), 10)
	require.NoError(t, err)
	require.NotEmpty(t, ranking)
	assert.Equal(t, int64(301), ranking[0].CabinSKUID)
	assert.Equal(t, int64(2), ranking[0].SoldCount)
}

func TestAnalyticsRepository_InventoryOverviewAndPageViewStats(t *testing.T) {
	repo := newAnalyticsTestRepo(t)
	db := repo.db

	require.NoError(t, db.AutoMigrate(&domain.CabinInventory{}))
	require.NoError(t, db.Exec("INSERT INTO cabin_inventories (cabin_sku_id, total, locked, sold, alert_threshold) VALUES (401, 10, 3, 6, 2)").Error)
	require.NoError(t, db.Exec("INSERT INTO cabin_inventories (cabin_sku_id, total, locked, sold, alert_threshold) VALUES (402, 5, 2, 3, 1)").Error)
	require.NoError(t, db.Exec("INSERT INTO bookings (user_id, voyage_id, cabin_sku_id, status, total_cents, created_at, updated_at) VALUES (5, 15, 401, 'created', 1000, datetime('now'), datetime('now'))").Error)
	require.NoError(t, db.Exec("INSERT INTO payments (order_id, provider, trade_no, amount_cents, status, created_at, updated_at) VALUES (99, 'wechat', 'PV1', 2000, 'paid', datetime('now'), datetime('now'))").Error)

	overview, err := repo.InventoryOverview(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(2), overview.TotalCabins)
	assert.Equal(t, int64(1), overview.OutOfStockCount)

	stats, err := repo.PageViewStats(context.Background())
	require.NoError(t, err)
	assert.Len(t, stats, 3)
}
