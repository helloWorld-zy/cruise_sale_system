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

// ListByCruise 查询指定邮轮下的所有设施，按排序权重和 ID 降序排列。
func (r *FacilityRepository) ListByCruise(ctx context.Context, cruiseID int64) ([]domain.Facility, error) {
	var items []domain.Facility
	if err := r.db.WithContext(ctx).Model(&domain.Facility{}).Where("cruise_id = ?", cruiseID).Order("sort_order desc, id desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// Delete 软删除指定的设施记录。
func (r *FacilityRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Facility{}, id).Error
}
