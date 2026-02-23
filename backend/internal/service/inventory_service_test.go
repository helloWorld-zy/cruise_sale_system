package service

import (
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
)

// fakeInventoryRepo satisfies the updated InventoryRepo interface (CRITICAL-01).
type fakeInventoryRepo struct {
	inv  domain.CabinInventory
	logs []domain.InventoryLog
}

func (f *fakeInventoryRepo) GetBySKU(id int64) (domain.CabinInventory, error) { return f.inv, nil }

// AdjustAtomic simulates the atomic SQL update: only applies if total+delta >= 0.
func (f *fakeInventoryRepo) AdjustAtomic(id int64, delta int) error {
	if f.inv.Total+delta < 0 {
		return domain.ErrInsufficientInventory
	}
	f.inv.Total += delta
	return nil
}

func (f *fakeInventoryRepo) AppendLog(log *domain.InventoryLog) error {
	f.logs = append(f.logs, *log)
	return nil
}

func TestInventoryServiceAdjust(t *testing.T) {
	repo := &fakeInventoryRepo{inv: domain.CabinInventory{CabinSKUID: 1, Total: 10}}
	svc := NewInventoryService(repo)
	if err := svc.Adjust(1, -2, "test sale"); err != nil {
		t.Fatal(err)
	}
	if repo.inv.Total != 8 {
		t.Fatalf("expected total 8, got %d", repo.inv.Total)
	}
	if len(repo.logs) != 1 || repo.logs[0].Change != -2 {
		t.Fatal("expected one inventory log entry with Change=-2")
	}
}

func TestInventoryServiceAdjust_InsufficientGuard(t *testing.T) {
	repo := &fakeInventoryRepo{inv: domain.CabinInventory{CabinSKUID: 1, Total: 1}}
	svc := NewInventoryService(repo)
	err := svc.Adjust(1, -5, "oversell")
	if err != domain.ErrInsufficientInventory {
		t.Fatalf("expected ErrInsufficientInventory, got %v", err)
	}
	if repo.inv.Total != 1 {
		t.Fatal("total must not change when adjustment is rejected")
	}
	if len(repo.logs) != 0 {
		t.Fatal("no log should be written when adjustment fails")
	}
}

func TestInventoryServiceAvailable(t *testing.T) {
	repo := &fakeInventoryRepo{inv: domain.CabinInventory{Total: 10, Locked: 2, Sold: 1}}
	svc := NewInventoryService(repo)
	avail, err := svc.Available(1)
	if err != nil || avail != 7 {
		t.Fatalf("expected 7 available, got %d (err=%v)", avail, err)
	}
}
