package domain

import "testing"

func TestFacilityExtendedFields(t *testing.T) {
	f := Facility{
		OpenHours:      "08:00-22:00",
		ExtraCharge:    true,
		ChargePriceTip: "参考价 ¥200/人",
		TargetAudience: "家庭,情侣,儿童",
	}
	if f.OpenHours == "" {
		t.Fatal("expected open_hours")
	}
	if !f.ExtraCharge {
		t.Fatal("expected extra_charge")
	}
	if f.ChargePriceTip == "" {
		t.Fatal("expected charge_price_tip")
	}
	if f.TargetAudience == "" {
		t.Fatal("expected target_audience")
	}
}

func TestFacilityCategoryExtendedFields(t *testing.T) {
	fc := FacilityCategory{Status: 1}
	if fc.Status == 0 {
		t.Fatal("expected status")
	}
}
