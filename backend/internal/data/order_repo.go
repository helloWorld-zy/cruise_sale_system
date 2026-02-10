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

func (r *OrderRepository) ListByUserID(ctx context.Context, userID string) ([]model.Order, error) {
	var orders []model.Order
	if err := DB.WithContext(ctx).Where("user_id = ?", userID).Preload("Items").Preload("Passengers").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id string) (*model.Order, error) {
	var order model.Order
	if err := DB.WithContext(ctx).Where("id = ?", id).Preload("Items").Preload("Passengers").First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, tx *gorm.DB, id string, status model.OrderStatus) error {
	db := DB
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *OrderRepository) List(ctx context.Context, status model.OrderStatus) ([]model.Order, error) {
	var orders []model.Order
	query := DB.WithContext(ctx).Preload("Items").Preload("Passengers")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *model.Order) error {
	return DB.WithContext(ctx).Save(order).Error
}

func (r *OrderRepository) GetTotalRevenue(ctx context.Context) (float64, error) {
	var total float64
	// Sum total_amount where status = paid/confirmed/completed
	err := DB.WithContext(ctx).Model(&model.Order{}).
		Where("status IN ?", []model.OrderStatus{model.OrderStatusPaid, model.OrderStatusConfirmed, model.OrderStatusCompleted}).
		Select("COALESCE(SUM(total_amount), 0)").Scan(&total).Error
	return total, err
}
