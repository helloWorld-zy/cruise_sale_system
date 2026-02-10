package data

import (
	"context"
	"cruise_booking_system/internal/model"

	"github.com/google/uuid"
)

type CruiseFilter struct {
	Destination string
	Date        string // Simple date string check for now (usually implies checking Voyages)
}

type CruiseRepository struct{}

func NewCruiseRepository() *CruiseRepository {
	return &CruiseRepository{}
}

func (r *CruiseRepository) List(ctx context.Context, filter CruiseFilter) ([]model.Cruise, error) {
	var cruises []model.Cruise
	query := DB.WithContext(ctx).Model(&model.Cruise{})

	// Filter logic would go here. 
	// Note: filtering by "Destination" or "Date" usually requires joining with Routes/Voyages 
	// which are not yet fully implemented. For MVP US1, we'll just list active cruises.
	query = query.Where("status = ?", model.CruiseStatusActive)

	if err := query.Find(&cruises).Error; err != nil {
		return nil, err
	}
	return cruises, nil
}

func (r *CruiseRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Cruise, error) {
	var cruise model.Cruise
	if err := DB.WithContext(ctx).First(&cruise, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &cruise, nil
}

func (r *CruiseRepository) GetCabinTypes(ctx context.Context, cruiseID uuid.UUID) ([]model.CabinType, error) {
	var cabinTypes []model.CabinType
	if err := DB.WithContext(ctx).Where("cruise_id = ?", cruiseID).Find(&cabinTypes).Error; err != nil {
		return nil, err
	}
	return cabinTypes, nil
}

func (r *CruiseRepository) Create(ctx context.Context, cruise *model.Cruise) error {
	return DB.WithContext(ctx).Create(cruise).Error
}

func (r *CruiseRepository) Update(ctx context.Context, cruise *model.Cruise) error {
	return DB.WithContext(ctx).Save(cruise).Error
}

func (r *CruiseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return DB.WithContext(ctx).Delete(&model.Cruise{}, "id = ?", id).Error
}
