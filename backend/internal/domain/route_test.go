package domain

import "testing"

func TestRouteModelFields(t *testing.T) {
	r := Route{Code: "ASIA-001", Name: "Asia Loop"}
	if r.Code == "" || r.Name == "" {
		t.Fatal("expected route fields")
	}
}

func TestVoyageModelFields(t *testing.T) {
	v := Voyage{Code: "V001", RouteID: 1, CruiseID: 2}
	if v.Code == "" || v.RouteID == 0 {
		t.Fatal("expected voyage fields")
	}
}

func TestCabinSKUModelFields(t *testing.T) {
	s := CabinSKU{Code: "SKU-001", VoyageID: 1, MaxGuests: 2}
	if s.Code == "" || s.MaxGuests == 0 {
		t.Fatal("expected cabin SKU fields")
	}
}

func TestCabinPriceModelFields(t *testing.T) {
	p := CabinPrice{CabinSKUID: 1, Occupancy: 2, PriceCents: 19900}
	if p.PriceCents == 0 {
		t.Fatal("expected price_cents")
	}
}

func TestCabinInventoryModelFields(t *testing.T) {
	inv := CabinInventory{CabinSKUID: 1, Total: 10, Locked: 2, Sold: 1}
	avail := inv.Total - inv.Locked - inv.Sold
	if avail != 7 {
		t.Fatalf("expected 7 available, got %d", avail)
	}
}

func TestInventoryLogModelFields(t *testing.T) {
	log := InventoryLog{CabinSKUID: 1, Change: -2, Reason: "sale"}
	if log.Change != -2 {
		t.Fatal("expected Change=-2")
	}
}

func TestErrInsufficientInventory(t *testing.T) {
	if ErrInsufficientInventory == nil {
		t.Fatal("expected non-nil sentinel error")
	}
}
