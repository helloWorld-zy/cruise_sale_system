package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

type fakeOrderTimeoutRepo struct {
	mu              sync.Mutex
	orders          map[int64]domain.Booking
	transitionCalls []struct {
		id     int64
		status string
	}
}

func (f *fakeOrderTimeoutRepo) FindExpiredOrders(ctx context.Context, timeout time.Duration) ([]domain.Booking, error) {
	_ = ctx
	_ = timeout
	f.mu.Lock()
	defer f.mu.Unlock()
	list := make([]domain.Booking, 0)
	for _, order := range f.orders {
		if order.Status == domain.OrderStatusPendingPayment {
			list = append(list, order)
		}
	}
	return list, nil
}

func (f *fakeOrderTimeoutRepo) TransitionStatus(ctx context.Context, id int64, status string, operatorID int64, remark string) error {
	_ = ctx
	f.mu.Lock()
	defer f.mu.Unlock()
	order, ok := f.orders[id]
	if !ok {
		return errors.New("order not found")
	}
	order.Status = status
	f.orders[id] = order
	f.transitionCalls = append(f.transitionCalls, struct {
		id     int64
		status string
	}{id: id, status: status})
	return nil
}

type fakeInventoryReleaser struct {
	mu           sync.Mutex
	releaseCalls int
	releaseErr   error
}

func (f *fakeInventoryReleaser) ReleaseLocked(ctx context.Context, skuID int64, quantity int) error {
	_ = ctx
	_ = skuID
	_ = quantity
	f.mu.Lock()
	defer f.mu.Unlock()
	f.releaseCalls++
	return f.releaseErr
}

func TestCloseExpiredOrdersRollbackOnInventoryReleaseFailure(t *testing.T) {
	repo := &fakeOrderTimeoutRepo{orders: map[int64]domain.Booking{
		1: {ID: 1, CabinSKUID: 101, Status: domain.OrderStatusPendingPayment},
	}}
	inv := &fakeInventoryReleaser{releaseErr: errors.New("release failed")}
	svc := NewOrderTimeoutService(repo, inv)

	closed, err := svc.CloseExpiredOrders(context.Background(), 15*time.Minute)
	if err != nil {
		t.Fatalf("expected no top-level error, got %v", err)
	}
	if closed != 0 {
		t.Fatalf("expected closed count 0, got %d", closed)
	}

	repo.mu.Lock()
	defer repo.mu.Unlock()
	if got := repo.orders[1].Status; got != domain.OrderStatusPendingPayment {
		t.Fatalf("expected rollback to pending_payment, got %s", got)
	}
	if len(repo.transitionCalls) != 2 {
		t.Fatalf("expected 2 transition calls (cancel + rollback), got %d", len(repo.transitionCalls))
	}
	if repo.transitionCalls[0].status != domain.OrderStatusCancelled || repo.transitionCalls[1].status != domain.OrderStatusPendingPayment {
		t.Fatalf("unexpected transition sequence: %+v", repo.transitionCalls)
	}
}

func TestCloseExpiredOrdersConcurrentIdempotent(t *testing.T) {
	repo := &fakeOrderTimeoutRepo{orders: map[int64]domain.Booking{
		1: {ID: 1, CabinSKUID: 101, Status: domain.OrderStatusPendingPayment},
	}}
	inv := &fakeInventoryReleaser{}
	svc := NewOrderTimeoutService(repo, inv)

	var wg sync.WaitGroup
	results := make(chan int, 2)
	errCh := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			closed, err := svc.CloseExpiredOrders(context.Background(), 15*time.Minute)
			results <- closed
			errCh <- err
		}()
	}
	wg.Wait()
	close(results)
	close(errCh)

	totalClosed := 0
	for n := range results {
		totalClosed += n
	}
	for err := range errCh {
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	}
	if totalClosed != 1 {
		t.Fatalf("expected total closed count 1 across concurrent calls, got %d", totalClosed)
	}
	inv.mu.Lock()
	defer inv.mu.Unlock()
	if inv.releaseCalls != 1 {
		t.Fatalf("expected inventory release once, got %d", inv.releaseCalls)
	}
}
