package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

// CabinTypeService 实现舱房类型相关的业务逻辑。
// 提供舱房类型的 CRUD 操作。
type CabinTypeService struct {
	repo domain.CabinTypeRepository // 舱房类型数据仓储
}

// NewCabinTypeService 创建舱房类型服务实例，通过依赖注入传入仓储。
func NewCabinTypeService(repo domain.CabinTypeRepository) *CabinTypeService {
	return &CabinTypeService{repo: repo}
}

// Create 创建新的舱房类型。
func (s *CabinTypeService) Create(ctx context.Context, ct *domain.CabinType) error {
	return s.repo.Create(ctx, ct)
}

// Update 保存对已有舱房类型的修改。
func (s *CabinTypeService) Update(ctx context.Context, ct *domain.CabinType) error {
	return s.repo.Update(ctx, ct)
}

// GetByID 根据主键查询舱房类型详情。
func (s *CabinTypeService) GetByID(ctx context.Context, id int64) (*domain.CabinType, error) {
	return s.repo.GetByID(ctx, id)
}

// List 返回指定邮轮下的分页舱房类型列表。
func (s *CabinTypeService) List(ctx context.Context, cruiseID int64, page, pageSize int) ([]domain.CabinType, int64, error) {
	return s.repo.ListByCruise(ctx, cruiseID, page, pageSize)
}

// Delete 删除指定的舱房类型。
func (s *CabinTypeService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
