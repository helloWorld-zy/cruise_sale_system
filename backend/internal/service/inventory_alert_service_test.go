package service

import (
	"context"
	"errors"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
)

type inventoryAlertRepoMock struct {
	items         []domain.CabinInventory
	listErr       error
	setErr        error
	lastSKU       int64
	lastThreshold int
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
