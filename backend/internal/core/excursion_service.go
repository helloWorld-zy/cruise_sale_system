package core

import (
	"context"
	"cruise_booking_system/internal/data"
	"cruise_booking_system/internal/model"

	"github.com/google/uuid"
)

type ExcursionService struct {
	repo *data.ExcursionRepository
}

func NewExcursionService(repo *data.ExcursionRepository) *ExcursionService {
	return &ExcursionService{repo: repo}
}

func (s *ExcursionService) ListExcursions(ctx context.Context, cruiseIDStr string) ([]model.Excursion, error) {
	id, err := uuid.Parse(cruiseIDStr)
	if err != nil {
		return nil, err
	}
	return s.repo.ListByCruiseID(ctx, id)
}
