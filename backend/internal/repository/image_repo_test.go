package repository

import (
	"context"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestImageRepository_CRUD(t *testing.T) {
	// 使用内存数据库验证图片仓储的基础 CRUD 与排序更新行为。
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	if err := db.AutoMigrate(&domain.Image{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	repo := NewImageRepository(db)
	ctx := context.Background()

	img := &domain.Image{EntityType: "cruise", EntityID: 1, URL: "https://img.com/1.jpg", IsPrimary: true}
	if err := repo.Create(ctx, img); err != nil {
		t.Fatal(err)
	}

	list, err := repo.ListByEntity(ctx, "cruise", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1, got %d", len(list))
	}
	if !list[0].IsPrimary {
		t.Fatal("expected primary")
	}

	// 验证排序更新能力，避免接口方法无测试覆盖。
	if err := repo.UpdateSortOrder(ctx, list[0].ID, 5); err != nil {
		t.Fatal(err)
	}
	updated, err := repo.ListByEntity(ctx, "cruise", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(updated) != 1 || updated[0].SortOrder != 5 {
		t.Fatalf("expected sort_order 5, got %d", updated[0].SortOrder)
	}

	if err := repo.DeleteByEntity(ctx, "cruise", 1); err != nil {
		t.Fatal(err)
	}
	list2, err := repo.ListByEntity(ctx, "cruise", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(list2) != 0 {
		t.Fatal("expected 0 after delete")
	}
}
