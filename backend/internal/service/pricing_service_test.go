package service

import (
	"context"
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

type fakePriceRepo struct {
	prices  []domain.CabinPrice
	created int
}

func (f *fakePriceRepo) ListBySKU(ctx context.Context, skuID int64) ([]domain.CabinPrice, error) {
	return f.prices, nil
}

func (f *fakePriceRepo) Create(ctx context.Context, price *domain.CabinPrice) error {
	f.created++
	return nil
}

func TestPricingServiceFindPrice(t *testing.T) {
	d := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	svc := NewPricingService(&fakePriceRepo{prices: []domain.CabinPrice{{CabinSKUID: 1, Date: d, Occupancy: 2, PriceCents: 19900}}})
	p, ok, err := svc.FindPrice(context.Background(), 1, d, 2)
	if err != nil || !ok || p != 19900 {
		t.Fatalf("expected price 19900 ok=true err=nil, got p=%d ok=%v err=%v", p, ok, err)
	}
}

func TestPricingServiceFindPrice_NotFound(t *testing.T) {
	svc := NewPricingService(&fakePriceRepo{prices: nil})
	_, ok, err := svc.FindPrice(context.Background(), 1, time.Now(), 2)
	if err != nil || ok {
		t.Fatal("expected ok=false, err=nil for empty repo")
	}
}

// HIGH-01: 验证时区安全性 — 不同时区中的相同日历日期必须匹配。
func TestPricingServiceFindPrice_TimezoneSafety(t *testing.T) {
	// 存储在 UTC 午夜
	utcDate := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	// 从 CST 客户端查询 (UTC+8)，即 2026-05-01 08:00 CST = 2026-05-01 00:00 UTC
	cst := time.FixedZone("CST", 8*3600)
	cstDate := time.Date(2026, 5, 1, 8, 0, 0, 0, cst)

	svc := NewPricingService(&fakePriceRepo{prices: []domain.CabinPrice{
		{CabinSKUID: 1, Date: utcDate, Occupancy: 2, PriceCents: 19900},
	}})

	p, ok, err := svc.FindPrice(context.Background(), 1, cstDate, 2)
	if err != nil || !ok || p != 19900 {
		t.Fatalf("timezone-safe lookup failed: p=%d ok=%v err=%v", p, ok, err)
	}
}

// HIGH-01 边缘情况：CST 午夜实际上是 UTC 的前一天 — 必须不匹配。
func TestPricingServiceFindPrice_TimezoneMismatch(t *testing.T) {
	utcDate := time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC)
	cst := time.FixedZone("CST", 8*3600)
	// 2026-05-01 00:00 CST = 2026-04-30 16:00 UTC → 不同的 UTC 日期
	cstDate := time.Date(2026, 5, 1, 0, 0, 0, 0, cst)

	svc := NewPricingService(&fakePriceRepo{prices: []domain.CabinPrice{
		{CabinSKUID: 1, Date: utcDate, Occupancy: 2, PriceCents: 19900},
	}})

	_, ok, err := svc.FindPrice(context.Background(), 1, cstDate, 2)
	if err != nil || ok {
		t.Fatal("expected no match for different UTC dates")
	}
}

func TestPricingServiceFindPriceWithType(t *testing.T) {
	d := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	svc := NewPricingService(&fakePriceRepo{prices: []domain.CabinPrice{
		{CabinSKUID: 1, Date: d, Occupancy: 2, PriceCents: 19900, ChildPriceCents: 9900, PriceType: "base"},
		{CabinSKUID: 1, Date: d, Occupancy: 2, PriceCents: 17900, ChildPriceCents: 8900, PriceType: "earlybird"},
	}})
	result, ok, err := svc.FindPriceByType(context.Background(), 1, d, 2, "earlybird")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected to find earlybird price")
	}
	if result.PriceCents != 17900 {
		t.Fatalf("expected earlybird price 17900, got %d", result.PriceCents)
	}
	if result.ChildPriceCents != 8900 {
		t.Fatalf("expected child price 8900, got %d", result.ChildPriceCents)
	}
}

func TestPricingServiceBatchSetPrice(t *testing.T) {
	repo := &fakePriceRepo{}
	svc := NewPricingService(repo)
	startDate := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 5, 3, 0, 0, 0, 0, time.UTC)
	err := svc.BatchSetPrice(context.Background(), 1, startDate, endDate, 2, 19900, 9900, 5000, "base")
	if err != nil {
		t.Fatal(err)
	}
	if repo.created != 3 {
		t.Fatalf("expected 3 prices created, got %d", repo.created)
	}
}
