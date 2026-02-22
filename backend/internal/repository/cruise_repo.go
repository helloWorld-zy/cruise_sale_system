package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

type CruiseRepository struct {
	db *gorm.DB
}

func NewCruiseRepository(db *gorm.DB) *CruiseRepository {
	return &CruiseRepository{db: db}
}

func (r *CruiseRepository) Create(ctx context.Context, cruise *domain.Cruise) error {
	return r.db.WithContext(ctx).Create(cruise).Error
}

func (r *CruiseRepository) Update(ctx context.Context, cruise *domain.Cruise) error {
	return r.db.WithContext(ctx).Save(cruise).Error
}

func (r *CruiseRepository) GetByID(ctx context.Context, id int64) (*domain.Cruise, error) {
	var cruise domain.Cruise
	if err := r.db.WithContext(ctx).First(&cruise, id).Error; err != nil {
		return nil, err
	}
	return &cruise, nil
}

func (r *CruiseRepository) List(ctx context.Context, companyID int64, page, pageSize int) ([]domain.Cruise, int64, error) {
	var items []domain.Cruise
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.Cruise{})
	if companyID > 0 {
		q = q.Where("company_id = ?", companyID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("sort_order desc, id desc").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *CruiseRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Cruise{}, id).Error
}
