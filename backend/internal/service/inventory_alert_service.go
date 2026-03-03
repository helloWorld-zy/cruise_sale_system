package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

// inventoryAlertRepo 抽象库存预警依赖的数据访问能力。
type inventoryAlertRepo interface {
	ListAllInventories(ctx context.Context) ([]domain.CabinInventory, error)
	SetAlertThreshold(ctx context.Context, skuID int64, threshold int) error
}

// InventoryAlertService 提供库存预警查询与阈值配置能力。
type InventoryAlertService struct {
	repo inventoryAlertRepo
}

// NewInventoryAlertService 创建库存预警服务。
func NewInventoryAlertService(repo inventoryAlertRepo) *InventoryAlertService {
	return &InventoryAlertService{repo: repo}
}

// CheckAlerts 返回所有低于阈值的库存项。
func (s *InventoryAlertService) CheckAlerts(ctx context.Context) ([]domain.InventoryAlert, error) {
	invs, err := s.repo.ListAllInventories(ctx)
	if err != nil {
		return nil, err
	}
	alerts := make([]domain.InventoryAlert, 0)
	for _, inv := range invs {
		available := inv.Total - inv.Locked - inv.Sold
		if inv.AlertThreshold > 0 && available < inv.AlertThreshold {
			alerts = append(alerts, domain.InventoryAlert{
				CabinSKUID:     inv.CabinSKUID,
				Available:      available,
				AlertThreshold: inv.AlertThreshold,
			})
		}
	}
	return alerts, nil
}

// SetThreshold 设置指定 SKU 的库存预警阈值。
func (s *InventoryAlertService) SetThreshold(ctx context.Context, skuID int64, threshold int) error {
	return s.repo.SetAlertThreshold(ctx, skuID, threshold)
}
