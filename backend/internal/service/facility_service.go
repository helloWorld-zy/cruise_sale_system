package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

// FacilityService 实现设施相关的业务逻辑。
type FacilityService struct {
	repo domain.FacilityRepository // 设施数据仓储
}

// NewFacilityService 创建设施服务实例。
func NewFacilityService(repo domain.FacilityRepository) *FacilityService {
	return &FacilityService{repo: repo}
}

// Create 创建新的设施。
func (s *FacilityService) Create(ctx context.Context, f *domain.Facility) error {
	return s.repo.Create(ctx, f)
}

// Update 更新设施信息。
func (s *FacilityService) Update(ctx context.Context, f *domain.Facility) error {
	return s.repo.Update(ctx, f)
}

// GetByID 查询单个设施详情。
func (s *FacilityService) GetByID(ctx context.Context, id int64) (*domain.Facility, error) {
	return s.repo.GetByID(ctx, id)
}

// ListByCruise 查询指定邮轮下的所有设施。
func (s *FacilityService) ListByCruise(ctx context.Context, cruiseID int64) ([]domain.Facility, error) {
	return s.repo.ListByCruise(ctx, cruiseID)
}

// ListByCruiseAndCategory 按邮轮和分类筛选设施列表。
func (s *FacilityService) ListByCruiseAndCategory(ctx context.Context, cruiseID, categoryID int64) ([]domain.Facility, error) {
	return s.repo.ListByCruiseAndCategory(ctx, cruiseID, categoryID)
}

// Delete 删除指定的设施。
func (s *FacilityService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
