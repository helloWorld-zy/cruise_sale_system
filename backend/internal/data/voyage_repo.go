package data

import (
	"context"
	"cruise_booking_system/internal/model"

	"github.com/google/uuid"
)

type VoyageRepository struct{}

func NewVoyageRepository() *VoyageRepository {
	return &VoyageRepository{}
}

func (r *VoyageRepository) Create(ctx context.Context, voyage *model.Voyage) error {
	return DB.WithContext(ctx).Create(voyage).Error
}

func (r *VoyageRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Voyage, error) {
	var voyage model.Voyage
	if err := DB.WithContext(ctx).First(&voyage, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &voyage, nil
}
