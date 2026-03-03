package repository

import (
	"context"
	"errors"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

var (
	// ErrInvalidOrderStatusTransition 表示请求的订单状态流转不符合状态机约束。
	ErrInvalidOrderStatusTransition = errors.New("invalid order status transition")
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
	return r.TransitionStatus(ctx, id, status, 0, "")
}

// TransitionStatus 通过统一入口变更订单状态，并在同一事务写入状态日志。
func (r *BookingRepository) TransitionStatus(ctx context.Context, id int64, status string, operatorID int64, remark string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var current domain.Booking
		if err := tx.First(&current, id).Error; err != nil {
			return err
		}
		if !current.CanTransitionTo(status) {
			return ErrInvalidOrderStatusTransition
		}
		if err := tx.Model(&domain.Booking{}).Where("id = ?", id).Update("status", status).Error; err != nil {
			return err
		}
		if remark == "" {
			remark = "status transition"
		}
		log := &domain.OrderStatusLog{
			OrderID:    id,
			FromStatus: current.Status,
			ToStatus:   status,
			OperatorID: operatorID,
			Remark:     remark,
		}
		if err := tx.Create(log).Error; err != nil {
			return err
		}
		return nil
	})
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

type BookingFilter struct {
	Status    string
	Phone     string
	RouteID   int64
	VoyageID  int64
	StartDate *string
	EndDate   *string
	BookingNo string
}

func (r *BookingRepository) ListWithFilter(ctx context.Context, filter BookingFilter, page, pageSize int) ([]domain.Booking, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	query := r.db.WithContext(ctx).Model(&domain.Booking{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.VoyageID > 0 {
		query = query.Where("voyage_id = ?", filter.VoyageID)
	}
	if filter.BookingNo != "" {
		query = query.Where("id LIKE ?", "%"+filter.BookingNo+"%")
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []domain.Booking
	err := query.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error
	return items, total, err
}
