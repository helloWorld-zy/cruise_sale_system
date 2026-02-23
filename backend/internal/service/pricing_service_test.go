package service

import (
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

type fakePriceRepo struct{ prices []domain.CabinPrice }

func (f fakePriceRepo) ListBySKU(skuID int64) ([]domain.CabinPrice, error) { return f.prices, nil }

func TestPricingServiceFindPrice(t *testing.T) {
	d := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	svc := NewPricingService(fakePriceRepo{prices: []domain.CabinPrice{{CabinSKUID: 1, Date: d, Occupancy: 2, PriceCents: 19900}}})
	p, ok := svc.FindPrice(1, d, 2)
	if !ok || p != 19900 {
		t.Fatal("expected price 19900")
	}
}
