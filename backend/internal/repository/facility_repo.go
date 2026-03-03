package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// FacilityRepository 提供设施实体的数据库操作。
type FacilityRepository struct {
	db *gorm.DB // 数据库连接实例
}

// NewFacilityRepository 创建设施仓储实例。
func NewFacilityRepository(db *gorm.DB) *FacilityRepository {
	return &FacilityRepository{db: db}
}

// Create 插入一条新的设施记录。
func (r *FacilityRepository) Create(ctx context.Context, facility *domain.Facility) error {
	return r.db.WithContext(ctx).Create(facility).Error
}

// Update 保存设施的所有字段修改。
func (r *FacilityRepository) Update(ctx context.Context, facility *domain.Facility) error {
	return r.db.WithContext(ctx).Save(facility).Error
}

// GetByID 根据主键查询设施记录。
func (r *FacilityRepository) GetByID(ctx context.Context, id int64) (*domain.Facility, error) {
	var item domain.Facility
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// ListByCruise 查询指定邮轮下的所有设施，按排序权重和 ID 降序排列。
func (r *FacilityRepository) ListByCruise(ctx context.Context, cruiseID int64) ([]domain.Facility, error) {
	var items []domain.Facility
	if err := r.db.WithContext(ctx).Model(&domain.Facility{}).Where("cruise_id = ?", cruiseID).Order("sort_order desc, id desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// ListByCruiseAndCategory 按邮轮和分类筛选设施。
func (r *FacilityRepository) ListByCruiseAndCategory(ctx context.Context, cruiseID, categoryID int64) ([]domain.Facility, error) {
	var items []domain.Facility
	q := r.db.WithContext(ctx).Model(&domain.Facility{}).Where("cruise_id = ?", cruiseID)
	if categoryID > 0 {
		q = q.Where("category_id = ?", categoryID)
	}
	if err := q.Order("sort_order desc, id desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// Delete 软删除指定的设施记录。
func (r *FacilityRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Facility{}, id).Error
}
