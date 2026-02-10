package core

import (
	"context"
	"cruise_booking_system/internal/data"
	"cruise_booking_system/internal/model"

	"github.com/google/uuid"
)

type CruiseRepo interface {
	List(ctx context.Context, filter data.CruiseFilter) ([]model.Cruise, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Cruise, error)
	GetCabinTypes(ctx context.Context, cruiseID uuid.UUID) ([]model.CabinType, error)
	Create(ctx context.Context, cruise *model.Cruise) error
	Update(ctx context.Context, cruise *model.Cruise) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CruiseService struct {
	repo CruiseRepo
}

func NewCruiseService(repo CruiseRepo) *CruiseService {
	return &CruiseService{repo: repo}
}

func (s *CruiseService) ListCruises(ctx context.Context, destination, date string) ([]model.Cruise, error) {
	filter := data.CruiseFilter{
		Destination: destination,
		Date:        date,
	}
	return s.repo.List(ctx, filter)
}

func (s *CruiseService) CreateCruise(ctx context.Context, cruise *model.Cruise) error {
	cruise.ID = uuid.New() // Ensure ID is generated
	return s.repo.Create(ctx, cruise)
}

func (s *CruiseService) UpdateCruise(ctx context.Context, cruise *model.Cruise) error {
	return s.repo.Update(ctx, cruise)
}

func (s *CruiseService) DeleteCruise(ctx context.Context, idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

type CruiseDetail struct {
	*model.Cruise
	CabinTypes []model.CabinType `json:"cabin_types"`
}

func (s *CruiseService) GetCruiseDetail(ctx context.Context, idStr string) (*CruiseDetail, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	cruise, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	cabinTypes, err := s.repo.GetCabinTypes(ctx, id)
	if err != nil {
		return nil, err // Or just log and return empty list
	}

	return &CruiseDetail{
		Cruise:     cruise,
		CabinTypes: cabinTypes,
	}, nil
}
