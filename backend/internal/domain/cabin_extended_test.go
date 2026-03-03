package domain

import "testing"

func TestCabinSKUExtendedFields(t *testing.T) {
	sku := CabinSKU{
		Position:    "mid",
		Orientation: "port",
		HasWindow:   true,
		HasBalcony:  true,
		BedType:     "大床",
		Amenities:   "卫浴,迷你吧,电视",
		Grade:       "premium",
		CabinTypeID: 1,
	}
	if sku.Position == "" {
		t.Fatal("expected position")
	}
	if sku.Orientation == "" {
		t.Fatal("expected orientation")
	}
	if !sku.HasWindow {
		t.Fatal("expected has_window")
	}
	if !sku.HasBalcony {
		t.Fatal("expected has_balcony")
	}
	if sku.BedType == "" {
		t.Fatal("expected bed_type")
	}
	if sku.Amenities == "" {
		t.Fatal("expected amenities")
	}
	if sku.Grade == "" {
		t.Fatal("expected grade")
	}
}

func TestCabinInventoryAlertThreshold(t *testing.T) {
	inv := CabinInventory{Total: 10, Locked: 3, Sold: 5, AlertThreshold: 5}
	available := inv.Total - inv.Locked - inv.Sold
	if available >= inv.AlertThreshold {
		t.Fatal("expected alert triggered")
	}
}

func TestCabinPriceType(t *testing.T) {
	p := CabinPrice{PriceType: "child", PriceCents: 5000}
	if p.PriceType == "" {
		t.Fatal("expected price_type")
	}
}

func TestCabinPriceExtendedFields(t *testing.T) {
	p := CabinPrice{
		CabinSKUID:            1,
		PriceCents:            199900,
		ChildPriceCents:       99900,
		SingleSupplementCents: 50000,
		PriceType:             "earlybird",
	}
	if p.ChildPriceCents == 0 {
		t.Fatal("expected child_price_cents")
	}
	if p.SingleSupplementCents == 0 {
		t.Fatal("expected single_supplement_cents")
	}
	if p.PriceType == "" {
		t.Fatal("expected price_type")
	}
}
