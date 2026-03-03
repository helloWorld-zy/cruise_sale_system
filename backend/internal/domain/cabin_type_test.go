package domain

import "testing"

func TestCabinTypeExtendedFields(t *testing.T) {
	ct := CabinType{
		Code:         "BALCONY",
		AreaMin:      28.0,
		AreaMax:      35.0,
		MaxCapacity:  4,
		BedType:      "双人床/可拆分为两张单人床",
		Tags:         "落地窗,私人阳台,管家服务",
		Amenities:    "独立卫浴,吹风机,保险箱,迷你吧,电视",
		FloorPlanURL: "https://example.com/plan.jpg",
	}
	if ct.Code == "" {
		t.Fatal("expected code")
	}
	if ct.AreaMin == 0 {
		t.Fatal("expected area_min")
	}
	if ct.AreaMax == 0 {
		t.Fatal("expected area_max")
	}
	if ct.MaxCapacity == 0 {
		t.Fatal("expected max_capacity")
	}
	if ct.BedType == "" {
		t.Fatal("expected bed_type")
	}
	if ct.Tags == "" {
		t.Fatal("expected tags")
	}
	if ct.Amenities == "" {
		t.Fatal("expected amenities")
	}
	if ct.FloorPlanURL == "" {
		t.Fatal("expected floor_plan_url")
	}
}
