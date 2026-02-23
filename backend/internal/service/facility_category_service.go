package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

// FacilityCategoryService 实现设施分类相关的业务逻辑。
type FacilityCategoryService struct {
	repo domain.FacilityCategoryRepository // 设施分类数据仓储
}

// NewFacilityCategoryService 创建设施分类服务实例。
func NewFacilityCategoryService(repo domain.FacilityCategoryRepository) *FacilityCategoryService {
	return &FacilityCategoryService{repo: repo}
}

// Create 创建新的设施分类。
func (s *FacilityCategoryService) Create(ctx context.Context, cat *domain.FacilityCategory) error {
	return s.repo.Create(ctx, cat)
}

// List 查询所有设施分类，按排序权重排列。
func (s *FacilityCategoryService) List(ctx context.Context) ([]domain.FacilityCategory, error) {
	return s.repo.List(ctx)
}

// Delete 删除设施分类。注意：调用方应先确认无设施引用此分类。
func (s *FacilityCategoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
