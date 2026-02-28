package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// 分析测试 — 使用更新的 domain.AnalyticsRepository 接口（带 context）。
type mockAnalyticsRepo struct{}

func (m mockAnalyticsRepo) TodaySales(_ context.Context) (int64, error) { return 0, nil }
func (m mockAnalyticsRepo) WeeklyTrend(_ context.Context) ([]int64, error) {
	return []int64{1, 2, 3, 4, 5, 6, 7}, nil
}
func (m mockAnalyticsRepo) TodayOrderCount(_ context.Context) (int64, error) { return 0, nil }

func TestAnalyticsWeeklyTrendEdge(t *testing.T) {
	svc := NewAnalyticsService(mockAnalyticsRepo{})
	res, err := svc.WeeklyTrend(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

// 支付网关测试 — 网关现在返回 (tradeNo, payURL, error)。
func TestPaymentGateways(t *testing.T) {
	wx := WechatGateway{}
	tradeNo, payURL, err := wx.CreatePay(1, 100)
	assert.NoError(t, err)
	assert.NotEmpty(t, tradeNo)
	assert.Equal(t, "wechat://pay", payURL)

	ali := AlipayGateway{}
	tradeNo, payURL, err = ali.CreatePay(1, 100)
	assert.NoError(t, err)
	assert.NotEmpty(t, tradeNo)
	assert.Equal(t, "alipay://pay", payURL)
}

// 搜索重试队列测试
func TestSearchRetryQueue(t *testing.T) {
	q := NewSearchRetryQueue(nil, 1, 1)
	q.Start()
	q.Start() // 命中提前返回 q.started == true
	q.Enqueue("test")
	q.Enqueue(make(chan int))         // 应当在 json 编码时失败
	time.Sleep(10 * time.Millisecond) // 让工作协程短暂处理它
}

// 占座服务测试
type mockHoldRepo struct {
	called bool
}

func (m *mockHoldRepo) CreateHoldTx(tx *gorm.DB, hold *domain.CabinHold) error {
	return nil
}
func (m *mockHoldRepo) AdjustInventoryTx(tx *gorm.DB, skuID int64, delta int, reason string) error {
	return errors.New("err")
}
func (m *mockHoldRepo) ExistsActiveHoldTx(tx *gorm.DB, skuID int64, userID int64, now time.Time) (bool, error) {
	return false, nil
}

func TestHoldService(t *testing.T) {
	svc := NewCabinHoldService(&mockHoldRepo{}, time.Second)
	ok := svc.HoldWithTx(nil, 1, 1, 1)
	assert.False(t, ok) // 因为 AdjustInventoryTx 返回错误

	svc2 := NewCabinHoldService(&mockHoldRepo{}, 0) // 命中 TTL <= 0
	ok2 := svc2.HoldWithTx(nil, 0, 0, 0)            // 命中提前返回
	assert.False(t, ok2)
}

// mockHoldRepoExistsErr 从 ExistsActiveHoldTx 返回错误
type mockHoldRepoExistsErr struct{}

func (m *mockHoldRepoExistsErr) ExistsActiveHoldTx(_ *gorm.DB, _ int64, _ int64, _ time.Time) (bool, error) {
	return false, errors.New("exists error")
}
func (m *mockHoldRepoExistsErr) CreateHoldTx(_ *gorm.DB, _ *domain.CabinHold) error { return nil }
func (m *mockHoldRepoExistsErr) AdjustInventoryTx(_ *gorm.DB, _ int64, _ int, _ string) error {
	return nil
}

func TestHoldServiceExistsError(t *testing.T) {
	svc := NewCabinHoldService(&mockHoldRepoExistsErr{}, time.Second)
	ok := svc.HoldWithTx(nil, 1, 1, 1)
	assert.False(t, ok)
}

func TestSearchRetryQueueDefaults(t *testing.T) {
	q := NewSearchRetryQueue(nil, 0, 0)
	assert.NotNil(t, q)
}

func TestSearchRetryQueueNilEnqueue(t *testing.T) {
	var q *SearchRetryQueue
	q.Enqueue("test") // 接收者为 nil 时不应 panic
}

func TestUserAuthNilStoreAndEmptyParams(t *testing.T) {
	// 仓储为 nil
	svc := NewUserAuthService(nil)
	err := svc.SendSMS("phone", "code")
	assert.ErrorIs(t, err, ErrCodeStoreUnavailable)
	ok := svc.VerifySMS("phone", "code")
	assert.False(t, ok)

	// 空参数
	svc2 := NewUserAuthService(&fakeCodeStore{})
	err2 := svc2.SendSMS("", "")
	assert.ErrorIs(t, err2, ErrPhoneOrCodeRequired)
}
