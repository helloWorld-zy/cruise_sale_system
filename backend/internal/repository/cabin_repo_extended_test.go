package repository

import (
	"context"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCabinRepository_ListWithFilters(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	if err := db.AutoMigrate(&domain.CabinSKU{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	repo := NewCabinRepository(db)
	ctx := context.Background()

	_ = repo.CreateSKU(ctx, &domain.CabinSKU{Code: "A101", VoyageID: 1, CabinTypeID: 1, Status: 1})
	_ = repo.CreateSKU(ctx, &domain.CabinSKU{Code: "B201", VoyageID: 1, CabinTypeID: 2, Status: 0})
	_ = repo.CreateSKU(ctx, &domain.CabinSKU{Code: "A102", VoyageID: 2, CabinTypeID: 1, Status: 1})

	filter := domain.CabinSKUFilter{VoyageID: 1, CabinTypeID: 1, Page: 1, PageSize: 10}
	items, total, err := repo.ListSKUFiltered(ctx, filter)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 {
		t.Fatalf("expected 1, got %d", total)
	}
	if len(items) != 1 || items[0].Code != "A101" {
		t.Fatalf("unexpected items: %+v", items)
	}
}

func TestCabinRepository_BatchUpdateStatus(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	if err := db.AutoMigrate(&domain.CabinSKU{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	repo := NewCabinRepository(db)
	ctx := context.Background()

	_ = repo.CreateSKU(ctx, &domain.CabinSKU{Code: "X1", Status: 1})
	_ = repo.CreateSKU(ctx, &domain.CabinSKU{Code: "X2", Status: 1})

	if err := repo.BatchUpdateStatus(ctx, []int64{1, 2}, 0); err != nil {
		t.Fatal(err)
	}
	one, err := repo.GetSKUByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	two, err := repo.GetSKUByID(ctx, 2)
	if err != nil {
		t.Fatal(err)
	}
	if one.Status != 0 || two.Status != 0 {
		t.Fatalf("expected both status=0, got one=%d two=%d", one.Status, two.Status)
	}
}
