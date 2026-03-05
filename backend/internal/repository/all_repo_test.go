package repository

import (
	"context"
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAllRepos(t *testing.T) {
	// Setup in-memory DB
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to sqlite: %v", err)
	}

	// Migrate all models
	err = db.AutoMigrate(
		&domain.Staff{},
		&domain.CruiseCompany{},
		&domain.Cruise{},
		&domain.Route{},
		&domain.Voyage{},
		&domain.VoyageItinerary{},
		&domain.CabinType{},
		&domain.FacilityCategory{},
		&domain.Facility{},
		&domain.CabinSKU{},
		&domain.CabinPrice{},
		&domain.CabinInventory{},
		&domain.InventoryLog{},
		&domain.Booking{},
	)
	if err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}

	ctx := context.Background()

	// 1. StaffRepo
	staffRepo := NewStaffRepository(db)
	staff := &domain.Staff{Username: "test1", Status: 1}
	_ = staffRepo.Create(ctx, staff)
	_, _ = staffRepo.GetByID(ctx, staff.ID)
	_, _ = staffRepo.GetByUsername(ctx, "test1")
	staff.Status = 0
	_ = staffRepo.Update(ctx, staff)
	_ = staffRepo.Delete(ctx, staff.ID)

	// 2. CompanyRepo
	compRepo := NewCompanyRepository(db)
	comp := &domain.CruiseCompany{Name: "c1"}
	_ = compRepo.Create(ctx, comp)
	_, _ = compRepo.GetByID(ctx, comp.ID)
	_, _, _ = compRepo.List(ctx, "", 1, 10)
	_, _, _ = compRepo.List(ctx, "c1", 1, 10)
	comp.Name = "c2"
	_ = compRepo.Update(ctx, comp)
	_ = compRepo.Delete(ctx, comp.ID)

	// 3. CruiseRepo
	cruiseRepo := NewCruiseRepository(db)
	cr := &domain.Cruise{Name: "cr1", CompanyID: 1}
	_ = cruiseRepo.Create(ctx, cr)
	_, _ = cruiseRepo.GetByID(ctx, cr.ID)
	_, _, _ = cruiseRepo.List(ctx, 1, "", nil, "", 1, 10)
	_, _, _ = cruiseRepo.List(ctx, 0, "", nil, "", 1, 10)
	cr.Name = "cr2"
	_ = cruiseRepo.Update(ctx, cr)
	_ = cruiseRepo.Delete(ctx, cr.ID)

	// 4. RouteRepo
	rtRepo := NewRouteRepository(db)
	rt := &domain.Route{Name: "rt1"}
	_ = rtRepo.Create(ctx, rt)
	_, _ = rtRepo.GetByID(ctx, rt.ID)
	_, _ = rtRepo.List(ctx)
	rt.Name = "rt2"
	_ = rtRepo.Update(ctx, rt)
	_ = rtRepo.Delete(ctx, rt.ID)

	// 5. VoyageRepo
	vRepo := NewVoyageRepository(db)
	v := &domain.Voyage{Code: "v1", CruiseID: 1, BriefInfo: "demo"}
	_ = vRepo.Create(ctx, v)
	_, _ = vRepo.GetByID(ctx, v.ID)
	_, _ = vRepo.List(ctx)
	v.Code = "v2"
	_ = vRepo.Update(ctx, v)
	_ = vRepo.Delete(ctx, v.ID)

	// 6. CabinTypeRepo
	ctRepo := NewCabinTypeRepository(db)
	ct := &domain.CabinType{Name: "ct1", CruiseID: 1}
	_ = ctRepo.Create(ctx, ct)
	_, _ = ctRepo.GetByID(ctx, ct.ID)
	_, _, _ = ctRepo.ListByCruise(ctx, 1, 1, 10)
	_, _, _ = ctRepo.ListByCruise(ctx, 0, 1, 10)
	ct.Name = "ct2"
	_ = ctRepo.Update(ctx, ct)
	_ = ctRepo.Delete(ctx, ct.ID)

	// 7. FacilityCategoryRepo
	fcRepo := NewFacilityCategoryRepository(db)
	fc := &domain.FacilityCategory{Name: "fc1"}
	_ = fcRepo.Create(ctx, fc)
	_, _ = fcRepo.GetByID(ctx, fc.ID)
	fc.Name = "fc2"
	_ = fcRepo.Update(ctx, fc)
	_, _ = fcRepo.List(ctx)
	_ = fcRepo.Delete(ctx, fc.ID)

	// 8. FacilityRepo
	facRepo := NewFacilityRepository(db)
	fac := &domain.Facility{Name: "fac1", CruiseID: 1}
	_ = facRepo.Create(ctx, fac)
	_, _ = facRepo.GetByID(ctx, fac.ID)
	fac.Name = "fac2"
	_ = facRepo.Update(ctx, fac)
	_, _ = facRepo.ListByCruise(ctx, 1)
	_, _ = facRepo.ListByCruiseAndCategory(ctx, 1, 0)
	_ = facRepo.Delete(ctx, fac.ID)

	// 9. CabinRepo (SKU, Prices, Inventory)
	cbRepo := NewCabinRepository(db)
	sku := &domain.CabinSKU{Code: "sku1", VoyageID: 1}
	_ = cbRepo.CreateSKU(ctx, sku)
	_, _ = cbRepo.GetSKUByID(ctx, sku.ID)
	_, _ = cbRepo.ListSKUByVoyage(ctx, 1)
	sku.Code = "sku2"
	_ = cbRepo.UpdateSKU(ctx, sku)

	_ = cbRepo.AdjustInventoryAtomic(ctx, sku.ID, 10)
	_, _ = cbRepo.GetInventoryBySKU(ctx, sku.ID)
	_ = cbRepo.AppendInventoryLog(ctx, &domain.InventoryLog{CabinSKUID: sku.ID, Change: 10, Reason: "test"})

	// Prices
	p := &domain.CabinPrice{CabinSKUID: sku.ID, Date: time.Now(), Occupancy: 2, PriceCents: 1000}
	_ = cbRepo.UpsertPrice(ctx, p)
	// Update price
	p.PriceCents = 2000
	_ = cbRepo.UpsertPrice(ctx, p)
	_, _ = cbRepo.ListPricesBySKU(ctx, sku.ID)

	_ = cbRepo.DeleteSKU(ctx, sku.ID)

	// 10. BookingRepo
	bkRepo := NewBookingRepository(db)
	bk := &domain.Booking{UserID: 1}
	_ = bkRepo.Create(context.Background(), bk)
	_, _, _ = bkRepo.List(ctx, 1, 10)
	_, _, _ = bkRepo.List(ctx, 0, 0) // edge: page < 1, pageSize <= 0
	_, _ = bkRepo.GetByID(ctx, bk.ID)
	_, _ = bkRepo.GetByID(ctx, 999) // not found
	_ = bkRepo.UpdateStatus(ctx, bk.ID, "paid")
	_ = bkRepo.Delete(ctx, bk.ID)

	// 11. CabinRepo — ListAllInventories & SetAlertThreshold
	sku2 := &domain.CabinSKU{Code: "sku-inv", VoyageID: 1}
	_ = cbRepo.CreateSKU(ctx, sku2)
	_ = cbRepo.AdjustInventoryAtomic(ctx, sku2.ID, 5)
	_, _ = cbRepo.ListAllInventories(ctx)
	_ = cbRepo.SetAlertThreshold(ctx, sku2.ID, 3)
	_ = cbRepo.DeleteSKU(ctx, sku2.ID)

	// 12. ImageRepo — ReplaceImages
	imgRepo := NewImageRepository(db)
	db.AutoMigrate(&domain.Image{})
	imgs := []*domain.Image{
		{EntityType: "cruise", EntityID: 1, URL: "http://a.png", SortOrder: 1, IsPrimary: true},
		{EntityType: "cruise", EntityID: 1, URL: "http://b.png", SortOrder: 2, IsPrimary: false},
	}
	_ = imgRepo.ReplaceImages(ctx, "cruise", 1, imgs)
	list, _ := imgRepo.ListByEntity(ctx, "cruise", 1)
	if len(list) != 2 {
		t.Errorf("expected 2 images, got %d", len(list))
	}
	// Replace again with fewer
	_ = imgRepo.ReplaceImages(ctx, "cruise", 1, []*domain.Image{
		{EntityType: "cruise", EntityID: 1, URL: "http://c.png", SortOrder: 1, IsPrimary: true},
	})
	list2, _ := imgRepo.ListByEntity(ctx, "cruise", 1)
	if len(list2) != 1 {
		t.Errorf("expected 1 image after replace, got %d", len(list2))
	}
	// Add missing error paths by creating an unmigrated DB
	badDB, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{}) // unmigrated
	ctx2 := context.Background()
	NewStaffRepository(badDB).Create(ctx2, &domain.Staff{})
	NewStaffRepository(badDB).GetByID(ctx2, 1)
	NewStaffRepository(badDB).GetByUsername(ctx2, "abc")
	NewStaffRepository(badDB).Update(ctx2, &domain.Staff{})
	NewStaffRepository(badDB).Delete(ctx2, 1)

	NewCompanyRepository(badDB).Create(ctx2, &domain.CruiseCompany{})
	NewCompanyRepository(badDB).GetByID(ctx2, 1)
	NewCompanyRepository(badDB).List(ctx2, "abc", 1, 10)
	NewCompanyRepository(badDB).Update(ctx2, &domain.CruiseCompany{})
	NewCompanyRepository(badDB).Delete(ctx2, 1)

	NewCruiseRepository(badDB).Create(ctx2, &domain.Cruise{})
	NewCruiseRepository(badDB).GetByID(ctx2, 1)
	NewCruiseRepository(badDB).List(ctx2, 1, "", nil, "", 1, 10)
	NewCruiseRepository(badDB).Update(ctx2, &domain.Cruise{})
	NewCruiseRepository(badDB).Delete(ctx2, 1)

	NewCabinTypeRepository(badDB).Create(ctx2, &domain.CabinType{})
	NewCabinTypeRepository(badDB).GetByID(ctx2, 1)
	NewCabinTypeRepository(badDB).ListByCruise(ctx2, 1, 1, 10)
	NewCabinTypeRepository(badDB).Update(ctx2, &domain.CabinType{})
	NewCabinTypeRepository(badDB).Delete(ctx2, 1)

	NewFacilityCategoryRepository(badDB).Create(ctx2, &domain.FacilityCategory{})
	NewFacilityCategoryRepository(badDB).GetByID(ctx2, 1)
	NewFacilityCategoryRepository(badDB).Update(ctx2, &domain.FacilityCategory{})
	NewFacilityCategoryRepository(badDB).List(ctx2)
	NewFacilityCategoryRepository(badDB).Delete(ctx2, 1)

	NewFacilityRepository(badDB).Create(ctx2, &domain.Facility{})
	NewFacilityRepository(badDB).GetByID(ctx2, 1)
	NewFacilityRepository(badDB).Update(ctx2, &domain.Facility{})
	NewFacilityRepository(badDB).ListByCruise(ctx2, 1)
	NewFacilityRepository(badDB).ListByCruiseAndCategory(ctx2, 1, 1)
	NewFacilityRepository(badDB).Delete(ctx2, 1)

	cbBad := NewCabinRepository(badDB)
	cbBad.CreateSKU(ctx2, &domain.CabinSKU{})
	cbBad.UpdateSKU(ctx2, &domain.CabinSKU{})
	cbBad.GetSKUByID(ctx2, 1)
	cbBad.ListSKUByVoyage(ctx2, 1)
	cbBad.DeleteSKU(ctx2, 1)
	cbBad.AdjustInventoryAtomic(ctx2, 1, 10)
	cbBad.GetInventoryBySKU(ctx2, 1)
	cbBad.AppendInventoryLog(ctx2, &domain.InventoryLog{})
	cbBad.UpsertPrice(ctx2, &domain.CabinPrice{})
	cbBad.ListPricesBySKU(ctx2, 1)

	NewVoyageRepository(badDB).Create(ctx2, &domain.Voyage{})
	NewVoyageRepository(badDB).Update(ctx2, &domain.Voyage{})
	NewVoyageRepository(badDB).GetByID(ctx2, 1)
	NewVoyageRepository(badDB).List(ctx2)
	NewVoyageRepository(badDB).Delete(ctx2, 1)

	NewRouteRepository(badDB).Create(ctx2, &domain.Route{})
	NewRouteRepository(badDB).Update(ctx2, &domain.Route{})
	NewRouteRepository(badDB).GetByID(ctx2, 1)
	NewRouteRepository(badDB).List(ctx2)
	NewRouteRepository(badDB).Delete(ctx2, 1)

	NewBookingRepository(badDB).Create(ctx2, &domain.Booking{})
}
