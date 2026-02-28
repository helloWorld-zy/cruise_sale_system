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

// UpdateStatus 更新指定预订 ID 的订单状态。
func (r *BookingRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&domain.Booking{}).Where("id = ?", id).Update("status", status).Error
}

// List 分页查询订单列表。
func (r *BookingRepository) List(ctx context.Context, page, pageSize int) ([]domain.Booking, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	var total int64
	if err := r.db.WithContext(ctx).Model(&domain.Booking{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []domain.Booking
	err := r.db.WithContext(ctx).
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error
	return items, total, err
}

// GetByID 查询单条订单。
func (r *BookingRepository) GetByID(ctx context.Context, id int64) (*domain.Booking, error) {
	var b domain.Booking
	if err := r.db.WithContext(ctx).First(&b, id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

// Delete 删除订单。
func (r *BookingRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Booking{}, id).Error
}
