package repository

import (
	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

type BookingRepository struct{ db *gorm.DB }

func NewBookingRepository(db *gorm.DB) *BookingRepository { return &BookingRepository{db: db} }

func (r *BookingRepository) Create(b *domain.Booking) error { return r.db.Create(b).Error }
func (r *BookingRepository) ListByUser(userID int64) ([]domain.Booking, error) {
	var out []domain.Booking
	return out, r.db.Where("user_id = ?", userID).Order("id desc").Find(&out).Error
}
