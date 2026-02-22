package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// StaffRepository provides database operations for the Staff domain entity.
type StaffRepository struct {
	db *gorm.DB
}

// NewStaffRepository creates a new StaffRepository backed by the given DB.
func NewStaffRepository(db *gorm.DB) *StaffRepository {
	return &StaffRepository{db: db}
}

// Create inserts a new Staff record.
func (r *StaffRepository) Create(ctx context.Context, staff *domain.Staff) error {
	return r.db.WithContext(ctx).Create(staff).Error
}

// GetByUsername retrieves a Staff by username (used for login).
func (r *StaffRepository) GetByUsername(ctx context.Context, username string) (*domain.Staff, error) {
	var staff domain.Staff
	if err := r.db.WithContext(ctx).Where("username = ? AND deleted_at IS NULL", username).First(&staff).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

// GetByID retrieves a Staff by primary key.
func (r *StaffRepository) GetByID(ctx context.Context, id int64) (*domain.Staff, error) {
	var staff domain.Staff
	if err := r.db.WithContext(ctx).First(&staff, id).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

// Update saves changes to an existing Staff.
func (r *StaffRepository) Update(ctx context.Context, staff *domain.Staff) error {
	return r.db.WithContext(ctx).Save(staff).Error
}

// Delete soft-deletes a Staff record.
func (r *StaffRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Staff{}, id).Error
}
