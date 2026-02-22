package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

type CabinTypeRepository struct {
	db *gorm.DB
}

func NewCabinTypeRepository(db *gorm.DB) *CabinTypeRepository {
	return &CabinTypeRepository{db: db}
}

func (r *CabinTypeRepository) Create(ctx context.Context, cabinType *domain.CabinType) error {
	return r.db.WithContext(ctx).Create(cabinType).Error
}

func (r *CabinTypeRepository) Update(ctx context.Context, cabinType *domain.CabinType) error {
	return r.db.WithContext(ctx).Save(cabinType).Error
}

func (r *CabinTypeRepository) GetByID(ctx context.Context, id int64) (*domain.CabinType, error) {
	var cabinType domain.CabinType
	if err := r.db.WithContext(ctx).First(&cabinType, id).Error; err != nil {
		return nil, err
	}
	return &cabinType, nil
}

func (r *CabinTypeRepository) ListByCruise(ctx context.Context, cruiseID int64, page, pageSize int) ([]domain.CabinType, int64, error) {
	var items []domain.CabinType
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.CabinType{}).Where("cruise_id = ?", cruiseID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("sort_order desc, id desc").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *CabinTypeRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.CabinType{}, id).Error
}
