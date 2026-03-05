package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type fakeVoyageCabinTypePriceRepo struct {
	created      *domain.VoyageCabinTypePriceVersion
	upserted     *domain.VoyageCabinTypeCurrent
	latest       *domain.VoyageCabinTypePriceVersion
	currentAt    *domain.VoyageCabinTypeCurrent
	currentAtErr error
}

func (f *fakeVoyageCabinTypePriceRepo) CreateVersion(ctx context.Context, version *domain.VoyageCabinTypePriceVersion) error {
	copied := *version
	f.created = &copied
	if f.created.ID == 0 {
		f.created.ID = 101
		version.ID = 101
	}
	return nil
}

func (f *fakeVoyageCabinTypePriceRepo) UpsertCurrent(ctx context.Context, current *domain.VoyageCabinTypeCurrent) error {
	copied := *current
	f.upserted = &copied
	return nil
}

func (f *fakeVoyageCabinTypePriceRepo) GetCurrent(ctx context.Context, voyageID, cabinTypeID int64) (*domain.VoyageCabinTypeCurrent, error) {
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeVoyageCabinTypePriceRepo) GetCurrentAt(ctx context.Context, voyageID, cabinTypeID int64, at time.Time) (*domain.VoyageCabinTypeCurrent, error) {
	if f.currentAtErr != nil {
		return nil, f.currentAtErr
	}
	if f.currentAt == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return f.currentAt, nil
}

func (f *fakeVoyageCabinTypePriceRepo) ListCurrentByVoyages(ctx context.Context, voyageIDs []int64) ([]domain.VoyageCabinTypeCurrent, error) {
	return []domain.VoyageCabinTypeCurrent{}, nil
}

func (f *fakeVoyageCabinTypePriceRepo) ListCurrentByVoyagesAt(ctx context.Context, voyageIDs []int64, at time.Time) ([]domain.VoyageCabinTypeCurrent, error) {
	return []domain.VoyageCabinTypeCurrent{}, nil
}

func (f *fakeVoyageCabinTypePriceRepo) GetLatestVersionAt(ctx context.Context, voyageID, cabinTypeID int64, at time.Time) (*domain.VoyageCabinTypePriceVersion, error) {
	return f.latest, nil
}

func (f *fakeVoyageCabinTypePriceRepo) ListVersions(ctx context.Context, voyageID, cabinTypeID int64, page, pageSize int) ([]domain.VoyageCabinTypePriceVersion, int64, error) {
	return []domain.VoyageCabinTypePriceVersion{}, 0, nil
}

func TestVoyageCabinTypePriceService_ApplyVersionAndRefreshCurrent(t *testing.T) {
	repo := &fakeVoyageCabinTypePriceRepo{
		latest: &domain.VoyageCabinTypePriceVersion{
			ID:                   88,
			VoyageID:             12,
			CabinTypeID:          34,
			InventoryTotal:       20,
			SettlementPriceCents: 12300,
			SalePriceCents:       16600,
			EffectiveAt:          time.Now().Add(-time.Minute),
		},
	}
	svc := NewVoyageCabinTypePriceService(repo)

	err := svc.ApplyVersionAndRefreshCurrent(context.Background(), &domain.VoyageCabinTypePriceVersion{
		VoyageID:             12,
		CabinTypeID:          34,
		InventoryTotal:       10,
		SettlementPriceCents: 12000,
		SalePriceCents:       16800,
		EffectiveAt:          time.Now().Add(2 * time.Hour),
	})

	assert.NoError(t, err)
	assert.NotNil(t, repo.created)
	assert.NotNil(t, repo.upserted)
	assert.Equal(t, int64(12), repo.upserted.VoyageID)
	assert.Equal(t, int64(34), repo.upserted.CabinTypeID)
	assert.Equal(t, int64(88), repo.upserted.VersionID)
	assert.Equal(t, int64(16600), repo.upserted.SalePriceCents)
}

func TestVoyageCabinTypePriceService_GetCurrentNotFoundReturnsNil(t *testing.T) {
	repo := &fakeVoyageCabinTypePriceRepo{currentAtErr: gorm.ErrRecordNotFound}
	svc := NewVoyageCabinTypePriceService(repo)

	current, err := svc.GetCurrent(context.Background(), 1, 2)
	assert.NoError(t, err)
	assert.Nil(t, current)
}

func TestVoyageCabinTypePriceService_GetCurrentPropagatesError(t *testing.T) {
	repo := &fakeVoyageCabinTypePriceRepo{currentAtErr: errors.New("db down")}
	svc := NewVoyageCabinTypePriceService(repo)

	current, err := svc.GetCurrent(context.Background(), 1, 2)
	assert.Error(t, err)
	assert.Nil(t, current)
}
