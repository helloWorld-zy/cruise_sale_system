package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// CabinTypeRepository 提供舱房类型实体的数据库操作。
type CabinTypeRepository struct {
	db *gorm.DB // 数据库连接实例
}

// NewCabinTypeRepository 创建舱房类型仓储实例。
func NewCabinTypeRepository(db *gorm.DB) *CabinTypeRepository {
	return &CabinTypeRepository{db: db}
}

// Create 插入一条新的舱房类型记录。
func (r *CabinTypeRepository) Create(ctx context.Context, cabinType *domain.CabinType) error {
	return r.db.WithContext(ctx).Create(cabinType).Error
}

// Update 保存舱房类型的所有字段修改。
func (r *CabinTypeRepository) Update(ctx context.Context, cabinType *domain.CabinType) error {
	return r.db.WithContext(ctx).Save(cabinType).Error
}

// GetByID 根据主键查询舱房类型记录。
func (r *CabinTypeRepository) GetByID(ctx context.Context, id int64) (*domain.CabinType, error) {
	var cabinType domain.CabinType
	if err := r.db.WithContext(ctx).First(&cabinType, id).Error; err != nil {
		return nil, err
	}
	return &cabinType, nil
}

// ListByCruise 分页查询指定邮轮下的舱房类型列表。
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

// Delete 软删除指定的舱房类型记录。
func (r *CabinTypeRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.CabinType{}, id).Error
}
