package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type fakeBookingRepo struct{ created bool }

func (f *fakeBookingRepo) Create(_ *domain.Booking) error { f.created = true; return nil }
func (f *fakeBookingRepo) InTx(fn func(tx *gorm.DB, create func(b *domain.Booking) error) error) error {
	return fn(nil, f.Create)
}

type fakeBookingRepoTxErr struct{}

func (f *fakeBookingRepoTxErr) Create(_ *domain.Booking) error { return nil }
func (f *fakeBookingRepoTxErr) InTx(fn func(tx *gorm.DB, create func(b *domain.Booking) error) error) error {
	_ = fn
	return errors.New("tx failed")
}

type fakePriceService struct{}

func (f fakePriceService) FindPrice(_ context.Context, _ int64, _ time.Time, _ int) (int64, bool, error) {
	return 10000, true, nil
}

type fakeHoldService struct{ ok bool }

func (f *fakeHoldService) HoldWithTx(_ *gorm.DB, _ int64, _ int64, _ int) bool {
	f.ok = true
	return true
}

func TestBookingServiceCreate(t *testing.T) {
	svc := NewBookingService(&fakeBookingRepo{}, fakePriceService{}, &fakeHoldService{})
	if err := svc.Create(1, 2, 3, 2); err != nil {
		t.Fatal(err)
	}
}

func TestBookingServiceCreate_TxFail(t *testing.T) {
	svc := NewBookingService(&fakeBookingRepoTxErr{}, fakePriceService{}, &fakeHoldService{})
	if err := svc.Create(1, 2, 3, 2); err == nil {
		t.Fatal("expected tx failure")
	}
}

func TestBookingServiceCreate_DependencyNotReady(t *testing.T) {
	svc := NewBookingService(nil, nil, nil)
	if err := svc.Create(1, 2, 3, 2); err == nil {
		t.Fatal("expected dependency error")
	}
}

type txFailBookingRepo struct{ db *gorm.DB }

func (r *txFailBookingRepo) Create(_ *domain.Booking) error { return nil }
func (r *txFailBookingRepo) InTx(fn func(tx *gorm.DB, create func(b *domain.Booking) error) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		create := func(b *domain.Booking) error {
			_ = b
			return errors.New("forced create failure")
		}
		return fn(tx, create)
	})
}

func TestBookingServiceCreate_RollbackHoldAndInventoryOnCreateFail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	if err := db.AutoMigrate(&domain.CabinInventory{}, &domain.InventoryLog{}, &domain.CabinHold{}, &domain.Booking{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&domain.CabinInventory{CabinSKUID: 3, Total: 1, Locked: 0, Sold: 0}).Error; err != nil {
		t.Fatal(err)
	}

	holdSvc := NewCabinHoldService(repository.NewCabinHoldRepository(db), time.Minute)
	svc := NewBookingService(&txFailBookingRepo{db: db}, fakePriceService{}, holdSvc)

	if err := svc.Create(1, 2, 3, 2); err == nil {
		t.Fatal("expected create to fail")
	}

	var inv domain.CabinInventory
	if err := db.Where("cabin_sku_id = ?", 3).First(&inv).Error; err != nil {
		t.Fatal(err)
	}
	if inv.Total != 1 {
		t.Fatalf("expected inventory rollback to total=1, got %d", inv.Total)
	}

	var holdCount int64
	if err := db.Model(&domain.CabinHold{}).Where("cabin_sku_id = ? AND user_id = ?", 3, 1).Count(&holdCount).Error; err != nil {
		t.Fatal(err)
	}
	if holdCount != 0 {
		t.Fatalf("expected hold rollback with zero records, got %d", holdCount)
	}
}
