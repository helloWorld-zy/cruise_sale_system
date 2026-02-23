package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// CruiseRepository 提供邮轮实体的数据库操作。
type CruiseRepository struct {
	db *gorm.DB // 数据库连接实例
}

// NewCruiseRepository 创建邮轮仓储实例。
func NewCruiseRepository(db *gorm.DB) *CruiseRepository {
	return &CruiseRepository{db: db}
}

// Create 插入一条新的邮轮记录。
func (r *CruiseRepository) Create(ctx context.Context, cruise *domain.Cruise) error {
	return r.db.WithContext(ctx).Create(cruise).Error
}

// Update 保存邮轮的所有字段修改。
func (r *CruiseRepository) Update(ctx context.Context, cruise *domain.Cruise) error {
	return r.db.WithContext(ctx).Save(cruise).Error
}

// GetByID 根据主键查询邮轮记录。
func (r *CruiseRepository) GetByID(ctx context.Context, id int64) (*domain.Cruise, error) {
	var cruise domain.Cruise
	if err := r.db.WithContext(ctx).First(&cruise, id).Error; err != nil {
		return nil, err
	}
	return &cruise, nil
}

// List 分页查询邮轮列表，可按公司 ID 过滤。
// 当 companyID > 0 时仅返回该公司下的邮轮。
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

// Delete 软删除指定的邮轮记录。
func (r *CruiseRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Cruise{}, id).Error
}
