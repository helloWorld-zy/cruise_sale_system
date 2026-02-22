package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

type FacilityRepository struct {
	db *gorm.DB
}

func NewFacilityRepository(db *gorm.DB) *FacilityRepository {
	return &FacilityRepository{db: db}
}

func (r *FacilityRepository) Create(ctx context.Context, facility *domain.Facility) error {
	return r.db.WithContext(ctx).Create(facility).Error
}

func (r *FacilityRepository) ListByCruise(ctx context.Context, cruiseID int64) ([]domain.Facility, error) {
	var items []domain.Facility
	if err := r.db.WithContext(ctx).Model(&domain.Facility{}).Where("cruise_id = ?", cruiseID).Order("sort_order desc, id desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *FacilityRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Facility{}, id).Error
}
