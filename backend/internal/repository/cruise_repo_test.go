package repository

import (
	"context"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCruiseRepository_ListWithFilters(t *testing.T) {
	// 使用内存数据库验证邮轮列表筛选能力（关键词/状态）。
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	if err := db.AutoMigrate(&domain.Cruise{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	repo := NewCruiseRepository(db)
	ctx := context.Background()

	if err := repo.Create(ctx, &domain.Cruise{Name: "海洋量子号", Code: "QNTS", Status: 1, Tonnage: 168000}); err != nil {
		t.Fatal(err)
	}
	if err := repo.Create(ctx, &domain.Cruise{Name: "地中海号", Code: "MEDIT", Status: 2, Tonnage: 90000}); err != nil {
		t.Fatal(err)
	}

	items, total, err := repo.List(ctx, 0, "海洋", nil, "", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatal("expected 1 result filtered by name")
	}

	statusOne := int16(1)
	items2, total2, err := repo.List(ctx, 0, "", &statusOne, "", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if total2 != 1 || len(items2) != 1 {
		t.Fatalf("expected 1 active, got %d", total2)
	}
}
