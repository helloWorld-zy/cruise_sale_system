package service

import (
	"context"
	"errors"

	"github.com/cruisebooking/backend/internal/domain"
)

// ErrCruiseHasCabins is returned when deleting a cruise that still has cabin types.
var ErrCruiseHasCabins = errors.New("cruise has cabins")

// CruiseService implements business logic for cruises.
type CruiseService struct {
	cruiseRepo  domain.CruiseRepository
	cabinRepo   domain.CabinTypeRepository
	companyRepo domain.CompanyRepository
}

// NewCruiseService creates a CruiseService with injected repositories.
func NewCruiseService(cruiseRepo domain.CruiseRepository, cabinRepo domain.CabinTypeRepository, companyRepo domain.CompanyRepository) *CruiseService {
	return &CruiseService{cruiseRepo: cruiseRepo, cabinRepo: cabinRepo, companyRepo: companyRepo}
}

// Create validates that the company exists before inserting the cruise.
func (s *CruiseService) Create(ctx context.Context, cruise *domain.Cruise) error {
	if _, err := s.companyRepo.GetByID(ctx, cruise.CompanyID); err != nil {
		return err
	}
	return s.cruiseRepo.Create(ctx, cruise)
}

// Update saves changes to an existing Cruise.
func (s *CruiseService) Update(ctx context.Context, cruise *domain.Cruise) error {
	return s.cruiseRepo.Update(ctx, cruise)
}

// GetByID retrieves a Cruise by primary key.
func (s *CruiseService) GetByID(ctx context.Context, id int64) (*domain.Cruise, error) {
	return s.cruiseRepo.GetByID(ctx, id)
}

// List returns a paginated list of cruises, optionally filtered by company.
func (s *CruiseService) List(ctx context.Context, companyID int64, page, pageSize int) ([]domain.Cruise, int64, error) {
	return s.cruiseRepo.List(ctx, companyID, page, pageSize)
}

// Delete prevents deletion when cabin types still reference the cruise.
func (s *CruiseService) Delete(ctx context.Context, id int64) error {
	cabins, total, err := s.cabinRepo.ListByCruise(ctx, id, 1, 1)
	if err != nil {
		return err
	}
	if total > 0 || len(cabins) > 0 {
		return ErrCruiseHasCabins
	}
	return s.cruiseRepo.Delete(ctx, id)
}
