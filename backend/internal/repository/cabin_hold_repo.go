package repository

import (
	"fmt"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CabinHoldRepository 提供舱位占座相关的数据访问实现。
type CabinHoldRepository struct {
	db *gorm.DB
}

// NewCabinHoldRepository 创建舱位占座仓储实例。
func NewCabinHoldRepository(db *gorm.DB) *CabinHoldRepository {
	return &CabinHoldRepository{db: db}
}

// ExistsActiveHoldTx 判断指定用户在当前时刻是否存在有效占座。
func (r *CabinHoldRepository) ExistsActiveHoldTx(tx *gorm.DB, skuID, userID int64, now time.Time) (bool, error) {
	db := tx
	if db == nil {
		db = r.db
	}

	var count int64
	err := db.Model(&domain.CabinHold{}).
		Where("cabin_sku_id = ? AND user_id = ? AND expires_at > ?", skuID, userID, now).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateHoldTx 创建占座记录，并清理该用户该 SKU 的过期占座。
func (r *CabinHoldRepository) CreateHoldTx(tx *gorm.DB, hold *domain.CabinHold) error {
	db := tx
	if db == nil {
		db = r.db
	}
	if err := db.Where("cabin_sku_id = ? AND user_id = ? AND expires_at <= ?", hold.CabinSKUID, hold.UserID, time.Now()).Delete(&domain.CabinHold{}).Error; err != nil {
		return err
	}
	return db.Create(hold).Error
}

// AdjustInventoryTx 在事务中按增量调整库存并写入库存日志。
func (r *CabinHoldRepository) AdjustInventoryTx(tx *gorm.DB, skuID int64, delta int, reason string) error {
	db := tx
	if db == nil {
		db = r.db
	}

	var inv domain.CabinInventory
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("cabin_sku_id = ?", skuID).First(&inv).Error; err != nil {
		return err
	}

	if inv.Total+delta < 0 {
		return fmt.Errorf("cabin_sku_id=%d: %w", skuID, domain.ErrInsufficientInventory)
	}

	inv.Total += delta
	if err := db.Save(&inv).Error; err != nil {
		return err
	}

	return db.Create(&domain.InventoryLog{
		CabinSKUID: skuID,
		Change:     delta,
		Reason:     reason,
	}).Error
}
