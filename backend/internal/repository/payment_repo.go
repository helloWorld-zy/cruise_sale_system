package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// PaymentRepository 基于 PostgreSQL 提供支付记录的持久化操作。
type PaymentRepository struct{ db *gorm.DB }

// NewPaymentRepository 创建支付仓储实例。
func NewPaymentRepository(db *gorm.DB) *PaymentRepository { return &PaymentRepository{db: db} }

// Create 持久化一条新的支付记录。
func (r *PaymentRepository) Create(ctx context.Context, p *domain.Payment) error {
	return r.db.WithContext(ctx).Create(p).Error
}

// FindByTradeNo 根据支付提供商的交易流水号查找支付记录。
// 当记录不存在时返回 gorm.ErrRecordNotFound 错误。
func (r *PaymentRepository) FindByTradeNo(ctx context.Context, tradeNo string) (*domain.Payment, error) {
	var p domain.Payment
	err := r.db.WithContext(ctx).Where("trade_no = ?", tradeNo).First(&p).Error
	return &p, err
}

// FindByID 根据主键查找支付记录。
func (r *PaymentRepository) FindByID(ctx context.Context, id int64) (*domain.Payment, error) {
	var p domain.Payment
	err := r.db.WithContext(ctx).First(&p, id).Error
	return &p, err
}

// UpdateStatus 更新指定支付记录的状态字段。
func (r *PaymentRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).
		Model(&domain.Payment{}).
		Where("id = ?", id).
		Update("status", status).Error
}
