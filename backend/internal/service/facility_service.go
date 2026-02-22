package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

// FacilityService implements business logic for facilities.
type FacilityService struct {
	repo domain.FacilityRepository
}

// NewFacilityService creates a FacilityService.
func NewFacilityService(repo domain.FacilityRepository) *FacilityService {
	return &FacilityService{repo: repo}
}

// Create inserts a new Facility.
func (s *FacilityService) Create(ctx context.Context, f *domain.Facility) error {
	return s.repo.Create(ctx, f)
}

// ListByCruise returns all facilities for the given cruise.
func (s *FacilityService) ListByCruise(ctx context.Context, cruiseID int64) ([]domain.Facility, error) {
	return s.repo.ListByCruise(ctx, cruiseID)
}

// Delete removes a Facility by ID.
func (s *FacilityService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
