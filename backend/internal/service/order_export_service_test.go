package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

type fakeOrderExportRepo struct {
	orders []domain.Booking
	err    error
	filter OrderFilter
}

func (f *fakeOrderExportRepo) ListWithFilter(ctx context.Context, filter OrderFilter) ([]domain.Booking, error) {
	_ = ctx
	f.filter = filter
	if f.err != nil {
		return nil, f.err
	}
	return f.orders, nil
}

func TestOrderExportServiceDeniedWithoutPermission(t *testing.T) {
	svc := NewOrderExportService(&fakeOrderExportRepo{})

	_, err := svc.ExportToExcel(context.Background(), OrderFilter{})
	if !errors.Is(err, ErrOrderExportForbidden) {
		t.Fatalf("expected ErrOrderExportForbidden, got %v", err)
	}
}

func TestOrderExportServiceRejectsOverLimit(t *testing.T) {
	orders := make([]domain.Booking, defaultOrderExportMaxRows+1)
	for i := range orders {
		orders[i] = domain.Booking{ID: int64(i + 1), Status: domain.OrderStatusCreated}
	}
	svc := NewOrderExportService(&fakeOrderExportRepo{orders: orders})

	ctx := WithOrderExportPermission(context.Background(), true)
	_, err := svc.ExportToExcel(ctx, OrderFilter{})
	if !errors.Is(err, ErrOrderExportExceededLimit) {
		t.Fatalf("expected ErrOrderExportExceededLimit, got %v", err)
	}
}

func TestOrderExportServiceSanitizesCSVInjection(t *testing.T) {
	now := time.Date(2026, 3, 3, 10, 0, 0, 0, time.UTC)
	svc := NewOrderExportService(&fakeOrderExportRepo{orders: []domain.Booking{
		{ID: 1, UserID: 10, VoyageID: 20, CabinSKUID: 30, Status: "=HYPERLINK(\"bad\")", TotalCents: 9900, PaidCents: 0, CreatedAt: now},
	}})

	ctx := WithOrderExportPermission(context.Background(), true)
	content, err := svc.ExportToExcel(ctx, OrderFilter{})
	if err != nil {
		t.Fatalf("expected export success, got %v", err)
	}
	text := string(content)
	if strings.Contains(text, "Excel content placeholder") {
		t.Fatal("placeholder content should not appear")
	}
	if !strings.Contains(text, "'=HYPERLINK") {
		t.Fatalf("expected CSV injection value to be prefixed with quote, got: %s", text)
	}
}

func TestOrderExportServiceSanitizesTabAndCR(t *testing.T) {
	now := time.Date(2026, 3, 3, 10, 0, 0, 0, time.UTC)
	svc := NewOrderExportService(&fakeOrderExportRepo{orders: []domain.Booking{
		{ID: 1, UserID: 10, VoyageID: 20, CabinSKUID: 30, Status: "\tcmd", TotalCents: 100, PaidCents: 0, CreatedAt: now},
	}})

	ctx := WithOrderExportPermission(context.Background(), true)
	content, err := svc.ExportToExcel(ctx, OrderFilter{})
	if err != nil {
		t.Fatalf("expected export success, got %v", err)
	}
	text := string(content)
	if !strings.Contains(text, "'\tcmd") {
		t.Fatalf("expected tab-prefixed value to be sanitized, got: %s", text)
	}
}

func TestOrderExportServicePassesExtendedFiltersAndRendersDerivedColumns(t *testing.T) {
	repo := &fakeOrderExportRepo{orders: []domain.Booking{{
		ID:         11,
		UserID:     22,
		VoyageID:   33,
		CabinSKUID: 44,
		Status:     domain.OrderStatusPaid,
		TotalCents: 18800,
		PaidCents:  18800,
		BookingNo:  "BK-11",
		Phone:      "13800000001",
		VoyageCode: "VOY-ALPHA",
		CruiseName: "海洋量子号",
		CreatedAt:  time.Date(2026, 3, 7, 8, 0, 0, 0, time.UTC),
	}}}
	svc := NewOrderExportService(repo)
	startDate := "2026-03-01"
	endDate := "2026-03-31"
	ctx := WithOrderExportPermission(context.Background(), true)

	content, err := svc.ExportToExcel(ctx, OrderFilter{
		Status:     domain.OrderStatusPaid,
		Phone:      "138",
		VoyageCode: "VOY",
		CruiseName: "量子",
		Keyword:    "BK-11",
		StartDate:  &startDate,
		EndDate:    &endDate,
	})
	if err != nil {
		t.Fatalf("expected export success, got %v", err)
	}
	if repo.filter.Phone != "138" || repo.filter.VoyageCode != "VOY" || repo.filter.CruiseName != "量子" || repo.filter.Keyword != "BK-11" {
		t.Fatalf("unexpected filter pass-through: %+v", repo.filter)
	}
	text := string(content)
	if !strings.Contains(text, "booking_no,phone,voyage_code,cruise_name") {
		t.Fatalf("expected derived export headers, got: %s", text)
	}
	if !strings.Contains(text, "BK-11,13800000001,VOY-ALPHA,海洋量子号") {
		t.Fatalf("expected derived order fields in export content, got: %s", text)
	}
}
