package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

// FacilityCategoryService implements business logic for facility categories.
type FacilityCategoryService struct {
	repo domain.FacilityCategoryRepository
}

// NewFacilityCategoryService creates a FacilityCategoryService.
func NewFacilityCategoryService(repo domain.FacilityCategoryRepository) *FacilityCategoryService {
	return &FacilityCategoryService{repo: repo}
}

// Create inserts a new FacilityCategory.
func (s *FacilityCategoryService) Create(ctx context.Context, cat *domain.FacilityCategory) error {
	return s.repo.Create(ctx, cat)
}

// List returns all facility categories ordered by sort_order.
func (s *FacilityCategoryService) List(ctx context.Context) ([]domain.FacilityCategory, error) {
	return s.repo.List(ctx)
}

// Delete removes a FacilityCategory. Note: caller should verify no facilities reference this category.
func (s *FacilityCategoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
