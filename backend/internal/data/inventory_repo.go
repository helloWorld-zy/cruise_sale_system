package data

import (
	"context"
	"cruise_booking_system/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryRepository struct{}

func NewInventoryRepository() *InventoryRepository {
	return &InventoryRepository{}
}

func (r *InventoryRepository) GetByVoyageAndType(ctx context.Context, voyageID, cabinTypeID uuid.UUID) (*model.Inventory, error) {
	var inv model.Inventory
	if err := DB.WithContext(ctx).Where("voyage_id = ? AND cabin_type_id = ?", voyageID, cabinTypeID).First(&inv).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

// DecreaseAvailableQty decreases available and increases reserved in a transaction
func (r *InventoryRepository) Reserve(ctx context.Context, tx *gorm.DB, voyageID, cabinTypeID uuid.UUID, qty int) error {
	result := tx.WithContext(ctx).Model(&model.Inventory{}).
		Where("voyage_id = ? AND cabin_type_id = ? AND available_qty >= ?", voyageID, cabinTypeID, qty).
		Updates(map[string]interface{}{
			"available_qty": gorm.Expr("available_qty - ?", qty),
			"reserved_qty":  gorm.Expr("reserved_qty + ?", qty),
		})
	
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Or "insufficient inventory"
	}
	return nil
}

func (r *InventoryRepository) Release(ctx context.Context, tx *gorm.DB, voyageID, cabinTypeID uuid.UUID, qty int) error {
	result := tx.WithContext(ctx).Model(&model.Inventory{}).
		Where("voyage_id = ? AND cabin_type_id = ?", voyageID, cabinTypeID).
		Updates(map[string]interface{}{
			"available_qty": gorm.Expr("available_qty + ?", qty),
			"reserved_qty":  gorm.Expr("reserved_qty - ?", qty),
		})
	
	return result.Error
}

func (r *InventoryRepository) Update(ctx context.Context, inv *model.Inventory) error {
	// Upsert or Update
	return DB.WithContext(ctx).Save(inv).Error
}
