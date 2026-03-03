package service

import (
	"context"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

// cabinAdminRepo 定义舱位管理相关的数据访问能力。
type cabinAdminRepo interface {
	CreateSKU(ctx context.Context, s *domain.CabinSKU) error
	UpdateSKU(ctx context.Context, s *domain.CabinSKU) error
	GetSKUByID(ctx context.Context, id int64) (*domain.CabinSKU, error)
	DeleteSKU(ctx context.Context, id int64) error
	ListSKUByVoyage(ctx context.Context, voyageID int64) ([]domain.CabinSKU, error)
	ListSKUFiltered(ctx context.Context, f domain.CabinSKUFilter) ([]domain.CabinSKU, int64, error)
	BatchUpdateStatus(ctx context.Context, ids []int64, status int16) error
	GetInventoryBySKU(ctx context.Context, skuID int64) (domain.CabinInventory, error)
	ListAllInventories(ctx context.Context) ([]domain.CabinInventory, error)
	SetAlertThreshold(ctx context.Context, skuID int64, threshold int) error
	AdjustInventoryAtomic(ctx context.Context, skuID int64, delta int) error
	AppendInventoryLog(ctx context.Context, log *domain.InventoryLog) error
	ListPricesBySKU(ctx context.Context, skuID int64) ([]domain.CabinPrice, error)
	UpsertPrice(ctx context.Context, p *domain.CabinPrice) error
	BatchSetPrice(ctx context.Context, skuID int64, start, end time.Time, occupancy int, priceCents, childPriceCents, singleSupplementCents int64, priceType string) error
	GetCategoryTree(ctx context.Context) (interface{}, error)
}

// CabinAdminService 提供后台舱位、库存与价格管理能力。
type CabinAdminService struct {
	repo cabinAdminRepo
}

// NewCabinAdminService 创建舱位后台管理服务。
func NewCabinAdminService(repo cabinAdminRepo) *CabinAdminService {
	return &CabinAdminService{repo: repo}
}

// ListByVoyage 按航次查询舱位 SKU 列表。
func (s *CabinAdminService) ListByVoyage(ctx context.Context, voyageID int64) ([]domain.CabinSKU, error) {
	return s.repo.ListSKUByVoyage(ctx, voyageID)
}

// FilteredList 按条件分页查询舱位 SKU。
func (s *CabinAdminService) FilteredList(ctx context.Context, f domain.CabinSKUFilter) ([]domain.CabinSKU, int64, error) {
	return s.repo.ListSKUFiltered(ctx, f)
}

// BatchUpdateStatus 批量更新舱位 SKU 状态。
func (s *CabinAdminService) BatchUpdateStatus(ctx context.Context, ids []int64, status int16) error {
	return s.repo.BatchUpdateStatus(ctx, ids, status)
}

// Create 创建舱位 SKU。
func (s *CabinAdminService) Create(ctx context.Context, sku *domain.CabinSKU) error {
	return s.repo.CreateSKU(ctx, sku)
}

// Update 更新舱位 SKU。
func (s *CabinAdminService) Update(ctx context.Context, sku *domain.CabinSKU) error {
	return s.repo.UpdateSKU(ctx, sku)
}

// GetByID 查询单个舱位 SKU 详情。
func (s *CabinAdminService) GetByID(ctx context.Context, id int64) (*domain.CabinSKU, error) {
	return s.repo.GetSKUByID(ctx, id)
}

// Delete 删除舱位 SKU。
func (s *CabinAdminService) Delete(ctx context.Context, id int64) error {
	return s.repo.DeleteSKU(ctx, id)
}

// GetInventory 查询指定 SKU 的当前库存。
func (s *CabinAdminService) GetInventory(ctx context.Context, skuID int64) (domain.CabinInventory, error) {
	return s.repo.GetInventoryBySKU(ctx, skuID)
}

// GetAlerts 查询低库存预警列表。
func (s *CabinAdminService) GetAlerts(ctx context.Context) ([]domain.InventoryAlert, error) {
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

// SetAlertThreshold 设置 SKU 的库存预警阈值。
func (s *CabinAdminService) SetAlertThreshold(ctx context.Context, skuID int64, threshold int) error {
	return s.repo.SetAlertThreshold(ctx, skuID, threshold)
}

// AdjustInventory 原子调整库存并追加库存变更日志。
func (s *CabinAdminService) AdjustInventory(ctx context.Context, skuID int64, delta int, reason string) error {
	if err := s.repo.AdjustInventoryAtomic(ctx, skuID, delta); err != nil {
		return err
	}
	return s.repo.AppendInventoryLog(ctx, &domain.InventoryLog{CabinSKUID: skuID, Change: delta, Reason: reason})
}

// ListPrices 查询指定 SKU 的价格列表。
func (s *CabinAdminService) ListPrices(ctx context.Context, skuID int64) ([]domain.CabinPrice, error) {
	return s.repo.ListPricesBySKU(ctx, skuID)
}

// UpsertPrice 新增或更新指定 SKU 的价格记录。
func (s *CabinAdminService) UpsertPrice(ctx context.Context, p *domain.CabinPrice) error {
	return s.repo.UpsertPrice(ctx, p)
}

// BatchSetPrice 按日期范围批量设置价格。
func (s *CabinAdminService) BatchSetPrice(ctx context.Context, skuID int64, start, end time.Time, occupancy int, priceCents, childPriceCents, singleSupplementCents int64, priceType string) error {
	return s.repo.BatchSetPrice(ctx, skuID, start, end, occupancy, priceCents, childPriceCents, singleSupplementCents, priceType)
}

// GetCategoryTree 获取邮轮→航线→舱型三级分类树。
func (s *CabinAdminService) GetCategoryTree(ctx context.Context) (interface{}, error) {
	return s.repo.GetCategoryTree(ctx)
}
