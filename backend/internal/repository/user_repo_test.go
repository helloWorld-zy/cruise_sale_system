package repository

import (
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserRepositoryFindOrCreateByPhone_Create(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		t.Fatal(err)
	}

	repo := NewUserRepository(db)
	u, err := repo.FindOrCreateByPhone("13800000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.ID == 0 {
		t.Fatal("expected user to be persisted with non-zero ID")
	}
	if u.Phone != "13800000001" {
		t.Fatalf("expected phone 13800000001, got %s", u.Phone)
	}
}

func TestUserRepositoryFindOrCreateByPhone_Idempotent(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		t.Fatal(err)
	}

	repo := NewUserRepository(db)
	u1, err := repo.FindOrCreateByPhone("13900000001")
	if err != nil {
		t.Fatal(err)
	}
	u2, err := repo.FindOrCreateByPhone("13900000001")
	if err != nil {
		t.Fatal(err)
	}
	if u1.ID != u2.ID {
		t.Fatalf("expected same user ID on second call, got %d vs %d", u1.ID, u2.ID)
	}
}
