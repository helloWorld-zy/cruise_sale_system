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
	p, ok, err := svc.FindPrice(1, d, 2)
	if err != nil || !ok || p != 19900 {
		t.Fatalf("expected price 19900 ok=true err=nil, got p=%d ok=%v err=%v", p, ok, err)
	}
}

func TestPricingServiceFindPrice_NotFound(t *testing.T) {
	svc := NewPricingService(fakePriceRepo{prices: nil})
	_, ok, err := svc.FindPrice(1, time.Now(), 2)
	if err != nil || ok {
		t.Fatal("expected ok=false, err=nil for empty repo")
	}
}

// HIGH-01: Verify timezone safety — same calendar date in different timezones must match.
func TestPricingServiceFindPrice_TimezoneSafety(t *testing.T) {
	// stored in UTC midnight
	utcDate := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	// queried from a CST client (UTC+8), which is 2026-05-01 08:00 CST = 2026-05-01 00:00 UTC
	cst := time.FixedZone("CST", 8*3600)
	cstDate := time.Date(2026, 5, 1, 8, 0, 0, 0, cst)

	svc := NewPricingService(fakePriceRepo{prices: []domain.CabinPrice{
		{CabinSKUID: 1, Date: utcDate, Occupancy: 2, PriceCents: 19900},
	}})

	p, ok, err := svc.FindPrice(1, cstDate, 2)
	if err != nil || !ok || p != 19900 {
		t.Fatalf("timezone-safe lookup failed: p=%d ok=%v err=%v", p, ok, err)
	}
}

// HIGH-01 edge case: CST midnight is actually UTC previous day — must NOT match.
func TestPricingServiceFindPrice_TimezoneMismatch(t *testing.T) {
	utcDate := time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC)
	cst := time.FixedZone("CST", 8*3600)
	// 2026-05-01 00:00 CST = 2026-04-30 16:00 UTC → different UTC day
	cstDate := time.Date(2026, 5, 1, 0, 0, 0, 0, cst)

	svc := NewPricingService(fakePriceRepo{prices: []domain.CabinPrice{
		{CabinSKUID: 1, Date: utcDate, Occupancy: 2, PriceCents: 19900},
	}})

	_, ok, err := svc.FindPrice(1, cstDate, 2)
	if err != nil || ok {
		t.Fatal("expected no match for different UTC dates")
	}
}
