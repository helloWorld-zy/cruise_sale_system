package service

import (
	"context"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

// OrderTimeoutRepo 订单超时查询接口。
type OrderTimeoutRepo interface {
	FindExpiredOrders(ctx context.Context, timeout time.Duration) ([]domain.Booking, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}

// InventoryReleaser 库存释放接口。
type InventoryReleaser interface {
	ReleaseLocked(ctx context.Context, skuID int64, quantity int) error
}

// OrderTimeoutService 处理订单超时自动关闭。
type OrderTimeoutService struct {
	orderRepo     OrderTimeoutRepo
	inventoryRepo InventoryReleaser
}

// NewOrderTimeoutService 创建订单超时服务。
func NewOrderTimeoutService(orderRepo OrderTimeoutRepo, inventoryRepo InventoryReleaser) *OrderTimeoutService {
	return &OrderTimeoutService{
		orderRepo:     orderRepo,
		inventoryRepo: inventoryRepo,
	}
}

// CloseExpiredOrders 关闭超时未支付的订单并释放库存。
func (s *OrderTimeoutService) CloseExpiredOrders(ctx context.Context, timeout time.Duration) (int, error) {
	orders, err := s.orderRepo.FindExpiredOrders(ctx, timeout)
	if err != nil {
		return 0, err
	}

	closed := 0
	for _, order := range orders {
		if err := s.orderRepo.UpdateStatus(ctx, order.ID, domain.OrderStatusCancelled); err != nil {
			continue
		}
		if err := s.inventoryRepo.ReleaseLocked(ctx, order.CabinSKUID, 1); err != nil {
			continue
		}
		closed++
	}
	return closed, nil
}
