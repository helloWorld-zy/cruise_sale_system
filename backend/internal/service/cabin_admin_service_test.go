package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

type fakeCabinAdminRepo struct {
	listErr          error
	createErr        error
	updateErr        error
	deleteErr        error
	getInventoryErr  error
	adjustErr        error
	appendLogErr     error
	listPricesErr    error
	upsertPriceErr   error
	listByVoyageData []domain.CabinSKU
	filteredItems    []domain.CabinSKU
	filteredTotal    int64
	inventories      []domain.CabinInventory
	inventory        domain.CabinInventory
	prices           []domain.CabinPrice
	lastLog          *domain.InventoryLog
	lastStatusIDs    []int64
	lastStatus       int16
	lastThresholdSKU int64
	lastThresholdVal int
}

func (f *fakeCabinAdminRepo) CreateSKU(ctx context.Context, s *domain.CabinSKU) error {
	return f.createErr
}
func (f *fakeCabinAdminRepo) UpdateSKU(ctx context.Context, s *domain.CabinSKU) error {
	return f.updateErr
}
func (f *fakeCabinAdminRepo) GetSKUByID(ctx context.Context, id int64) (*domain.CabinSKU, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	return &domain.CabinSKU{ID: id, Code: "SKU-TEST"}, nil
}
func (f *fakeCabinAdminRepo) DeleteSKU(ctx context.Context, id int64) error { return f.deleteErr }
func (f *fakeCabinAdminRepo) ListSKUByVoyage(ctx context.Context, voyageID int64) ([]domain.CabinSKU, error) {
	return f.listByVoyageData, f.listErr
}
func (f *fakeCabinAdminRepo) ListSKUFiltered(ctx context.Context, filter domain.CabinSKUFilter) ([]domain.CabinSKU, int64, error) {
	_ = filter
	return f.filteredItems, f.filteredTotal, f.listErr
}
func (f *fakeCabinAdminRepo) BatchUpdateStatus(ctx context.Context, ids []int64, status int16) error {
	f.lastStatusIDs = ids
	f.lastStatus = status
	return nil
}
func (f *fakeCabinAdminRepo) GetInventoryBySKU(ctx context.Context, skuID int64) (domain.CabinInventory, error) {
	return f.inventory, f.getInventoryErr
}
func (f *fakeCabinAdminRepo) ListAllInventories(ctx context.Context) ([]domain.CabinInventory, error) {
	return f.inventories, f.getInventoryErr
}
func (f *fakeCabinAdminRepo) SetAlertThreshold(ctx context.Context, skuID int64, threshold int) error {
	f.lastThresholdSKU = skuID
	f.lastThresholdVal = threshold
	return nil
}
func (f *fakeCabinAdminRepo) AdjustInventoryAtomic(ctx context.Context, skuID int64, delta int) error {
	return f.adjustErr
}
func (f *fakeCabinAdminRepo) AppendInventoryLog(ctx context.Context, log *domain.InventoryLog) error {
	f.lastLog = log
	return f.appendLogErr
}
func (f *fakeCabinAdminRepo) ListPricesBySKU(ctx context.Context, skuID int64) ([]domain.CabinPrice, error) {
	return f.prices, f.listPricesErr
}
func (f *fakeCabinAdminRepo) UpsertPrice(ctx context.Context, p *domain.CabinPrice) error {
	return f.upsertPriceErr
}

func (f *fakeCabinAdminRepo) BatchSetPrice(ctx context.Context, skuID int64, start, end time.Time, occupancy int, priceCents, childPriceCents, singleSupplementCents int64, priceType string) error {
	_ = ctx
	_ = skuID
	_ = start
	_ = end
	_ = occupancy
	_ = priceCents
	_ = childPriceCents
	_ = singleSupplementCents
	_ = priceType
	return nil
}

func (f *fakeCabinAdminRepo) GetCategoryTree(ctx context.Context) (interface{}, error) {
	return []interface{}{}, nil
}

func TestCabinAdminServiceCrudAndPricing(t *testing.T) {
	repo := &fakeCabinAdminRepo{
		listByVoyageData: []domain.CabinSKU{{ID: 1, Code: "SKU-1"}},
		filteredItems:    []domain.CabinSKU{{ID: 1, Code: "SKU-1"}},
		filteredTotal:    1,
		inventories:      []domain.CabinInventory{{CabinSKUID: 1, Total: 5, Locked: 2, Sold: 2, AlertThreshold: 2}},
		inventory:        domain.CabinInventory{CabinSKUID: 1, Total: 3},
		prices:           []domain.CabinPrice{{CabinSKUID: 1, PriceCents: 9999, Date: time.Now()}},
	}
	svc := NewCabinAdminService(repo)
	ctx := context.Background()

	if _, err := svc.ListByVoyage(ctx, 1); err != nil {
		t.Fatalf("ListByVoyage failed: %v", err)
	}
	if _, total, err := svc.FilteredList(ctx, domain.CabinSKUFilter{VoyageID: 1, Page: 1, PageSize: 10}); err != nil || total != 1 {
		t.Fatalf("FilteredList failed: err=%v total=%d", err, total)
	}
	if err := svc.BatchUpdateStatus(ctx, []int64{1}, 0); err != nil {
		t.Fatalf("BatchUpdateStatus failed: %v", err)
	}
	if err := svc.Create(ctx, &domain.CabinSKU{Code: "SKU-2"}); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if err := svc.Update(ctx, &domain.CabinSKU{ID: 1, Code: "SKU-3"}); err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if err := svc.Delete(ctx, 1); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if _, err := svc.GetInventory(ctx, 1); err != nil {
		t.Fatalf("GetInventory failed: %v", err)
	}
	if _, err := svc.GetAlerts(ctx); err != nil {
		t.Fatalf("GetAlerts failed: %v", err)
	}
	if err := svc.SetAlertThreshold(ctx, 1, 3); err != nil {
		t.Fatalf("SetAlertThreshold failed: %v", err)
	}
	if _, err := svc.ListPrices(ctx, 1); err != nil {
		t.Fatalf("ListPrices failed: %v", err)
	}
	if err := svc.UpsertPrice(ctx, &domain.CabinPrice{CabinSKUID: 1, PriceCents: 12000}); err != nil {
		t.Fatalf("UpsertPrice failed: %v", err)
	}
}

func TestCabinAdminServiceAdjustInventoryBranches(t *testing.T) {
	repo := &fakeCabinAdminRepo{}
	svc := NewCabinAdminService(repo)

	ctx := context.Background()
	if err := svc.AdjustInventory(ctx, 1, -1, "sell"); err != nil {
		t.Fatalf("AdjustInventory should succeed: %v", err)
	}
	if repo.lastLog == nil || repo.lastLog.Change != -1 || repo.lastLog.Reason != "sell" {
		t.Fatalf("expected inventory log to be appended, got %+v", repo.lastLog)
	}

	repo.adjustErr = errors.New("adjust failed")
	if err := svc.AdjustInventory(ctx, 1, -1, "sell"); err == nil {
		t.Fatal("expected adjust error")
	}

	repo.adjustErr = nil
	repo.appendLogErr = errors.New("log failed")
	if err := svc.AdjustInventory(ctx, 1, -1, "sell"); err == nil {
		t.Fatal("expected append log error")
	}
}
