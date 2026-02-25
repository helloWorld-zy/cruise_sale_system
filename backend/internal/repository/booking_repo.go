package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// BookingRepository 提供预订实体的数据持久化能力。
type BookingRepository struct{ db *gorm.DB }

// NewBookingRepository 创建预订仓储实例。
func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

// Create 写入一条预订记录。
func (r *BookingRepository) Create(ctx context.Context, b *domain.Booking) error {
	err := r.db.WithContext(ctx).Create(b).Error
	return err
}

// InTx 在事务上下文内执行预订创建流程。
//
//go:noinline
func (r *BookingRepository) InTx(fn func(tx *gorm.DB, create func(b *domain.Booking) error) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		create := func(b *domain.Booking) error {
			return tx.Create(b).Error
		}
		return fn(tx, create)
	})
}
