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
