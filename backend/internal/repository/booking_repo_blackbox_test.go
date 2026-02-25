package repository_test

import (
	"errors"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestBookingRepositoryInTxBlackbox(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	if err := db.AutoMigrate(&domain.Booking{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}

	repo := repository.NewBookingRepository(db)
	if err := repo.InTx(func(tx *gorm.DB, create func(b *domain.Booking) error) error {
		return create(&domain.Booking{UserID: 1, VoyageID: 2, CabinSKUID: 3, Status: "created", TotalCents: 100})
	}); err != nil {
		t.Fatalf("InTx create failed: %v", err)
	}

	if err := repo.InTx(func(tx *gorm.DB, create func(b *domain.Booking) error) error {
		_ = tx
		_ = create
		return errors.New("force rollback")
	}); err == nil {
		t.Fatal("expected rollback error")
	}
}
