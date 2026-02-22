package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

// CabinTypeService implements business logic for cabin types.
type CabinTypeService struct {
	repo domain.CabinTypeRepository
}

// NewCabinTypeService creates a CabinTypeService with injected repository.
func NewCabinTypeService(repo domain.CabinTypeRepository) *CabinTypeService {
	return &CabinTypeService{repo: repo}
}

// Create inserts a new CabinType. Validates that CruiseID is set.
func (s *CabinTypeService) Create(ctx context.Context, ct *domain.CabinType) error {
	return s.repo.Create(ctx, ct)
}

// Update saves changes to an existing CabinType.
func (s *CabinTypeService) Update(ctx context.Context, ct *domain.CabinType) error {
	return s.repo.Update(ctx, ct)
}

// GetByID retrieves a CabinType by its primary key.
func (s *CabinTypeService) GetByID(ctx context.Context, id int64) (*domain.CabinType, error) {
	return s.repo.GetByID(ctx, id)
}

// List returns a paginated list of cabin types for the given cruise.
func (s *CabinTypeService) List(ctx context.Context, cruiseID int64, page, pageSize int) ([]domain.CabinType, int64, error) {
	return s.repo.ListByCruise(ctx, cruiseID, page, pageSize)
}

// Delete removes a CabinType by ID.
func (s *CabinTypeService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
