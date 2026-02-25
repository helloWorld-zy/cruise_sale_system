package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

// InventoryRepo 定义库存操作的端口接口。
// AdjustAtomic 必须使用单条 SQL 原子化更新，防止并发超卖
// （CRITICAL-01 修复项）。
type InventoryRepo interface {
	GetBySKU(skuID int64) (domain.CabinInventory, error) // 根据 SKU ID 查询库存
	// AdjustAtomic 原子化调整库存总量。
	// 必须使用单条 SQL：UPDATE ... SET total = total + delta WHERE ... AND total + delta >= 0
	// 当约束不满足时返回 domain.ErrInsufficientInventory。
	AdjustAtomic(ctx context.Context, id int64, delta int) error
	AppendLog(ctx context.Context, log *domain.InventoryLog) error // 追加库存变动审计日志
}

// InventoryService 提供舱房库存管理的业务逻辑。
type InventoryService struct{ repo InventoryRepo }

// NewInventoryService 创建库存服务实例。
func NewInventoryService(repo InventoryRepo) *InventoryService { return &InventoryService{repo: repo} }

// Adjust 对指定舱房 SKU 的库存进行原子化调整，并记录审计日志。
// delta 为正数表示增加库存，负数表示减少库存。
// 当库存不足时返回 domain.ErrInsufficientInventory。
func (s *InventoryService) Adjust(skuID int64, delta int, reason string) error {
	ctx := context.Background()
	if err := s.repo.AdjustAtomic(ctx, skuID, delta); err != nil {
		return err
	}
	return s.repo.AppendLog(ctx, &domain.InventoryLog{
		CabinSKUID: skuID,
		Change:     delta,
		Reason:     reason,
	})
}

// Available 计算并返回指定 SKU 的可用库存数量。
// 可用库存 = 总量 - 锁定量 - 已售量。
func (s *InventoryService) Available(skuID int64) (int, error) {
	inv, err := s.repo.GetBySKU(skuID)
	if err != nil {
		return 0, err
	}
	return inv.Total - inv.Locked - inv.Sold, nil
}
