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

// ListByCruise 查询指定邮轮下的所有设施。
func (s *FacilityService) ListByCruise(ctx context.Context, cruiseID int64) ([]domain.Facility, error) {
	return s.repo.ListByCruise(ctx, cruiseID)
}

// Delete 删除指定的设施。
func (s *FacilityService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
