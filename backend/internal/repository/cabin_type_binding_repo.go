package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// CabinTypeBindingRepository 提供舱型与邮轮绑定关系的数据访问实现。
type CabinTypeBindingRepository struct {
	db *gorm.DB
}

func NewCabinTypeBindingRepository(db *gorm.DB) *CabinTypeBindingRepository {
	return &CabinTypeBindingRepository{db: db}
}

func (r *CabinTypeBindingRepository) ReplaceCruiseBindings(ctx context.Context, cabinTypeID int64, cruiseIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("cabin_type_id = ?", cabinTypeID).Delete(&domain.CabinTypeCruiseBinding{}).Error; err != nil {
			return err
		}
		if len(cruiseIDs) == 0 {
			return nil
		}
		rows := make([]domain.CabinTypeCruiseBinding, 0, len(cruiseIDs))
		for _, cruiseID := range cruiseIDs {
			rows = append(rows, domain.CabinTypeCruiseBinding{CabinTypeID: cabinTypeID, CruiseID: cruiseID})
		}
		return tx.Create(&rows).Error
	})
}

func (r *CabinTypeBindingRepository) ListCruiseIDsByCabinType(ctx context.Context, cabinTypeID int64) ([]int64, error) {
	var ids []int64
	if err := r.db.WithContext(ctx).
		Model(&domain.CabinTypeCruiseBinding{}).
		Where("cabin_type_id = ?", cabinTypeID).
		Order("cruise_id asc").
		Pluck("cruise_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *CabinTypeBindingRepository) ListCabinTypeIDsByCruise(ctx context.Context, cruiseID int64) ([]int64, error) {
	var ids []int64
	if err := r.db.WithContext(ctx).
		Model(&domain.CabinTypeCruiseBinding{}).
		Where("cruise_id = ?", cruiseID).
		Order("cabin_type_id asc").
		Pluck("cabin_type_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *CabinTypeBindingRepository) HasCabinTypesByCruise(ctx context.Context, cruiseID int64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.CabinTypeCruiseBinding{}).Where("cruise_id = ?", cruiseID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
