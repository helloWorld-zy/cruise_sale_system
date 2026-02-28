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

// newPaymentTestRepo 创建支付仓储测试实例。
func newPaymentTestRepo(t *testing.T) *PaymentRepository {
	t.Helper()
	// 使用独立内存库，避免测试数据串扰。
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&domain.Payment{}))
	return NewPaymentRepository(db)
}

func TestPaymentRepository_Create(t *testing.T) {
	repo := newPaymentTestRepo(t)
	p := &domain.Payment{OrderID: 101, Provider: "wechat", TradeNo: "TN-1", AmountCents: 999, Status: "pending"}
	require.NoError(t, repo.Create(context.Background(), p))
	assert.Greater(t, p.ID, int64(0))
}

func TestPaymentRepository_FindByTradeNo_Found(t *testing.T) {
	repo := newPaymentTestRepo(t)
	seed := &domain.Payment{OrderID: 101, Provider: "wechat", TradeNo: "TN-FOUND", AmountCents: 999, Status: "pending"}
	require.NoError(t, repo.Create(context.Background(), seed))

	p, err := repo.FindByTradeNo(context.Background(), "TN-FOUND")
	require.NoError(t, err)
	assert.Equal(t, seed.ID, p.ID)
}

func TestPaymentRepository_FindByTradeNo_NotFound(t *testing.T) {
	repo := newPaymentTestRepo(t)
	_, err := repo.FindByTradeNo(context.Background(), "NOT-EXIST")
	assert.Error(t, err)
}

func TestPaymentRepository_FindByID_Found(t *testing.T) {
	repo := newPaymentTestRepo(t)
	seed := &domain.Payment{OrderID: 101, Provider: "wechat", TradeNo: "TN-ID", AmountCents: 999, Status: "pending"}
	require.NoError(t, repo.Create(context.Background(), seed))

	p, err := repo.FindByID(context.Background(), seed.ID)
	require.NoError(t, err)
	assert.Equal(t, "TN-ID", p.TradeNo)
}

func TestPaymentRepository_FindByID_NotFound(t *testing.T) {
	repo := newPaymentTestRepo(t)
	_, err := repo.FindByID(context.Background(), 99999)
	assert.Error(t, err)
}

func TestPaymentRepository_UpdateStatus(t *testing.T) {
	repo := newPaymentTestRepo(t)
	seed := &domain.Payment{OrderID: 101, Provider: "wechat", TradeNo: "TN-UPD", AmountCents: 999, Status: "pending"}
	require.NoError(t, repo.Create(context.Background(), seed))

	require.NoError(t, repo.UpdateStatus(context.Background(), seed.ID, "paid"))
	got, err := repo.FindByID(context.Background(), seed.ID)
	require.NoError(t, err)
	assert.Equal(t, "paid", got.Status)
}
