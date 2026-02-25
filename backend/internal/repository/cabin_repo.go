package repository

import (
	"context"
	"fmt"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// CabinRepository 实现 domain.CabinSKURepository 接口，
// 提供舱房 SKU、库存和价格的数据库操作。
type CabinRepository struct{ db *gorm.DB }

// NewCabinRepository 创建舱房仓储实例。
func NewCabinRepository(db *gorm.DB) *CabinRepository { return &CabinRepository{db: db} }

// CreateSKU 创建舱房 SKU 记录。
func (r *CabinRepository) CreateSKU(ctx context.Context, v *domain.CabinSKU) error {
	return r.db.WithContext(ctx).Create(v).Error
}

// UpdateSKU 更新舱房 SKU 记录。
func (r *CabinRepository) UpdateSKU(ctx context.Context, v *domain.CabinSKU) error {
	return r.db.WithContext(ctx).Save(v).Error
}

// GetSKUByID 根据 ID 查询舱房 SKU。
func (r *CabinRepository) GetSKUByID(ctx context.Context, id int64) (*domain.CabinSKU, error) {
	var out domain.CabinSKU
	if err := r.db.WithContext(ctx).First(&out, id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

// ListSKUByVoyage 查询指定航次下的所有舱房 SKU。
func (r *CabinRepository) ListSKUByVoyage(ctx context.Context, voyageID int64) ([]domain.CabinSKU, error) {
	var out []domain.CabinSKU
	return out, r.db.WithContext(ctx).Where("voyage_id = ?", voyageID).Order("id desc").Find(&out).Error
}

// DeleteSKU 删除指定的舱房 SKU。
func (r *CabinRepository) DeleteSKU(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.CabinSKU{}, id).Error
}

// AdjustInventoryAtomic 使用单条原子化 SQL 更新库存总量，
// 防止并发请求导致竞态条件（CRITICAL-01 修复项）。
// 当 total+delta 为负时 WHERE 子句不匹配，会返回 ErrInsufficientInventory。
func (r *CabinRepository) AdjustInventoryAtomic(ctx context.Context, skuID int64, delta int) error {
	result := r.db.WithContext(ctx).Exec(
		`UPDATE cabin_inventories SET total = total + ?, updated_at = CURRENT_TIMESTAMP
		 WHERE cabin_sku_id = ? AND total + ? >= 0`,
		delta, skuID, delta,
	)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("cabin_sku_id=%d: %w", skuID, domain.ErrInsufficientInventory)
	}
	return nil
}

// GetInventoryBySKU 根据 SKU ID 查询库存信息。
func (r *CabinRepository) GetInventoryBySKU(ctx context.Context, skuID int64) (domain.CabinInventory, error) {
	var out domain.CabinInventory
	return out, r.db.WithContext(ctx).Where("cabin_sku_id = ?", skuID).First(&out).Error
}

// AppendInventoryLog 追加一条库存变动审计日志。
func (r *CabinRepository) AppendInventoryLog(ctx context.Context, log *domain.InventoryLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// ListPricesBySKU 查询指定 SKU 的价格列表，按日期和入住人数排序。
func (r *CabinRepository) ListPricesBySKU(ctx context.Context, skuID int64) ([]domain.CabinPrice, error) {
	var out []domain.CabinPrice
	return out, r.db.WithContext(ctx).Where("cabin_sku_id = ?", skuID).Order("date asc, occupancy asc").Find(&out).Error
}

// ListBySKU 兼容 PricingService 的方法。
func (r *CabinRepository) ListBySKU(ctx context.Context, skuID int64) ([]domain.CabinPrice, error) {
	return r.ListPricesBySKU(ctx, skuID)
}

// UpsertPrice 新增或更新价格记录。
//
//go:noinline
func (r *CabinRepository) UpsertPrice(ctx context.Context, p *domain.CabinPrice) error {
	err := r.db.WithContext(ctx).Save(p).Error
	return err
}

// 编译时接口实现检查
var _ domain.CabinSKURepository = (*CabinRepository)(nil)
