package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create(ctx context.Context, company *domain.CruiseCompany) error {
	return r.db.WithContext(ctx).Create(company).Error
}

func (r *CompanyRepository) Update(ctx context.Context, company *domain.CruiseCompany) error {
	return r.db.WithContext(ctx).Save(company).Error
}

func (r *CompanyRepository) GetByID(ctx context.Context, id int64) (*domain.CruiseCompany, error) {
	var company domain.CruiseCompany
	if err := r.db.WithContext(ctx).First(&company, id).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *CompanyRepository) List(ctx context.Context, keyword string, page, pageSize int) ([]domain.CruiseCompany, int64, error) {
	var items []domain.CruiseCompany
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.CruiseCompany{})
	if keyword != "" {
		q = q.Where("name LIKE ?", "%"+keyword+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("sort_order desc, id desc").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *CompanyRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.CruiseCompany{}, id).Error
}
