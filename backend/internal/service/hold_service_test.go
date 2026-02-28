package service

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

type fakeHoldRepo struct {
	mu             sync.Mutex
	exists         bool
	activeHolds    map[string]bool
	createError    error
	createCalls    int
	adjustErr      error
	adjustCalls    []int
	ExistsActiveFn func(tx *gorm.DB, skuID, userID int64, now time.Time) (bool, error)
}

func (f *fakeHoldRepo) ExistsActiveHoldTx(tx *gorm.DB, skuID, userID int64, now time.Time) (bool, error) {
	if f.ExistsActiveFn != nil {
		return f.ExistsActiveFn(tx, skuID, userID, now)
	}
	_ = tx
	_ = now
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.activeHolds != nil {
		return f.activeHolds[fmt.Sprintf("%d:%d", skuID, userID)], nil
	}
	return f.exists, nil
}

func (f *fakeHoldRepo) CreateHoldTx(tx *gorm.DB, hold *domain.CabinHold) error {
	_ = tx
	f.mu.Lock()
	defer f.mu.Unlock()
	f.createCalls++
	if f.activeHolds == nil {
		f.activeHolds = make(map[string]bool)
	}
	f.activeHolds[fmt.Sprintf("%d:%d", hold.CabinSKUID, hold.UserID)] = true
	return f.createError
}

func (f *fakeHoldRepo) AdjustInventoryTx(tx *gorm.DB, skuID int64, delta int, reason string) error {
	_ = tx
	_ = skuID
	_ = reason
	f.mu.Lock()
	defer f.mu.Unlock()
	f.adjustCalls = append(f.adjustCalls, delta)
	if f.adjustErr != nil {
		return f.adjustErr
	}
	return nil
}

func TestCabinHoldServiceHold_Idempotent(t *testing.T) {
	repo := &fakeHoldRepo{exists: true}
	svc := NewCabinHoldService(repo, time.Minute)
	if !svc.Hold(1, 1, 1) {
		t.Fatal("expected existing hold to be treated as success")
	}
}

func TestCabinHoldServiceHold_AdjustFailure(t *testing.T) {
	repo := &fakeHoldRepo{adjustErr: errors.New("insufficient")}
	svc := NewCabinHoldService(repo, time.Minute)
	if svc.Hold(1, 1, 1) {
		t.Fatal("expected hold to fail when inventory adjust fails")
	}
}

func TestCabinHoldServiceHold_RollbackWhenCreateFails(t *testing.T) {
	repo := &fakeHoldRepo{createError: errors.New("db error")}
	svc := NewCabinHoldService(repo, time.Minute)

	if svc.Hold(1, 1, 1) {
		t.Fatal("expected hold to fail when hold creation fails")
	}

	if len(repo.adjustCalls) != 2 || repo.adjustCalls[0] != -1 || repo.adjustCalls[1] != 1 {
		t.Fatalf("expected adjust and rollback calls, got %+v", repo.adjustCalls)
	}
}

func TestCabinHoldServiceHold_UsesInjectableClock(t *testing.T) {
	fixedTime := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	var capturedNow time.Time

	repo := &fakeHoldRepo{
		ExistsActiveFn: func(_ *gorm.DB, _, _ int64, now time.Time) (bool, error) {
			capturedNow = now
			return false, nil
		},
	}
	svc := NewCabinHoldService(repo, time.Minute)
	svc.now = func() time.Time { return fixedTime }

	_ = svc.Hold(1, 1, 1) // 调用；应当设置 capturedNow

	if !capturedNow.Equal(fixedTime) {
		t.Fatalf("expected now=%v to be passed to repo, got %v", fixedTime, capturedNow)
	}
}

func TestCabinHoldServiceHold_ConcurrentIdempotent(t *testing.T) {
	repo := &fakeHoldRepo{activeHolds: make(map[string]bool)}
	svc := NewCabinHoldService(repo, time.Minute)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !svc.Hold(1, 1, 1) {
				t.Error("expected concurrent hold to succeed")
			}
		}()
	}
	wg.Wait()

	repo.mu.Lock()
	createCalls := repo.createCalls
	repo.mu.Unlock()
	if createCalls != 1 {
		t.Fatalf("expected one hold creation under concurrency, got %d", createCalls)
	}

	repo.mu.Lock()
	adjustCalls := append([]int(nil), repo.adjustCalls...)
	repo.mu.Unlock()
	if len(adjustCalls) != 1 || adjustCalls[0] != -1 {
		t.Fatalf("expected single inventory deduction, got %+v", adjustCalls)
	}
}
