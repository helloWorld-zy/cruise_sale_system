package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

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

// ListSKUFiltered 按条件分页查询舱房 SKU。
func (r *CabinRepository) ListSKUFiltered(ctx context.Context, f domain.CabinSKUFilter) ([]domain.CabinSKU, int64, error) {
	var out []domain.CabinSKU
	var total int64

	page := f.Page
	if page <= 0 {
		page = 1
	}
	pageSize := f.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	q := r.db.WithContext(ctx).Model(&domain.CabinSKU{})
	if f.VoyageID > 0 {
		q = q.Where("voyage_id = ?", f.VoyageID)
	}
	if f.CabinTypeID > 0 {
		q = q.Where("cabin_type_id = ?", f.CabinTypeID)
	}
	if f.Status != nil {
		q = q.Where("status = ?", *f.Status)
	}
	if strings.TrimSpace(f.Keyword) != "" {
		q = q.Where("code LIKE ?", "%"+strings.TrimSpace(f.Keyword)+"%")
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&out).Error; err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

// BatchUpdateStatus 批量更新舱房 SKU 的上下架状态。
func (r *CabinRepository) BatchUpdateStatus(ctx context.Context, ids []int64, status int16) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&domain.CabinSKU{}).Where("id IN ?", ids).Update("status", status)
		if res.Error != nil {
			return res.Error
		}
		if int(res.RowsAffected) != len(ids) {
			return fmt.Errorf("batch update cabin status affected=%d expected=%d", res.RowsAffected, len(ids))
		}
		return nil
	})
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

// ListAllInventories 查询全部库存记录。
func (r *CabinRepository) ListAllInventories(ctx context.Context) ([]domain.CabinInventory, error) {
	var out []domain.CabinInventory
	return out, r.db.WithContext(ctx).Order("cabin_sku_id asc").Find(&out).Error
}

// SetAlertThreshold 设置指定 SKU 的库存预警阈值。
func (r *CabinRepository) SetAlertThreshold(ctx context.Context, skuID int64, threshold int) error {
	return r.db.WithContext(ctx).Model(&domain.CabinInventory{}).Where("cabin_sku_id = ?", skuID).Update("alert_threshold", threshold).Error
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

// Create 实现 PriceRepo 接口的 Create 方法。
func (r *CabinRepository) Create(ctx context.Context, p *domain.CabinPrice) error {
	return r.db.WithContext(ctx).Create(p).Error
}

// BatchSetPrice 按日期区间批量设置价格。
func (r *CabinRepository) BatchSetPrice(ctx context.Context, skuID int64, start, end time.Time, occupancy int, priceCents, childPriceCents, singleSupplementCents int64, priceType string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
			p := &domain.CabinPrice{
				CabinSKUID:            skuID,
				Date:                  d,
				Occupancy:             occupancy,
				PriceCents:            priceCents,
				ChildPriceCents:       childPriceCents,
				SingleSupplementCents: singleSupplementCents,
				PriceType:             priceType,
			}
			if err := tx.Save(p).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetCategoryTree 获取邮轮→航线→舱型三级分类树。
func (r *CabinRepository) GetCategoryTree(ctx context.Context) (interface{}, error) {
	var cruises []domain.Cruise
	if err := r.db.WithContext(ctx).Where("status = ?", 1).Find(&cruises).Error; err != nil {
		return nil, err
	}
	type routeWithVoyages struct {
		domain.Route
		Voyages []domain.Voyage `json:"voyages"`
	}
	type cruiseWithRoutes struct {
		domain.Cruise
		Routes []routeWithVoyages `json:"routes"`
	}
	result := make([]cruiseWithRoutes, 0, len(cruises))
	for _, c := range cruises {
		var routes []domain.Route
		if err := r.db.WithContext(ctx).Where("status = ?", 1).Find(&routes).Error; err != nil {
			continue
		}
		cruiseRoutes := make([]routeWithVoyages, 0, len(routes))
		for _, rt := range routes {
			var voyages []domain.Voyage
			if err := r.db.WithContext(ctx).Where("route_id = ? AND cruise_id = ?", rt.ID, c.ID).Find(&voyages).Error; err != nil {
				continue
			}
			cruiseRoutes = append(cruiseRoutes, routeWithVoyages{
				Route:   rt,
				Voyages: voyages,
			})
		}
		result = append(result, cruiseWithRoutes{
			Cruise: c,
			Routes: cruiseRoutes,
		})
	}
	return result, nil
}

// 编译时接口实现检查
var _ domain.CabinSKURepository = (*CabinRepository)(nil)
