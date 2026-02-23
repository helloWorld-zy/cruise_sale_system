package service

import (
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

type BookingRepo interface{ Create(b *domain.Booking) error }
type PriceService interface {
	FindPrice(skuID int64, date time.Time, occupancy int) (int64, bool)
}
type HoldService interface {
	Hold(skuID int64, userID int64, qty int) bool
}

type BookingService struct {
	repo  BookingRepo
	price PriceService
	hold  HoldService
}

func NewBookingService(repo BookingRepo, price PriceService, hold HoldService) *BookingService {
	return &BookingService{repo: repo, price: price, hold: hold}
}

func (s *BookingService) Create(userID, voyageID, skuID int64, guests int) error {
	if !s.hold.Hold(skuID, userID, 1) {
		return nil
	}
	price, _ := s.price.FindPrice(skuID, time.Now(), guests)
	return s.repo.Create(&domain.Booking{UserID: userID, VoyageID: voyageID, CabinSKUID: skuID, Status: "created", TotalCents: price})
}
