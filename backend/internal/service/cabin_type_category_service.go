package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

// CabinTypeCategoryService 提供舱型大类字典的业务能力。
type CabinTypeCategoryService struct {
	repo domain.CabinTypeCategoryRepository
}

func NewCabinTypeCategoryService(repo domain.CabinTypeCategoryRepository) *CabinTypeCategoryService {
	return &CabinTypeCategoryService{repo: repo}
}

func (s *CabinTypeCategoryService) Create(ctx context.Context, category *domain.CabinTypeCategory) error {
	return s.repo.Create(ctx, category)
}

func (s *CabinTypeCategoryService) Update(ctx context.Context, category *domain.CabinTypeCategory) error {
	return s.repo.Update(ctx, category)
}

func (s *CabinTypeCategoryService) GetByID(ctx context.Context, id int64) (*domain.CabinTypeCategory, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CabinTypeCategoryService) List(ctx context.Context) ([]domain.CabinTypeCategory, error) {
	return s.repo.List(ctx)
}

func (s *CabinTypeCategoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
