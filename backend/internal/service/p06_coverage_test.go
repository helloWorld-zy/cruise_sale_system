package service

import (
	"context"
	"errors"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

// --- mock for CabinAdminService.GetByID ---
type mockCabinAdminRepo struct {
	allMockCabinRepo // embed to inherit methods
}

func (m *mockCabinAdminRepo) GetSKUByID(_ context.Context, id int64) (*domain.CabinSKU, error) {
	if id == 99 {
		return nil, errors.New("not found")
	}
	return &domain.CabinSKU{ID: id}, nil
}

type allMockCabinRepo struct{}

func (m *allMockCabinRepo) CreateSKU(_ context.Context, s *domain.CabinSKU) error { return nil }
func (m *allMockCabinRepo) UpdateSKU(_ context.Context, s *domain.CabinSKU) error { return nil }
func (m *allMockCabinRepo) DeleteSKU(_ context.Context, _ int64) error            { return nil }
func (m *allMockCabinRepo) GetSKUByID(_ context.Context, id int64) (*domain.CabinSKU, error) {
	if id == 99 {
		return nil, errors.New("not found")
	}
	return &domain.CabinSKU{ID: id}, nil
}
func (m *allMockCabinRepo) ListSKUByVoyage(_ context.Context, _ int64) ([]domain.CabinSKU, error) {
	return nil, nil
}
func (m *allMockCabinRepo) ListSKUFiltered(_ context.Context, _ domain.CabinSKUFilter) ([]domain.CabinSKU, int64, error) {
	return nil, 0, nil
}
func (m *allMockCabinRepo) BatchUpdateStatus(_ context.Context, _ []int64, _ int16) error {
	return nil
}
func (m *allMockCabinRepo) GetInventoryBySKU(_ context.Context, _ int64) (domain.CabinInventory, error) {
	return domain.CabinInventory{}, nil
}
func (m *allMockCabinRepo) ListAllInventories(_ context.Context) ([]domain.CabinInventory, error) {
	return nil, nil
}
func (m *allMockCabinRepo) SetAlertThreshold(_ context.Context, _ int64, _ int) error { return nil }
func (m *allMockCabinRepo) AdjustInventoryAtomic(_ context.Context, _ int64, _ int) error {
	return nil
}
func (m *allMockCabinRepo) AppendInventoryLog(_ context.Context, _ *domain.InventoryLog) error {
	return nil
}
func (m *allMockCabinRepo) ListPricesBySKU(_ context.Context, _ int64) ([]domain.CabinPrice, error) {
	return nil, nil
}
func (m *allMockCabinRepo) UpsertPrice(_ context.Context, _ *domain.CabinPrice) error { return nil }
func (m *allMockCabinRepo) GetCategoryTree(_ context.Context) (interface{}, error)    { return nil, nil }

// === CabinAdminService.GetByID ===
func TestCabinAdminService_GetByID(t *testing.T) {
	svc := NewCabinAdminService(&allMockCabinRepo{})
	ctx := context.Background()

	item, err := svc.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), item.ID)

	_, err = svc.GetByID(ctx, 99)
	assert.Error(t, err)
}

// === CruiseService.ListWithFilters ===
func TestCruiseService_ListWithFilters(t *testing.T) {
	svc := NewCruiseService(&allMockCruiseRepo{}, &allMockCabinTypeRepo{}, &allMockCompanyRepo{})
	ctx := context.Background()

	items, total, err := svc.ListWithFilters(ctx, 1, "", nil, "", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Nil(t, items)

	_, _, err = svc.ListWithFilters(ctx, 99, "", nil, "", 1, 10)
	assert.Error(t, err)
}

// === FacilityCategoryService.Update / GetByID ===
func TestFacilityCategoryService_UpdateAndGetByID(t *testing.T) {
	svc := NewFacilityCategoryService(&allMockFacilityCategoryRepo{})
	ctx := context.Background()

	err := svc.Update(ctx, &domain.FacilityCategory{Name: "test"})
	assert.NoError(t, err)

	err = svc.Update(ctx, &domain.FacilityCategory{Name: "error"})
	assert.Error(t, err)

	item, err := svc.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), item.ID)

	_, err = svc.GetByID(ctx, 99)
	assert.Error(t, err)
}

// === FacilityService.Update / GetByID / ListByCruiseAndCategory ===
func TestFacilityService_UpdateGetByIDAndListByCruiseAndCategory(t *testing.T) {
	svc := NewFacilityService(&allMockFacilityRepo{})
	ctx := context.Background()

	err := svc.Update(ctx, &domain.Facility{Name: "test"})
	assert.NoError(t, err)

	err = svc.Update(ctx, &domain.Facility{Name: "error"})
	assert.Error(t, err)

	item, err := svc.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), item.ID)

	_, err = svc.GetByID(ctx, 99)
	assert.Error(t, err)

	items, err := svc.ListByCruiseAndCategory(ctx, 1, 1)
	assert.NoError(t, err)
	assert.Nil(t, items)

	_, err = svc.ListByCruiseAndCategory(ctx, 99, 1)
	assert.Error(t, err)
}

// === CruiseService.BatchUpdateStatus paths ===
func TestCruiseService_BatchUpdateStatusP06(t *testing.T) {
	svc := NewCruiseService(&allMockCruiseRepo{}, &allMockCabinTypeRepo{}, &allMockCompanyRepo{})
	ctx := context.Background()

	// success path (id=1 exists)
	err := svc.BatchUpdateStatus(ctx, []int64{1}, 1)
	assert.NoError(t, err)

	// GetByID error (id=99)
	err = svc.BatchUpdateStatus(ctx, []int64{99}, 1)
	assert.Error(t, err)

	// empty, no-op
	err = svc.BatchUpdateStatus(ctx, []int64{}, 1)
	assert.NoError(t, err)
}
