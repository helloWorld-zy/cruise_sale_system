package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

// cabinAdminRepo 定义舱位管理相关的数据访问能力。
type cabinAdminRepo interface {
	CreateSKU(ctx context.Context, s *domain.CabinSKU) error
	UpdateSKU(ctx context.Context, s *domain.CabinSKU) error
	GetSKUByID(ctx context.Context, id int64) (*domain.CabinSKU, error)
	DeleteSKU(ctx context.Context, id int64) error
	ListSKUByVoyage(ctx context.Context, voyageID int64) ([]domain.CabinSKU, error)
	GetInventoryBySKU(ctx context.Context, skuID int64) (domain.CabinInventory, error)
	AdjustInventoryAtomic(ctx context.Context, skuID int64, delta int) error
	AppendInventoryLog(ctx context.Context, log *domain.InventoryLog) error
	ListPricesBySKU(ctx context.Context, skuID int64) ([]domain.CabinPrice, error)
	UpsertPrice(ctx context.Context, p *domain.CabinPrice) error
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
