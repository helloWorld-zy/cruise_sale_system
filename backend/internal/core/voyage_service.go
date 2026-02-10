package core

import (
	"context"
	"cruise_booking_system/internal/data"
	"cruise_booking_system/internal/model"

	"github.com/google/uuid"
)

type VoyageService struct {
	repo *data.VoyageRepository
}

func NewVoyageService(repo *data.VoyageRepository) *VoyageService {
	return &VoyageService{repo: repo}
}

func (s *VoyageService) CreateVoyage(ctx context.Context, voyage *model.Voyage) error {
	voyage.ID = uuid.New()
	return s.repo.Create(ctx, voyage)
}
