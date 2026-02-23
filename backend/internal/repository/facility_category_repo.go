package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// FacilityCategoryRepository 提供设施分类实体的数据库操作。
type FacilityCategoryRepository struct {
	db *gorm.DB // 数据库连接实例
}

// NewFacilityCategoryRepository 创建设施分类仓储实例。
func NewFacilityCategoryRepository(db *gorm.DB) *FacilityCategoryRepository {
	return &FacilityCategoryRepository{db: db}
}

// Create 插入一条新的设施分类记录。
func (r *FacilityCategoryRepository) Create(ctx context.Context, category *domain.FacilityCategory) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// List 查询所有设施分类，按排序权重和 ID 降序排列。
func (r *FacilityCategoryRepository) List(ctx context.Context) ([]domain.FacilityCategory, error) {
	var items []domain.FacilityCategory
	if err := r.db.WithContext(ctx).Model(&domain.FacilityCategory{}).Order("sort_order desc, id desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// Delete 删除指定的设施分类记录。
func (r *FacilityCategoryRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.FacilityCategory{}, id).Error
}
