package service

import (
	"context"
	"errors"

	"github.com/cruisebooking/backend/internal/domain"
)

// ErrCompanyHasCruises is returned when deleting a company that still has cruises.
var ErrCompanyHasCruises = errors.New("company has cruises")

// CompanyService implements business logic for cruise companies.
type CompanyService struct {
	repo       domain.CompanyRepository
	cruiseRepo domain.CruiseRepository
}

// NewCompanyService creates a CompanyService.
func NewCompanyService(repo domain.CompanyRepository, cruiseRepo domain.CruiseRepository) *CompanyService {
	return &CompanyService{repo: repo, cruiseRepo: cruiseRepo}
}

// Create inserts a new CruiseCompany.
func (s *CompanyService) Create(ctx context.Context, company *domain.CruiseCompany) error {
	return s.repo.Create(ctx, company)
}

// Update saves changes to an existing CruiseCompany.
func (s *CompanyService) Update(ctx context.Context, company *domain.CruiseCompany) error {
	return s.repo.Update(ctx, company)
}

// GetByID retrieves a CruiseCompany by primary key.
func (s *CompanyService) GetByID(ctx context.Context, id int64) (*domain.CruiseCompany, error) {
	return s.repo.GetByID(ctx, id)
}

// List returns a paginated, optionally keyword-filtered list of companies.
func (s *CompanyService) List(ctx context.Context, keyword string, page, pageSize int) ([]domain.CruiseCompany, int64, error) {
	return s.repo.List(ctx, keyword, page, pageSize)
}

// Delete prevents deletion when cruises still reference the company.
// HI-02 FIX: cascade check was previously missing (just had a comment).
func (s *CompanyService) Delete(ctx context.Context, id int64) error {
	_, total, err := s.cruiseRepo.List(ctx, id, 1, 1)
	if err != nil {
		return err
	}
	if total > 0 {
		return ErrCompanyHasCruises
	}
	return s.repo.Delete(ctx, id)
}
