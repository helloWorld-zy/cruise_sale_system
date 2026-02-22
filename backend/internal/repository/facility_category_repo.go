package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

type FacilityCategoryRepository struct {
	db *gorm.DB
}

func NewFacilityCategoryRepository(db *gorm.DB) *FacilityCategoryRepository {
	return &FacilityCategoryRepository{db: db}
}

func (r *FacilityCategoryRepository) Create(ctx context.Context, category *domain.FacilityCategory) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *FacilityCategoryRepository) List(ctx context.Context) ([]domain.FacilityCategory, error) {
	var items []domain.FacilityCategory
	if err := r.db.WithContext(ctx).Model(&domain.FacilityCategory{}).Order("sort_order desc, id desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *FacilityCategoryRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.FacilityCategory{}, id).Error
}
