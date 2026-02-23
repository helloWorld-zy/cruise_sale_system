package repository

import (
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
