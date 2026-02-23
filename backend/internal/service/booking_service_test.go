package service

import (
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

type fakeBookingRepo struct{ created bool }

func (f *fakeBookingRepo) Create(_ *domain.Booking) error { f.created = true; return nil }

type fakePriceService struct{}

func (f fakePriceService) FindPrice(_ int64, _ time.Time, _ int) (int64, bool) { return 10000, true }

type fakeHoldService struct{ ok bool }

func (f *fakeHoldService) Hold(_ int64, _ int64, _ int) bool { f.ok = true; return true }

func TestBookingServiceCreate(t *testing.T) {
	svc := NewBookingService(&fakeBookingRepo{}, fakePriceService{}, &fakeHoldService{})
	if err := svc.Create(1, 2, 3, 2); err != nil {
		t.Fatal(err)
	}
}
