package repository

import (
	"context"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	return db
}

func TestRouteRepoCreateAndList(t *testing.T) {
	db := openTestDB(t)
	_ = db.AutoMigrate(&domain.Route{})
	repo := NewRouteRepository(db)
	ctx := context.Background()
	if err := repo.Create(ctx, &domain.Route{Code: "R1", Name: "Route 1"}); err != nil {
		t.Fatal(err)
	}
	list, err := repo.List(ctx)
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 route, got %d (err=%v)", len(list), err)
	}
}

func TestRouteRepoGetByID(t *testing.T) {
	db := openTestDB(t)
	_ = db.AutoMigrate(&domain.Route{})
	repo := NewRouteRepository(db)
	ctx := context.Background()
	_ = repo.Create(ctx, &domain.Route{Code: "R2", Name: "Route 2"})
	list, _ := repo.List(ctx)
	got, err := repo.GetByID(ctx, list[0].ID)
	if err != nil || got.Code != "R2" {
		t.Fatalf("GetByID failed: %v", err)
	}
}

func TestRouteRepoUpdateAndDelete(t *testing.T) {
	db := openTestDB(t)
	_ = db.AutoMigrate(&domain.Route{})
	repo := NewRouteRepository(db)
	ctx := context.Background()
	_ = repo.Create(ctx, &domain.Route{Code: "R3", Name: "Route 3"})
	list, _ := repo.List(ctx)
	r := list[0]
	r.Name = "Updated"
	if err := repo.Update(ctx, &r); err != nil {
		t.Fatal(err)
	}
	if err := repo.Delete(ctx, r.ID); err != nil {
		t.Fatal(err)
	}
	list2, _ := repo.List(ctx)
	if len(list2) != 0 {
		t.Fatal("expected empty list after delete")
	}
}
