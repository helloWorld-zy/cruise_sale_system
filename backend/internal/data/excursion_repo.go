package data

import (
	"context"
	"cruise_booking_system/internal/model"

	"github.com/google/uuid"
)

type ExcursionRepository struct{}

func NewExcursionRepository() *ExcursionRepository {
	return &ExcursionRepository{}
}

func (r *ExcursionRepository) ListByCruiseID(ctx context.Context, cruiseID uuid.UUID) ([]model.Excursion, error) {
	var excursions []model.Excursion
	if err := DB.WithContext(ctx).Where("cruise_id = ?", cruiseID).Find(&excursions).Error; err != nil {
		return nil, err
	}
	return excursions, nil
}
