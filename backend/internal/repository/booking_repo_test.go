package repository

import (
	"errors"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestBookingRepoCreate(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&domain.Booking{})
	repo := NewBookingRepository(db)
	err := repo.Create(&domain.Booking{UserID: 1, VoyageID: 2, CabinSKUID: 3, Status: "created", TotalCents: 100})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBookingRepoInTx(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&domain.Booking{})
	repo := NewBookingRepository(db)

	err := repo.InTx(func(tx *gorm.DB, create func(b *domain.Booking) error) error {
		if tx == nil {
			t.Fatal("tx should not be nil")
		}
		return create(&domain.Booking{UserID: 1, VoyageID: 2, CabinSKUID: 3, Status: "created", TotalCents: 100})
	})
	if err != nil {
		t.Fatalf("expected InTx success, got %v", err)
	}

	err = repo.InTx(func(tx *gorm.DB, create func(b *domain.Booking) error) error {
		_ = tx
		_ = create
		return errors.New("forced tx error")
	})
	if err == nil {
		t.Fatal("expected InTx failure")
	}
}
