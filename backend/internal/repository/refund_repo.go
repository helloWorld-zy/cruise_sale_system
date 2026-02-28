package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// RefundRepository 基于 PostgreSQL 提供退款记录的持久化操作。
type RefundRepository struct{ db *gorm.DB }

// NewRefundRepository 创建退款仓储实例。
func NewRefundRepository(db *gorm.DB) *RefundRepository { return &RefundRepository{db: db} }

// Create 持久化一条新的退款记录。
func (r *RefundRepository) Create(ctx context.Context, refund *domain.Refund) error {
	return r.db.WithContext(ctx).Create(refund).Error
}

// SumByPaymentID 返回指定支付 ID 下所有未取消退款的总金额（单位：分）。
// 用于强制执行"退款总额 ≤ 原始支付金额"的业务规则。
func (r *RefundRepository) SumByPaymentID(ctx context.Context, paymentID int64) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&domain.Refund{}).
		Where("payment_id = ? AND status != ?", paymentID, "cancelled").
		Select("COALESCE(SUM(amount_cents), 0)").
		Scan(&total).Error
	return total, err
}
