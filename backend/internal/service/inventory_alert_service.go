package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

const (
	defaultInventoryAlertTemplate = "inventory_alert"
	defaultInventoryAlertWindow   = 30 * time.Minute
)

// inventoryAlertRepo 抽象库存预警依赖的数据访问能力。
type inventoryAlertRepo interface {
	ListAllInventories(ctx context.Context) ([]domain.CabinInventory, error)
	GetInventoryBySKU(ctx context.Context, skuID int64) (domain.CabinInventory, error)
	SetAlertThreshold(ctx context.Context, skuID int64, threshold int) error
}

type inventoryAlertNotifier interface {
	Enqueue(ctx context.Context, userID int64, channel, template, payload string) error
}

// InventoryAlertService 提供库存预警查询与阈值配置能力。
type InventoryAlertService struct {
	repo         inventoryAlertRepo
	notifier     inventoryAlertNotifier
	dedupeWindow time.Duration
	lastSent     map[string]time.Time
	now          func() time.Time
	mu           sync.Mutex
}

// NewInventoryAlertService 创建库存预警服务。
func NewInventoryAlertService(repo inventoryAlertRepo) *InventoryAlertService {
	return NewInventoryAlertServiceWithNotify(repo, nil, defaultInventoryAlertWindow)
}

func NewInventoryAlertServiceWithNotify(repo inventoryAlertRepo, notifier inventoryAlertNotifier, dedupeWindow time.Duration) *InventoryAlertService {
	if dedupeWindow <= 0 {
		dedupeWindow = defaultInventoryAlertWindow
	}
	return &InventoryAlertService{
		repo:         repo,
		notifier:     notifier,
		dedupeWindow: dedupeWindow,
		lastSent:     make(map[string]time.Time),
		now:          time.Now,
	}
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
		if inv.AlertThreshold > 0 && available <= inv.AlertThreshold {
			alerts = append(alerts, domain.InventoryAlert{
				CabinSKUID:     inv.CabinSKUID,
				Available:      available,
				AlertThreshold: inv.AlertThreshold,
			})
		}
	}
	return alerts, nil
}

// AvailableWithAlert 返回指定 SKU 的可用库存与是否触发预警。
// 预警规则：threshold > 0 且 available <= threshold。
func (s *InventoryAlertService) AvailableWithAlert(ctx context.Context, skuID int64) (int, bool, error) {
	inv, err := s.repo.GetInventoryBySKU(ctx, skuID)
	if err != nil {
		return 0, false, err
	}
	available := inv.Total - inv.Locked - inv.Sold
	alert := inv.AlertThreshold > 0 && available <= inv.AlertThreshold
	return available, alert, nil
}

// SetThreshold 设置指定 SKU 的库存预警阈值。
func (s *InventoryAlertService) SetThreshold(ctx context.Context, skuID int64, threshold int) error {
	return s.repo.SetAlertThreshold(ctx, skuID, threshold)
}

// ScanAndNotify scans current alerts and sends deduplicated notifications.
func (s *InventoryAlertService) ScanAndNotify(ctx context.Context, userID int64, channel string) (int, error) {
	alerts, err := s.CheckAlerts(ctx)
	if err != nil {
		return 0, err
	}
	if s.notifier == nil || len(alerts) == 0 {
		return 0, nil
	}

	now := s.now()
	notified := 0
	for _, alert := range alerts {
		key := fmt.Sprintf("sku:%d", alert.CabinSKUID)

		s.mu.Lock()
		last, seen := s.lastSent[key]
		if seen && now.Sub(last) < s.dedupeWindow {
			s.mu.Unlock()
			continue
		}
		s.mu.Unlock()

		payload := fmt.Sprintf(`{"cabin_sku_id":%d,"available":%d,"threshold":%d}`, alert.CabinSKUID, alert.Available, alert.AlertThreshold)
		if err := s.notifier.Enqueue(ctx, userID, channel, defaultInventoryAlertTemplate, payload); err != nil {
			return notified, err
		}

		s.mu.Lock()
		s.lastSent[key] = now
		s.mu.Unlock()
		notified++
	}

	return notified, nil
}
