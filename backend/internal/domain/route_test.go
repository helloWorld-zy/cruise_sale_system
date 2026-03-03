package domain

import "testing"

// TestRouteModelFields 测试 Route 模型字段。
func TestRouteModelFields(t *testing.T) {
	r := Route{Code: "ASIA-001", Name: "Asia Loop"}
	if r.Code == "" || r.Name == "" {
		t.Fatal("expected route fields")
	}
}

// TestRouteExtendedFields 测试 Route 扩展字段。
func TestRouteExtendedFields(t *testing.T) {
	r := Route{
		Code:          "ASIA-001",
		Name:          "亚洲环线",
		DeparturePort: "上海",
		ArrivalPort:   "上海",
		Stops:         "长崎,福冈,济州岛",
		SortOrder:     10,
	}
	if r.DeparturePort == "" {
		t.Fatal("expected departure_port")
	}
	if r.ArrivalPort == "" {
		t.Fatal("expected arrival_port")
	}
	if r.Stops == "" {
		t.Fatal("expected stops")
	}
	if r.SortOrder == 0 {
		t.Fatal("expected sort_order")
	}
}

// TestVoyageModelFields 测试 Voyage 模型字段。
func TestVoyageModelFields(t *testing.T) {
	v := Voyage{Code: "V001", RouteID: 1, CruiseID: 2}
	if v.Code == "" || v.RouteID == 0 {
		t.Fatal("expected voyage fields")
	}
}

// TestCabinSKUModelFields 测试 CabinSKU 模型字段。
func TestCabinSKUModelFields(t *testing.T) {
	s := CabinSKU{Code: "SKU-001", VoyageID: 1, MaxGuests: 2}
	if s.Code == "" || s.MaxGuests == 0 {
		t.Fatal("expected cabin SKU fields")
	}
}

// TestCabinPriceModelFields 测试 CabinPrice 模型字段。
func TestCabinPriceModelFields(t *testing.T) {
	p := CabinPrice{CabinSKUID: 1, Occupancy: 2, PriceCents: 19900}
	if p.PriceCents == 0 {
		t.Fatal("expected price_cents")
	}
}

// TestCabinInventoryModelFields 测试 CabinInventory 模型字段。
func TestCabinInventoryModelFields(t *testing.T) {
	inv := CabinInventory{CabinSKUID: 1, Total: 10, Locked: 2, Sold: 1}
	avail := inv.Total - inv.Locked - inv.Sold
	if avail != 7 {
		t.Fatalf("expected 7 available, got %d", avail)
	}
}

// TestInventoryLogModelFields 测试 InventoryLog 模型字段。
func TestInventoryLogModelFields(t *testing.T) {
	log := InventoryLog{CabinSKUID: 1, Change: -2, Reason: "sale"}
	if log.Change != -2 {
		t.Fatal("expected Change=-2")
	}
}

// TestErrInsufficientInventory 测试 ErrInsufficientInventory 错误。
func TestErrInsufficientInventory(t *testing.T) {
	if ErrInsufficientInventory == nil {
		t.Fatal("expected non-nil sentinel error")
	}
}

// TestPassengerExtendedFields 测试 Passenger 扩展字段。
func TestPassengerExtendedFields(t *testing.T) {
	p := Passenger{
		Name:             "张三",
		EnglishName:      "Zhang San",
		IDType:           "passport",
		IDNumber:         "E12345678",
		Phone:            "13800000001",
		Email:            "zhangsan@example.com",
		EmergencyContact: "李四",
		EmergencyPhone:   "13800000002",
		SpecialNeeds:     "素食,无障碍座位",
		IsFavorite:       true,
	}
	if p.EnglishName == "" {
		t.Fatal("expected english_name")
	}
	if p.Phone == "" {
		t.Fatal("expected phone")
	}
	if p.Email == "" {
		t.Fatal("expected email")
	}
	if p.EmergencyContact == "" {
		t.Fatal("expected emergency_contact")
	}
	if p.EmergencyPhone == "" {
		t.Fatal("expected emergency_phone")
	}
	if p.SpecialNeeds == "" {
		t.Fatal("expected special_needs")
	}
	if !p.IsFavorite {
		t.Fatal("expected is_favorite")
	}
}

// TestOrderStatusTransition 测试订单状态转换。
func TestOrderStatusTransition(t *testing.T) {
	o := Booking{Status: OrderStatusCreated}
	if !o.CanTransitionTo(OrderStatusPendingPayment) {
		t.Fatal("expected valid transition from created to pending_payment")
	}
	if o.CanTransitionTo(OrderStatusCompleted) {
		t.Fatal("expected invalid transition from created to completed")
	}
}

// TestOrderStatusLog 测试订单状态日志。
func TestOrderStatusLog(t *testing.T) {
	log := OrderStatusLog{
		OrderID:    1,
		FromStatus: OrderStatusPendingPayment,
		ToStatus:   OrderStatusPaid,
		OperatorID: 100,
		Remark:     "微信支付成功",
	}
	if log.OperatorID == 0 {
		t.Fatal("expected operator")
	}
	if log.FromStatus != OrderStatusPendingPayment {
		t.Fatal("expected from_status")
	}
	if log.ToStatus != OrderStatusPaid {
		t.Fatal("expected to_status")
	}
}
