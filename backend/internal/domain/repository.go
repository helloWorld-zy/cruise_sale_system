package domain

import "context"

type CompanyRepository interface {
	Create(ctx context.Context, company *CruiseCompany) error
	Update(ctx context.Context, company *CruiseCompany) error
	GetByID(ctx context.Context, id int64) (*CruiseCompany, error)
	List(ctx context.Context, keyword string, page, pageSize int) ([]CruiseCompany, int64, error)
	Delete(ctx context.Context, id int64) error
}

type CruiseRepository interface {
	Create(ctx context.Context, cruise *Cruise) error
	Update(ctx context.Context, cruise *Cruise) error
	GetByID(ctx context.Context, id int64) (*Cruise, error)
	List(ctx context.Context, companyID int64, page, pageSize int) ([]Cruise, int64, error)
	Delete(ctx context.Context, id int64) error
}

type CabinTypeRepository interface {
	Create(ctx context.Context, cabinType *CabinType) error
	Update(ctx context.Context, cabinType *CabinType) error
	GetByID(ctx context.Context, id int64) (*CabinType, error)
	ListByCruise(ctx context.Context, cruiseID int64, page, pageSize int) ([]CabinType, int64, error)
	Delete(ctx context.Context, id int64) error
}

type FacilityCategoryRepository interface {
	Create(ctx context.Context, category *FacilityCategory) error
	List(ctx context.Context) ([]FacilityCategory, error)
	Delete(ctx context.Context, id int64) error
}

type FacilityRepository interface {
	Create(ctx context.Context, facility *Facility) error
	ListByCruise(ctx context.Context, cruiseID int64) ([]Facility, error)
	Delete(ctx context.Context, id int64) error
}
