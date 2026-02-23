package service

import (
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

type PriceRepo interface {
	ListBySKU(skuID int64) ([]domain.CabinPrice, error)
}

type PricingService struct{ repo PriceRepo }

func NewPricingService(repo PriceRepo) *PricingService { return &PricingService{repo: repo} }

func sameDay(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (s *PricingService) FindPrice(skuID int64, date time.Time, occupancy int) (int64, bool) {
	list, _ := s.repo.ListBySKU(skuID)
	for _, v := range list {
		if sameDay(v.Date, date) && v.Occupancy == occupancy {
			return v.PriceCents, true
		}
	}
	return 0, false
}
