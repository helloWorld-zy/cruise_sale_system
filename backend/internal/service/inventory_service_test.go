package service

import (
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
)

type fakeInventoryRepo struct{ inv domain.CabinInventory }

func (f *fakeInventoryRepo) GetBySKU(id int64) (domain.CabinInventory, error) { return f.inv, nil }
func (f *fakeInventoryRepo) Update(inv *domain.CabinInventory) error          { f.inv = *inv; return nil }

func TestInventoryServiceAdjust(t *testing.T) {
	repo := &fakeInventoryRepo{inv: domain.CabinInventory{CabinSKUID: 1, Total: 10, Locked: 0, Sold: 0}}
	svc := NewInventoryService(repo)
	if err := svc.Adjust(1, -2); err != nil {
		t.Fatal(err)
	}
	if repo.inv.Total != 8 {
		t.Fatal("expected total 8")
	}
}
