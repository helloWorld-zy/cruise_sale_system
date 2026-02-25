package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openCabinRepoTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	if err := db.AutoMigrate(
		&domain.Route{},
		&domain.Voyage{},
		&domain.CabinSKU{},
		&domain.CabinInventory{},
		&domain.InventoryLog{},
		&domain.CabinPrice{},
		&domain.Booking{},
	); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	return db
}

func seedVoyage(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	route := domain.Route{Code: "R-COV", Name: "Route Coverage"}
	if err := db.Create(&route).Error; err != nil {
		t.Fatalf("create route failed: %v", err)
	}
	voyage := domain.Voyage{Code: "V-COV", RouteID: route.ID, CruiseID: 1}
	if err := db.Create(&voyage).Error; err != nil {
		t.Fatalf("create voyage failed: %v", err)
	}
	return voyage.ID
}

func TestCabinRepositoryAliasMethodsAndInventoryBranches(t *testing.T) {
	db := openCabinRepoTestDB(t)
	repo := NewCabinRepository(db)
	ctx := context.Background()
	voyageID := seedVoyage(t, db)

	sku := &domain.CabinSKU{Code: "SKU-COV-1", VoyageID: voyageID, CabinTypeID: 1, MaxGuests: 2}
	if err := repo.CreateSKU(ctx, sku); err != nil {
		t.Fatalf("CreateSKU failed: %v", err)
	}

	if err := db.Create(&domain.CabinInventory{CabinSKUID: sku.ID, Total: 3, Locked: 0, Sold: 0}).Error; err != nil {
		t.Fatalf("seed inventory failed: %v", err)
	}

	items, err := repo.ListSKUByVoyage(ctx, voyageID)
	if err != nil || len(items) == 0 {
		t.Fatalf("ListSKUByVoyage failed, err=%v len=%d", err, len(items))
	}

	sku.Deck = "D1"
	if err := repo.UpdateSKU(ctx, sku); err != nil {
		t.Fatalf("UpdateSKU failed: %v", err)
	}

	inv, err := repo.GetInventoryBySKU(ctx, sku.ID)
	if err != nil {
		t.Fatalf("GetInventoryBySKU failed: %v", err)
	}
	if inv.Total != 3 {
		t.Fatalf("expected total=3 got=%d", inv.Total)
	}

	if err := repo.AdjustInventoryAtomic(ctx, sku.ID, -1); err != nil {
		t.Fatalf("AdjustInventoryAtomic failed: %v", err)
	}
	if err := repo.AppendInventoryLog(ctx, &domain.InventoryLog{CabinSKUID: sku.ID, Change: -1, Reason: "sell-one"}); err != nil {
		t.Fatalf("AppendInventoryLog failed: %v", err)
	}

	inv, err = repo.GetInventoryBySKU(ctx, sku.ID)
	if err != nil {
		t.Fatalf("GetInventoryBySKU after adjust failed: %v", err)
	}
	if inv.Total != 2 {
		t.Fatalf("expected total=2 got=%d", inv.Total)
	}

	if err := repo.AdjustInventoryAtomic(ctx, sku.ID, -999); err == nil || !errors.Is(err, domain.ErrInsufficientInventory) {
		t.Fatalf("expected ErrInsufficientInventory, got %v", err)
	}

	price := &domain.CabinPrice{CabinSKUID: sku.ID, Date: time.Now(), Occupancy: 2, PriceCents: 12345}
	if err := repo.UpsertPrice(ctx, price); err != nil {
		t.Fatalf("UpsertPrice failed: %v", err)
	}
	prices, err := repo.ListPricesBySKU(ctx, sku.ID)
	if err != nil || len(prices) == 0 {
		t.Fatalf("ListPricesBySKU failed, err=%v len=%d", err, len(prices))
	}
	prices2, err := repo.ListBySKU(ctx, sku.ID)
	if err != nil || len(prices2) == 0 {
		t.Fatalf("ListBySKU failed, err=%v len=%d", err, len(prices2))
	}

	if err := repo.DeleteSKU(ctx, sku.ID); err != nil {
		t.Fatalf("DeleteSKU failed: %v", err)
	}
}

func TestCabinRepositoryAdjustInventoryAtomic_DBError(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}

	repo := NewCabinRepository(db)
	if err := repo.AdjustInventoryAtomic(context.Background(), 1, -1); err == nil {
		t.Fatal("expected db error when cabin_inventories table does not exist")
	}
}
