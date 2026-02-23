package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// routeRepositoryImpl 实现 domain.RouteRepository 接口。
// HIGH-04 修复：与 Company/Cruise 仓储保持一致的 DDD 端口模式。
type routeRepositoryImpl struct{ db *gorm.DB }

// NewRouteRepository 创建航线仓储实例。
func NewRouteRepository(db *gorm.DB) *routeRepositoryImpl { return &routeRepositoryImpl{db: db} }

// Create 插入一条新的航线记录。
func (r *routeRepositoryImpl) Create(ctx context.Context, v *domain.Route) error {
	return r.db.WithContext(ctx).Create(v).Error
}

// Update 保存航线的所有字段修改。
func (r *routeRepositoryImpl) Update(ctx context.Context, v *domain.Route) error {
	return r.db.WithContext(ctx).Save(v).Error
}

// GetByID 根据主键查询航线记录。
func (r *routeRepositoryImpl) GetByID(ctx context.Context, id int64) (*domain.Route, error) {
	var out domain.Route
	if err := r.db.WithContext(ctx).First(&out, id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

// List 查询所有航线，按 ID 降序排列。
func (r *routeRepositoryImpl) List(ctx context.Context) ([]domain.Route, error) {
	var out []domain.Route
	return out, r.db.WithContext(ctx).Order("id desc").Find(&out).Error
}

// Delete 删除指定的航线记录。
func (r *routeRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Route{}, id).Error
}

// 编译时接口实现检查
var _ domain.RouteRepository = (*routeRepositoryImpl)(nil)
