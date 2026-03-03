package repository

import (
	"context"
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
	err := repo.Create(context.Background(), &domain.Booking{UserID: 1, VoyageID: 2, CabinSKUID: 3, Status: "created", TotalCents: 100})
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

func TestBookingRepoUpdateStatusWritesLog(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&domain.Booking{}, &domain.OrderStatusLog{})
	repo := NewBookingRepository(db)

	seed := &domain.Booking{UserID: 1, VoyageID: 2, CabinSKUID: 3, Status: domain.OrderStatusCreated, TotalCents: 100}
	if err := repo.Create(context.Background(), seed); err != nil {
		t.Fatal(err)
	}

	if err := repo.TransitionStatus(context.Background(), seed.ID, domain.OrderStatusPendingPayment, 7, "admin update"); err != nil {
		t.Fatalf("expected transition success, got %v", err)
	}

	got, err := repo.GetByID(context.Background(), seed.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.OrderStatusPendingPayment {
		t.Fatalf("expected pending_payment, got %s", got.Status)
	}

	var logs []domain.OrderStatusLog
	if err := db.Where("order_id = ?", seed.ID).Find(&logs).Error; err != nil {
		t.Fatal(err)
	}
	if len(logs) != 1 {
		t.Fatalf("expected 1 status log, got %d", len(logs))
	}
	if logs[0].FromStatus != domain.OrderStatusCreated || logs[0].ToStatus != domain.OrderStatusPendingPayment {
		t.Fatalf("unexpected status log: %+v", logs[0])
	}
	if logs[0].OperatorID != 7 {
		t.Fatalf("expected operator 7, got %d", logs[0].OperatorID)
	}
}

func TestBookingRepoUpdateStatusRejectsInvalidTransition(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&domain.Booking{}, &domain.OrderStatusLog{})
	repo := NewBookingRepository(db)

	seed := &domain.Booking{UserID: 1, VoyageID: 2, CabinSKUID: 3, Status: domain.OrderStatusCreated, TotalCents: 100}
	if err := repo.Create(context.Background(), seed); err != nil {
		t.Fatal(err)
	}

	err := repo.TransitionStatus(context.Background(), seed.ID, domain.OrderStatusPaid, 7, "skip flow")
	if !errors.Is(err, ErrInvalidOrderStatusTransition) {
		t.Fatalf("expected ErrInvalidOrderStatusTransition, got %v", err)
	}

	got, err := repo.GetByID(context.Background(), seed.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.OrderStatusCreated {
		t.Fatalf("expected status unchanged, got %s", got.Status)
	}

	var count int64
	if err := db.Model(&domain.OrderStatusLog{}).Where("order_id = ?", seed.ID).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expected no status log for invalid transition, got %d", count)
	}
}
