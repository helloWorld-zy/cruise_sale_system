package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// CompanyRepository 提供邮轮公司实体的数据库操作。
type CompanyRepository struct {
	db *gorm.DB // 数据库连接实例
}

// NewCompanyRepository 创建公司仓储实例。
func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

// Create 插入一条新的公司记录。
func (r *CompanyRepository) Create(ctx context.Context, company *domain.CruiseCompany) error {
	return r.db.WithContext(ctx).Create(company).Error
}

// Update 保存公司的所有字段修改。
func (r *CompanyRepository) Update(ctx context.Context, company *domain.CruiseCompany) error {
	return r.db.WithContext(ctx).Save(company).Error
}

// GetByID 根据主键查询公司记录。
func (r *CompanyRepository) GetByID(ctx context.Context, id int64) (*domain.CruiseCompany, error) {
	var company domain.CruiseCompany
	if err := r.db.WithContext(ctx).First(&company, id).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

// List 分页查询公司列表，支持按名称关键词模糊搜索。
// 返回值：公司列表、总记录数、错误信息。
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

// Delete 软删除指定的公司记录。
func (r *CompanyRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.CruiseCompany{}, id).Error
}
