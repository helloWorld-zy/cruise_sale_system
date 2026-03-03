package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

// OrderTimeoutRepo 订单超时查询接口。
type OrderTimeoutRepo interface {
	FindExpiredOrders(ctx context.Context, timeout time.Duration) ([]domain.Booking, error)
	TransitionStatus(ctx context.Context, id int64, status string, operatorID int64, remark string) error
}

// InventoryReleaser 库存释放接口。
type InventoryReleaser interface {
	ReleaseLocked(ctx context.Context, skuID int64, quantity int) error
}

// OrderTimeoutService 处理订单超时自动关闭。
type OrderTimeoutService struct {
	orderRepo     OrderTimeoutRepo
	inventoryRepo InventoryReleaser
	mu            sync.Mutex
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
	s.mu.Lock()
	defer s.mu.Unlock()

	orders, err := s.orderRepo.FindExpiredOrders(ctx, timeout)
	if err != nil {
		return 0, err
	}

	closed := 0
	for _, order := range orders {
		if order.Status != domain.OrderStatusPendingPayment {
			continue
		}
		if err := s.orderRepo.TransitionStatus(ctx, order.ID, domain.OrderStatusCancelled, 0, "timeout auto close"); err != nil {
			continue
		}
		if err := s.inventoryRepo.ReleaseLocked(ctx, order.CabinSKUID, 1); err != nil {
			if rbErr := s.orderRepo.TransitionStatus(ctx, order.ID, order.Status, 0, "rollback: inventory release failed"); rbErr != nil {
				log.Printf("order_timeout: rollback order %d failed: %v (original: %v)", order.ID, rbErr, err)
			}
			continue
		}
		closed++
	}
	return closed, nil
}
