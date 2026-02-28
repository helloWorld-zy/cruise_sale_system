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

// newRefundTestRepo 创建退款仓储测试实例。
func newRefundTestRepo(t *testing.T) (*RefundRepository, *PaymentRepository) {
	t.Helper()
	// 使用独立内存库，避免测试数据串扰。
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	// 同时迁移 Payment，避免 Refund 的外键引用缺失。
	require.NoError(t, db.AutoMigrate(&domain.Payment{}, &domain.Refund{}))
	return NewRefundRepository(db), NewPaymentRepository(db)
}

func TestRefundRepository_Create(t *testing.T) {
	refundRepo, payRepo := newRefundTestRepo(t)
	payment := &domain.Payment{OrderID: 88, Provider: "wechat", TradeNo: "R-C1", AmountCents: 10000, Status: "paid"}
	require.NoError(t, payRepo.Create(context.Background(), payment))

	r := &domain.Refund{PaymentID: payment.ID, AmountCents: 500, Reason: "test", Status: "pending"}
	require.NoError(t, refundRepo.Create(context.Background(), r))
	assert.Greater(t, r.ID, int64(0))
}

func TestRefundRepository_SumByPaymentID_NoRefunds(t *testing.T) {
	refundRepo, payRepo := newRefundTestRepo(t)
	payment := &domain.Payment{OrderID: 88, Provider: "wechat", TradeNo: "R-NO", AmountCents: 10000, Status: "paid"}
	require.NoError(t, payRepo.Create(context.Background(), payment))

	total, err := refundRepo.SumByPaymentID(context.Background(), payment.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
}

func TestRefundRepository_SumByPaymentID_WithRefunds(t *testing.T) {
	refundRepo, payRepo := newRefundTestRepo(t)
	payment := &domain.Payment{OrderID: 88, Provider: "wechat", TradeNo: "R-SUM", AmountCents: 10000, Status: "paid"}
	require.NoError(t, payRepo.Create(context.Background(), payment))

	require.NoError(t, refundRepo.Create(context.Background(), &domain.Refund{PaymentID: payment.ID, AmountCents: 1000, Reason: "a", Status: "pending"}))
	require.NoError(t, refundRepo.Create(context.Background(), &domain.Refund{PaymentID: payment.ID, AmountCents: 2300, Reason: "b", Status: "approved"}))

	total, err := refundRepo.SumByPaymentID(context.Background(), payment.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(3300), total)
}

func TestRefundRepository_SumByPaymentID_ExcludesCancelled(t *testing.T) {
	refundRepo, payRepo := newRefundTestRepo(t)
	payment := &domain.Payment{OrderID: 88, Provider: "wechat", TradeNo: "R-CAN", AmountCents: 10000, Status: "paid"}
	require.NoError(t, payRepo.Create(context.Background(), payment))

	require.NoError(t, refundRepo.Create(context.Background(), &domain.Refund{PaymentID: payment.ID, AmountCents: 1000, Reason: "a", Status: "pending"}))
	require.NoError(t, refundRepo.Create(context.Background(), &domain.Refund{PaymentID: payment.ID, AmountCents: 5000, Reason: "b", Status: "cancelled"}))

	total, err := refundRepo.SumByPaymentID(context.Background(), payment.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1000), total)
}
