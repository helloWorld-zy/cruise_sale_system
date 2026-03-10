package service

import (
	"context"
	"errors"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
)

// --- Mock 实现 ---

type mockCompanyRepo struct{}

func (m *mockCompanyRepo) Create(ctx context.Context, company *domain.CruiseCompany) error {
	return nil
}
func (m *mockCompanyRepo) Update(ctx context.Context, company *domain.CruiseCompany) error {
	return nil
}
func (m *mockCompanyRepo) GetByID(ctx context.Context, id int64) (*domain.CruiseCompany, error) {
	if id == 99 {
		return nil, errors.New("not found")
	}
	return &domain.CruiseCompany{ID: id}, nil
}
func (m *mockCompanyRepo) List(ctx context.Context, keyword string, page, pageSize int) ([]domain.CruiseCompany, int64, error) {
	return nil, 0, nil
}
func (m *mockCompanyRepo) ListPublic(ctx context.Context, page, pageSize int) ([]domain.CruiseCompany, int64, error) {
	return []domain.CruiseCompany{{ID: 1, Status: 1}}, 1, nil
}
func (m *mockCompanyRepo) Delete(ctx context.Context, id int64) error { return nil }

type mockCruiseRepo struct{ created bool }

func (m *mockCruiseRepo) Create(ctx context.Context, cruise *domain.Cruise) error {
	m.created = true
	return nil
}
func (m *mockCruiseRepo) Update(ctx context.Context, cruise *domain.Cruise) error { return nil }
func (m *mockCruiseRepo) GetByID(ctx context.Context, id int64) (*domain.Cruise, error) {
	return nil, nil
}
func (m *mockCruiseRepo) List(ctx context.Context, companyID int64, keyword string, status *int16, sortBy string, page, pageSize int) ([]domain.Cruise, int64, error) {
	_ = keyword
	_ = status
	_ = sortBy
	if companyID == 55 {
		// 模拟公司 55 有邮轮 → 删除应失败
		return []domain.Cruise{{ID: 1}}, 1, nil
	}
	return nil, 0, nil
}
func (m *mockCruiseRepo) ListPublic(ctx context.Context, companyID int64, keyword string, sortBy string, page, pageSize int) ([]domain.Cruise, int64, error) {
	if companyID > 0 {
		return []domain.Cruise{{ID: 2, CompanyID: companyID, Status: 1}}, 1, nil
	}
	return []domain.Cruise{{ID: 1, CompanyID: 1, Status: 1}}, 1, nil
}
func (m *mockCruiseRepo) Delete(ctx context.Context, id int64) error { return nil }

type mockCruiseRepoFKDelete struct{}

func (m *mockCruiseRepoFKDelete) Create(ctx context.Context, cruise *domain.Cruise) error { return nil }
func (m *mockCruiseRepoFKDelete) Update(ctx context.Context, cruise *domain.Cruise) error { return nil }
func (m *mockCruiseRepoFKDelete) GetByID(ctx context.Context, id int64) (*domain.Cruise, error) {
	return &domain.Cruise{ID: id}, nil
}
func (m *mockCruiseRepoFKDelete) List(ctx context.Context, companyID int64, keyword string, status *int16, sortBy string, page, pageSize int) ([]domain.Cruise, int64, error) {
	return nil, 0, nil
}
func (m *mockCruiseRepoFKDelete) ListPublic(ctx context.Context, companyID int64, keyword string, sortBy string, page, pageSize int) ([]domain.Cruise, int64, error) {
	return nil, 0, nil
}
func (m *mockCruiseRepoFKDelete) Delete(ctx context.Context, id int64) error {
	return errors.New(`ERROR: update or delete on table "cruises" violates foreign key constraint "voyages_cruise_id_fkey" on table "voyages" (SQLSTATE 23503)`)
}

type mockCabinRepo struct{}

func (m *mockCabinRepo) Create(ctx context.Context, cabinType *domain.CabinType) error { return nil }
func (m *mockCabinRepo) Update(ctx context.Context, cabinType *domain.CabinType) error { return nil }
func (m *mockCabinRepo) GetByID(ctx context.Context, id int64) (*domain.CabinType, error) {
	return nil, nil
}
func (m *mockCabinRepo) ListByCruise(ctx context.Context, cruiseID int64, page, pageSize int) ([]domain.CabinType, int64, error) {
	if cruiseID == 1 {
		return []domain.CabinType{{ID: 1}}, 1, nil
	}
	return nil, 0, nil
}
func (m *mockCabinRepo) Delete(ctx context.Context, id int64) error { return nil }

// --- CruiseService 测试 ---

func TestCruiseService_CreateRequiresCompany(t *testing.T) {
	svc := NewCruiseService(&mockCruiseRepo{}, &mockCabinRepo{}, &mockCompanyRepo{})
	err := svc.Create(context.Background(), &domain.Cruise{CompanyID: 99})
	if err == nil {
		t.Fatal("expected error when company not found")
	}
}

func TestCruiseService_CreateSucceeds(t *testing.T) {
	cr := &mockCruiseRepo{}
	svc := NewCruiseService(cr, &mockCabinRepo{}, &mockCompanyRepo{})
	err := svc.Create(context.Background(), &domain.Cruise{CompanyID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cr.created {
		t.Fatal("expected cruise repo Create to be called")
	}
}

func TestCruiseService_DeleteFailsWhenCabinsExist(t *testing.T) {
	svc := NewCruiseService(&mockCruiseRepo{}, &mockCabinRepo{}, &mockCompanyRepo{})
	err := svc.Delete(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error when cabins exist")
	}
	if !errors.Is(err, ErrCruiseHasCabins) {
		t.Fatalf("expected ErrCruiseHasCabins, got %v", err)
	}
}

func TestCruiseService_DeleteSucceedsWhenNoCabins(t *testing.T) {
	svc := NewCruiseService(&mockCruiseRepo{}, &mockCabinRepo{}, &mockCompanyRepo{})
	err := svc.Delete(context.Background(), 99) // cruiseID 99 → 无舱房
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestCruiseService_DeleteFailsWhenVoyagesExist(t *testing.T) {
	svc := NewCruiseService(&mockCruiseRepoFKDelete{}, &mockCabinRepo{}, &mockCompanyRepo{})
	err := svc.Delete(context.Background(), 99)
	if err == nil {
		t.Fatal("expected error when voyages exist")
	}
	if !errors.Is(err, ErrCruiseHasVoyages) {
		t.Fatalf("expected ErrCruiseHasVoyages, got %v", err)
	}
}

func TestCruiseService_UpdateRequiresValidCompany(t *testing.T) {
	svc := NewCruiseService(&mockCruiseRepo{}, &mockCabinRepo{}, &mockCompanyRepo{})
	err := svc.Update(context.Background(), &domain.Cruise{ID: 1, CompanyID: 99, Name: "x"})
	if err == nil {
		t.Fatal("expected error when updating cruise with invalid company")
	}
}

// --- CompanyService 测试 ---

func TestCompanyService_DeleteFailsWhenCruisesExist(t *testing.T) {
	svc := NewCompanyService(&mockCompanyRepo{}, &mockCruiseRepo{})
	// Mock 中公司 55 有邮轮
	err := svc.Delete(context.Background(), 55)
	if err == nil {
		t.Fatal("expected error when company has cruises")
	}
	if !errors.Is(err, ErrCompanyHasCruises) {
		t.Fatalf("expected ErrCompanyHasCruises, got %v", err)
	}
}

func TestCompanyService_DeleteSucceedsWhenNoCruises(t *testing.T) {
	svc := NewCompanyService(&mockCompanyRepo{}, &mockCruiseRepo{})
	// Mock 中公司 1 没有邮轮
	err := svc.Delete(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestCompanyService_ListPublic(t *testing.T) {
	svc := NewCompanyService(&mockCompanyRepo{}, &mockCruiseRepo{})
	items, total, err := svc.ListPublic(context.Background(), 1, 50)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].Status != 1 {
		t.Fatalf("expected one enabled public company, got total=%d len=%d", total, len(items))
	}
}

func TestCruiseService_ListPublic(t *testing.T) {
	svc := NewCruiseService(&mockCruiseRepo{}, &mockCabinRepo{}, &mockCompanyRepo{})
	items, total, err := svc.ListPublic(context.Background(), 2, "", "", 1, 20)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].CompanyID != 2 {
		t.Fatalf("expected one public cruise for company 2, got total=%d len=%d", total, len(items))
	}
}
