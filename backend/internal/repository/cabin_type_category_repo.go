package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// CabinTypeCategoryRepository 提供舱型大类字典的数据访问实现。
type CabinTypeCategoryRepository struct {
	db *gorm.DB
}

func NewCabinTypeCategoryRepository(db *gorm.DB) *CabinTypeCategoryRepository {
	return &CabinTypeCategoryRepository{db: db}
}

func (r *CabinTypeCategoryRepository) Create(ctx context.Context, category *domain.CabinTypeCategory) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *CabinTypeCategoryRepository) Update(ctx context.Context, category *domain.CabinTypeCategory) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *CabinTypeCategoryRepository) GetByID(ctx context.Context, id int64) (*domain.CabinTypeCategory, error) {
	var item domain.CabinTypeCategory
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CabinTypeCategoryRepository) List(ctx context.Context) ([]domain.CabinTypeCategory, error) {
	var items []domain.CabinTypeCategory
	if err := r.db.WithContext(ctx).Model(&domain.CabinTypeCategory{}).Order("sort_order desc, id desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CabinTypeCategoryRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.CabinTypeCategory{}, id).Error
}
