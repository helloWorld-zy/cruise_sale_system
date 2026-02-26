package service

import (
	"errors"
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Analytics Test
type mockAnalyticsRepo struct{}

func (m mockAnalyticsRepo) TodaySales() (int64, error) { return 0, nil }

func TestAnalyticsWeeklyTrend(t *testing.T) {
	svc := NewAnalyticsService(mockAnalyticsRepo{})
	res, err := svc.WeeklyTrend()
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

// Payment Gateway Test
func TestPaymentGateways(t *testing.T) {
	wx := WechatGateway{}
	res, err := wx.CreatePay(1, 100)
	assert.NoError(t, err)
	assert.Equal(t, "wechat://pay", res)

	ali := AlipayGateway{}
	res, err = ali.CreatePay(1, 100)
	assert.NoError(t, err)
	assert.Equal(t, "alipay://pay", res)
}

// Search Retry Queue Test
func TestSearchRetryQueue(t *testing.T) {
	q := NewSearchRetryQueue(nil, 1, 1)
	q.Start()
	q.Start() // hits early return q.started == true
	q.Enqueue("test")
	q.Enqueue(make(chan int))         // should fail json encode
	time.Sleep(10 * time.Millisecond) // let worker process it briefly
}

// Hold Service Test
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
	assert.False(t, ok) // since AdjustInventoryTx returns error

	svc2 := NewCabinHoldService(&mockHoldRepo{}, 0) // hits TTL <= 0
	ok2 := svc2.HoldWithTx(nil, 0, 0, 0)            // hits early return
	assert.False(t, ok2)
}

// mockHoldRepoExistsErr returns error from ExistsActiveHoldTx
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
	q.Enqueue("test") // should not panic on nil receiver
}

func TestUserAuthNilStoreAndEmptyParams(t *testing.T) {
	// nil store
	svc := NewUserAuthService(nil)
	err := svc.SendSMS("phone", "code")
	assert.ErrorIs(t, err, ErrCodeStoreUnavailable)
	ok := svc.VerifySMS("phone", "code")
	assert.False(t, ok)

	// empty params
	svc2 := NewUserAuthService(&fakeCodeStore{})
	err2 := svc2.SendSMS("", "")
	assert.ErrorIs(t, err2, ErrPhoneOrCodeRequired)
}
