package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// Replicate mocks here for service testing

type allMockCompanyRepo struct{}

func (m *allMockCompanyRepo) Create(ctx context.Context, company *domain.CruiseCompany) error {
	if company.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *allMockCompanyRepo) Update(ctx context.Context, company *domain.CruiseCompany) error {
	if company.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *allMockCompanyRepo) GetByID(ctx context.Context, id int64) (*domain.CruiseCompany, error) {
	if id == 1 {
		return &domain.CruiseCompany{ID: 1}, nil
	}
	if id == 99 {
		return nil, errors.New("error")
	}
	return nil, nil // not found
}
func (m *allMockCompanyRepo) List(ctx context.Context, keyword string, page, pageSize int) ([]domain.CruiseCompany, int64, error) {
	if keyword == "error" {
		return nil, 0, errors.New("error")
	}
	return []domain.CruiseCompany{{ID: 1}}, 1, nil
}
func (m *allMockCompanyRepo) Delete(ctx context.Context, id int64) error {
	if id == 99 {
		return errors.New("error")
	}
	return nil
}

type allMockCruiseRepo struct{}

func (m *allMockCruiseRepo) Create(ctx context.Context, cruise *domain.Cruise) error {
	if cruise.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *allMockCruiseRepo) Update(ctx context.Context, cruise *domain.Cruise) error {
	if cruise.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *allMockCruiseRepo) GetByID(ctx context.Context, id int64) (*domain.Cruise, error) {
	if id == 1 || id == 88 {
		return &domain.Cruise{ID: 1, CompanyID: 1}, nil
	}
	if id == 99 {
		return nil, errors.New("error")
	}
	return nil, nil // not found
}
func (m *allMockCruiseRepo) List(ctx context.Context, companyID int64, page, pageSize int) ([]domain.Cruise, int64, error) {
	if companyID == 99 {
		return nil, 0, errors.New("error")
	}
	if companyID == 88 {
		return []domain.Cruise{{ID: 1}}, 1, nil
	} // has cruises
	return nil, 0, nil
}
func (m *allMockCruiseRepo) Delete(ctx context.Context, id int64) error {
	if id == 99 {
		return errors.New("error")
	}
	return nil
}

type allMockCabinTypeRepo struct{}

func (m *allMockCabinTypeRepo) Create(ctx context.Context, c *domain.CabinType) error {
	if c.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *allMockCabinTypeRepo) Update(ctx context.Context, c *domain.CabinType) error {
	if c.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *allMockCabinTypeRepo) GetByID(ctx context.Context, id int64) (*domain.CabinType, error) {
	if id == 1 {
		return &domain.CabinType{ID: 1}, nil
	}
	if id == 99 {
		return nil, errors.New("error")
	}
	return nil, nil
}
func (m *allMockCabinTypeRepo) ListByCruise(ctx context.Context, id int64, p, ps int) ([]domain.CabinType, int64, error) {
	if id == 99 {
		return nil, 0, errors.New("error")
	}
	if id == 88 {
		return []domain.CabinType{{ID: 1}}, 1, nil
	}
	return nil, 0, nil
}
func (m *allMockCabinTypeRepo) Delete(ctx context.Context, id int64) error {
	if id == 99 {
		return errors.New("error")
	}
	return nil
}

type allMockFacilityCategoryRepo struct{}

func (m *allMockFacilityCategoryRepo) Create(ctx context.Context, c *domain.FacilityCategory) error {
	if c.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *allMockFacilityCategoryRepo) List(ctx context.Context) ([]domain.FacilityCategory, error) {
	return nil, nil
}
func (m *allMockFacilityCategoryRepo) Delete(ctx context.Context, id int64) error {
	if id == 99 {
		return errors.New("error")
	}
	return nil
}

type allMockFacilityRepo struct{}

func (m *allMockFacilityRepo) Create(ctx context.Context, f *domain.Facility) error {
	if f.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *allMockFacilityRepo) ListByCruise(ctx context.Context, cruiseID int64) ([]domain.Facility, error) {
	if cruiseID == 99 {
		return nil, errors.New("error")
	}
	return nil, nil
}
func (m *allMockFacilityRepo) Delete(ctx context.Context, id int64) error {
	if id == 99 {
		return errors.New("error")
	}
	return nil
}

func TestCompanyServiceAll(t *testing.T) {
	svc := NewCompanyService(&allMockCompanyRepo{}, &allMockCruiseRepo{})
	ctx := context.Background()

	svc.Create(ctx, &domain.CruiseCompany{Name: "test"})
	svc.Create(ctx, &domain.CruiseCompany{Name: "error"})

	svc.Update(ctx, &domain.CruiseCompany{Name: "test"})
	svc.Update(ctx, &domain.CruiseCompany{Name: "error"})

	svc.GetByID(ctx, 1)
	svc.GetByID(ctx, 2)
	svc.GetByID(ctx, 99)

	svc.List(ctx, "test", 1, 10)
	svc.List(ctx, "error", 1, 10)

	svc.Delete(ctx, 1)
	svc.Delete(ctx, 88) // has cruises
	svc.Delete(ctx, 99)
}

func TestCruiseServiceAll(t *testing.T) {
	svc := NewCruiseService(&allMockCruiseRepo{}, &allMockCabinTypeRepo{}, &allMockCompanyRepo{})
	ctx := context.Background()

	svc.Create(ctx, &domain.Cruise{Name: "test", CompanyID: 1})
	svc.Create(ctx, &domain.Cruise{Name: "test", CompanyID: 2})  // invalid comp
	svc.Create(ctx, &domain.Cruise{Name: "test", CompanyID: 99}) // error comp
	svc.Create(ctx, &domain.Cruise{Name: "error", CompanyID: 1}) // create err

	svc.Update(ctx, &domain.Cruise{Name: "test"})
	svc.Update(ctx, &domain.Cruise{Name: "error"})

	svc.GetByID(ctx, 1)
	svc.GetByID(ctx, 2)
	svc.GetByID(ctx, 99)

	svc.List(ctx, 1, 1, 10)
	svc.List(ctx, 99, 1, 10)

	svc.Delete(ctx, 1)  // success
	svc.Delete(ctx, 88) // has cabin types
	svc.Delete(ctx, 99)
}

func TestCabinTypeServiceAll(t *testing.T) {
	svc := NewCabinTypeService(&allMockCabinTypeRepo{})
	ctx := context.Background()

	svc.Create(ctx, &domain.CabinType{Name: "test"})
	svc.Create(ctx, &domain.CabinType{Name: "error"})

	svc.Update(ctx, &domain.CabinType{Name: "test"})
	svc.Update(ctx, &domain.CabinType{Name: "error"})

	svc.GetByID(ctx, 1)
	svc.GetByID(ctx, 2)
	svc.GetByID(ctx, 99)

	svc.List(ctx, 1, 1, 10)
	svc.List(ctx, 99, 1, 10)

	svc.Delete(ctx, 1)
	svc.Delete(ctx, 99)
}

func TestFacilitySvcAll(t *testing.T) {
	fcs := NewFacilityCategoryService(&allMockFacilityCategoryRepo{})
	ctx := context.Background()
	fcs.Create(ctx, &domain.FacilityCategory{Name: "test"})
	fcs.Create(ctx, &domain.FacilityCategory{Name: "error"})
	fcs.List(ctx)
	fcs.Delete(ctx, 1)
	fcs.Delete(ctx, 99)

	fs := NewFacilityService(&allMockFacilityRepo{})
	fs.Create(ctx, &domain.Facility{Name: "test"})
	fs.Create(ctx, &domain.Facility{Name: "error"})
	fs.ListByCruise(ctx, 1)
	fs.ListByCruise(ctx, 99)
	fs.Delete(ctx, 1)
	fs.Delete(ctx, 99)
}

// User AuthService
type dummyCodeStore struct{}

func (m *dummyCodeStore) Save(phone, code string) error {
	if phone == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *dummyCodeStore) Verify(phone, code string) bool {
	return code == "1234"
}

func TestUserAuthServiceAll(t *testing.T) {
	svc := NewUserAuthService(&dummyCodeStore{})
	svc.SendSMS("123", "1234")
	svc.SendSMS("error", "1234")
	svc.VerifySMS("123", "1234")
	svc.VerifySMS("123", "wrong")
	svc.WechatLogin("code")
}

// Misc Service tests
type mockInvRepo struct{}

func (m *mockInvRepo) GetBySKU(id int64) (domain.CabinInventory, error) {
	if id == 99 {
		return domain.CabinInventory{}, errors.New("err")
	}
	return domain.CabinInventory{}, nil
}
func (m *mockInvRepo) AdjustAtomic(ctx context.Context, id int64, d int) error { return nil }
func (m *mockInvRepo) AppendLog(ctx context.Context, l *domain.InventoryLog) error {
	return nil
}

type mockPriceRepo struct{}

func (m *mockPriceRepo) ListBySKU(ctx context.Context, id int64) ([]domain.CabinPrice, error) {
	if id == 99 {
		return nil, errors.New("error")
	}
	return []domain.CabinPrice{}, nil
}

type mockHoldSvc struct{}

func (m *mockHoldSvc) HoldWithTx(tx *gorm.DB, sku, u int64, q int) bool {
	_ = tx
	return sku != 99
}

type mockPriceSvc struct{}

func (m *mockPriceSvc) FindPrice(ctx context.Context, skuID int64, date time.Time, occ int) (int64, bool, error) {
	if skuID == 88 {
		return 0, false, errors.New("error")
	}
	return 100, true, nil
}

type mockBkRepo struct{}

func (m *mockBkRepo) Create(_ context.Context, b *domain.Booking) error { return nil }
func (m *mockBkRepo) InTx(fn func(tx *gorm.DB, create func(b *domain.Booking) error) error) error {
	return fn(nil, func(b *domain.Booking) error { return m.Create(context.Background(), b) })
}

func TestMiscServices(t *testing.T) {
	inv := NewInventoryService(&mockInvRepo{})
	inv.Available(1)
	inv.Available(99) // error branch

	price := NewPricingService(&mockPriceRepo{})
	price.FindPrice(context.Background(), 1, time.Now(), 2)
	price.FindPrice(context.Background(), 99, time.Now(), 2)

	bk := NewBookingService(&mockBkRepo{}, &mockPriceSvc{}, &mockHoldSvc{})
	bk.Create(context.Background(), 1, 1, 1, 2)
	bk.Create(context.Background(), 1, 1, 99, 2) // hold fail
	bk.Create(context.Background(), 1, 1, 88, 2) // price fail
}
