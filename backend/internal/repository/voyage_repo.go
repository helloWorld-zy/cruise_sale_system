package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// VoyageRepository 实现 domain.VoyageRepository 接口，提供航次实体的数据库操作。
type VoyageRepository struct{ db *gorm.DB }

// NewVoyageRepository 创建航次仓储实例。
func NewVoyageRepository(db *gorm.DB) *VoyageRepository { return &VoyageRepository{db: db} }

// Create 插入一条新的航次记录。
func (r *VoyageRepository) Create(ctx context.Context, v *domain.Voyage) error {
	return r.db.WithContext(ctx).Create(v).Error
}

// Update 保存航次的所有字段修改。
func (r *VoyageRepository) Update(ctx context.Context, v *domain.Voyage) error {
	return r.db.WithContext(ctx).Save(v).Error
}

// GetByID 根据主键查询航次记录。
func (r *VoyageRepository) GetByID(ctx context.Context, id int64) (*domain.Voyage, error) {
	var out domain.Voyage
	if err := r.db.WithContext(ctx).First(&out, id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

// ListByRoute 查询指定航线下的所有航次，按出发日期升序排列。
func (r *VoyageRepository) ListByRoute(ctx context.Context, routeID int64) ([]domain.Voyage, error) {
	var out []domain.Voyage
	return out, r.db.WithContext(ctx).Where("route_id = ?", routeID).Order("depart_date asc").Find(&out).Error
}

// Delete 删除指定的航次记录。
func (r *VoyageRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Voyage{}, id).Error
}

// 编译时接口实现检查
var _ domain.VoyageRepository = (*VoyageRepository)(nil)
