package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

type inventoryAlertRepoMock struct {
	items         []domain.CabinInventory
	listErr       error
	getErr        error
	setErr        error
	lastSKU       int64
	lastThreshold int
}

type inventoryAlertNotifierMock struct {
	calls int
	last  struct {
		userID   int64
		channel  string
		template string
		payload  string
	}
}

func (m *inventoryAlertNotifierMock) Enqueue(ctx context.Context, userID int64, channel, template, payload string) error {
	_ = ctx
	m.calls++
	m.last.userID = userID
	m.last.channel = channel
	m.last.template = template
	m.last.payload = payload
	return nil
}

func (m *inventoryAlertRepoMock) ListAllInventories(ctx context.Context) ([]domain.CabinInventory, error) {
	_ = ctx
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.items, nil
}

func (m *inventoryAlertRepoMock) SetAlertThreshold(ctx context.Context, skuID int64, threshold int) error {
	_ = ctx
	m.lastSKU = skuID
	m.lastThreshold = threshold
	return m.setErr
}

func (m *inventoryAlertRepoMock) GetInventoryBySKU(ctx context.Context, skuID int64) (domain.CabinInventory, error) {
	_ = ctx
	if m.getErr != nil {
		return domain.CabinInventory{}, m.getErr
	}
	for _, item := range m.items {
		if item.CabinSKUID == skuID {
			return item, nil
		}
	}
	return domain.CabinInventory{}, nil
}

func TestInventoryAlertService_CheckAlerts(t *testing.T) {
	repo := &inventoryAlertRepoMock{items: []domain.CabinInventory{
		{CabinSKUID: 1, Total: 10, Locked: 2, Sold: 6, AlertThreshold: 3}, // available=2, 命中
		{CabinSKUID: 2, Total: 10, Locked: 1, Sold: 2, AlertThreshold: 3}, // available=7, 不命中
		{CabinSKUID: 3, Total: 8, Locked: 1, Sold: 6, AlertThreshold: 0},  // threshold=0, 忽略
	}}
	svc := NewInventoryAlertService(repo)

	alerts, err := svc.CheckAlerts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].CabinSKUID != 1 || alerts[0].Available != 2 {
		t.Fatalf("unexpected alert: %+v", alerts[0])
	}
}

func TestInventoryAlertService_SetThreshold(t *testing.T) {
	repo := &inventoryAlertRepoMock{}
	svc := NewInventoryAlertService(repo)

	if err := svc.SetThreshold(context.Background(), 9, 5); err != nil {
		t.Fatal(err)
	}
	if repo.lastSKU != 9 || repo.lastThreshold != 5 {
		t.Fatalf("unexpected set call: sku=%d threshold=%d", repo.lastSKU, repo.lastThreshold)
	}
}

func TestInventoryAlertService_PropagateRepoError(t *testing.T) {
	repo := &inventoryAlertRepoMock{listErr: errors.New("db error")}
	svc := NewInventoryAlertService(repo)

	if _, err := svc.CheckAlerts(context.Background()); err == nil {
		t.Fatal("expected list error")
	}

	repo2 := &inventoryAlertRepoMock{setErr: errors.New("set error")}
	svc2 := NewInventoryAlertService(repo2)
	if err := svc2.SetThreshold(context.Background(), 1, 2); err == nil {
		t.Fatal("expected set error")
	}
}

func TestInventoryAlertService_AvailableWithAlertBoundaries(t *testing.T) {
	repo := &inventoryAlertRepoMock{items: []domain.CabinInventory{
		{CabinSKUID: 11, Total: 10, Locked: 2, Sold: 5, AlertThreshold: 3}, // available=3, == threshold, 命中
		{CabinSKUID: 12, Total: 10, Locked: 2, Sold: 6, AlertThreshold: 3}, // available=2, < threshold, 命中
		{CabinSKUID: 13, Total: 10, Locked: 2, Sold: 8, AlertThreshold: 0}, // threshold=0, 忽略
	}}
	svc := NewInventoryAlertService(repo)

	available, alert, err := svc.AvailableWithAlert(context.Background(), 11)
	if err != nil || available != 3 || !alert {
		t.Fatalf("expected available=3 alert=true, got available=%d alert=%v err=%v", available, alert, err)
	}

	available, alert, err = svc.AvailableWithAlert(context.Background(), 12)
	if err != nil || available != 2 || !alert {
		t.Fatalf("expected available=2 alert=true, got available=%d alert=%v err=%v", available, alert, err)
	}

	available, alert, err = svc.AvailableWithAlert(context.Background(), 13)
	if err != nil || available != 0 || alert {
		t.Fatalf("expected available=0 alert=false, got available=%d alert=%v err=%v", available, alert, err)
	}
}

func TestInventoryAlertService_ScanAndNotifyDeduplicatesWithinWindow(t *testing.T) {
	repo := &inventoryAlertRepoMock{items: []domain.CabinInventory{{CabinSKUID: 21, Total: 10, Locked: 2, Sold: 7, AlertThreshold: 2}}}
	notify := &inventoryAlertNotifierMock{}

	now := time.Date(2026, 3, 3, 10, 0, 0, 0, time.UTC)
	svc := NewInventoryAlertServiceWithNotify(repo, notify, time.Hour)
	svc.now = func() time.Time { return now }

	count, err := svc.ScanAndNotify(context.Background(), 1001, ChannelInbox)
	if err != nil {
		t.Fatalf("expected first scan success, got %v", err)
	}
	if count != 1 || notify.calls != 1 {
		t.Fatalf("expected first scan notify once, count=%d calls=%d", count, notify.calls)
	}

	count, err = svc.ScanAndNotify(context.Background(), 1001, ChannelInbox)
	if err != nil {
		t.Fatalf("expected second scan success, got %v", err)
	}
	if count != 0 || notify.calls != 1 {
		t.Fatalf("expected deduped second scan, count=%d calls=%d", count, notify.calls)
	}

	now = now.Add(2 * time.Hour)
	count, err = svc.ScanAndNotify(context.Background(), 1001, ChannelInbox)
	if err != nil {
		t.Fatalf("expected third scan success, got %v", err)
	}
	if count != 1 || notify.calls != 2 {
		t.Fatalf("expected notify after window expires, count=%d calls=%d", count, notify.calls)
	}
}
