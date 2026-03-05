package repository

import (
	"context"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestCompanyRepository_List(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("db error: %v", err)
	}
	if err := db.AutoMigrate(&domain.CruiseCompany{}); err != nil {
		t.Fatalf("migrate error: %v", err)
	}
	repo := NewCompanyRepository(db)
	err = repo.Create(context.Background(), &domain.CruiseCompany{Name: "A", SortOrder: 1})
	if err != nil {
		t.Fatalf("create error: %v", err)
	}

	items, total, err := repo.List(context.Background(), "A", 1, 10)
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if total == 0 || len(items) == 0 {
		t.Fatal("expected items")
	}
}

func TestCompanyRepository_List_SearchEnglishName(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("db error: %v", err)
	}
	if err := db.AutoMigrate(&domain.CruiseCompany{}); err != nil {
		t.Fatalf("migrate error: %v", err)
	}
	repo := NewCompanyRepository(db)
	err = repo.Create(context.Background(), &domain.CruiseCompany{Name: "皇家加勒比", EnglishName: "Royal Caribbean", SortOrder: 1})
	if err != nil {
		t.Fatalf("create error: %v", err)
	}

	items, total, err := repo.List(context.Background(), "Royal", 1, 10)
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("expected 1 item by english name search, got total=%d len=%d", total, len(items))
	}
}
