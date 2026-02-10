package core

import (
	"context"
	"cruise_booking_system/internal/data"
	"cruise_booking_system/internal/model"
)

type RecommendationService struct {
	cruiseRepo CruiseRepo // Reusing CruiseRepo interface
}

func NewRecommendationService(repo CruiseRepo) *RecommendationService {
	return &RecommendationService{repo: repo}
}

func (s *RecommendationService) GetRecommendations(ctx context.Context, userID string) ([]model.Cruise, error) {
	// Stub implementation: Just return list of cruises for now
	return s.cruiseRepo.List(ctx, data.CruiseFilter{})
}
