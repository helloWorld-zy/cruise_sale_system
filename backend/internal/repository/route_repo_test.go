package repository

import (
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRouteRepoCreateAndList(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	_ = db.AutoMigrate(&domain.Route{})
	repo := NewRouteRepository(db)
	if err := repo.Create(&domain.Route{Code: "R1", Name: "Route 1"}); err != nil {
		t.Fatal(err)
	}
	list, err := repo.List()
	if err != nil || len(list) != 1 {
		t.Fatal("expected 1 route")
	}
}
