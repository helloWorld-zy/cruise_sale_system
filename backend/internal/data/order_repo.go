package data

import (
	"context"
	"cruise_booking_system/internal/model"

	"gorm.io/gorm"
)

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (r *OrderRepository) Create(ctx context.Context, tx *gorm.DB, order *model.Order) error {
	return tx.WithContext(ctx).Create(order).Error
}
