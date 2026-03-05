package service

import (
	"context"
	"errors"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

var shanghaiLocation = mustLoadLocation("Asia/Shanghai")

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.FixedZone("CST", 8*60*60)
	}
	return loc
}

// VoyageCabinTypePriceService 提供航次舱型价格版本与当前态管理能力。
type VoyageCabinTypePriceService struct {
	repo domain.VoyageCabinTypePriceRepository
}

func NewVoyageCabinTypePriceService(repo domain.VoyageCabinTypePriceRepository) *VoyageCabinTypePriceService {
	return &VoyageCabinTypePriceService{repo: repo}
}

func (s *VoyageCabinTypePriceService) CreateVersion(ctx context.Context, version *domain.VoyageCabinTypePriceVersion) error {
	return s.repo.CreateVersion(ctx, version)
}

func (s *VoyageCabinTypePriceService) UpsertCurrent(ctx context.Context, current *domain.VoyageCabinTypeCurrent) error {
	return s.repo.UpsertCurrent(ctx, current)
}

// ApplyVersionAndSetCurrent 先写入历史版本，再把该版本设置为当前价。
func (s *VoyageCabinTypePriceService) ApplyVersionAndSetCurrent(ctx context.Context, version *domain.VoyageCabinTypePriceVersion) error {
	if version.EffectiveAt.IsZero() {
		version.EffectiveAt = time.Now().In(shanghaiLocation)
	}
	if err := s.repo.CreateVersion(ctx, version); err != nil {
		return err
	}

	// Keep future prices in history only; current table must represent the latest effective record at "now".
	if version.EffectiveAt.After(time.Now().In(shanghaiLocation)) {
		return nil
	}

	current := &domain.VoyageCabinTypeCurrent{
		VoyageID:             version.VoyageID,
		CabinTypeID:          version.CabinTypeID,
		InventoryTotal:       version.InventoryTotal,
		SettlementPriceCents: version.SettlementPriceCents,
		SalePriceCents:       version.SalePriceCents,
		EffectiveAt:          version.EffectiveAt,
		VersionID:            version.ID,
		UpdatedAt:            time.Now(),
	}
	return s.repo.UpsertCurrent(ctx, current)
}

// ApplyVersionAndRefreshCurrent writes history and recalculates current by effective_at <= now(Asia/Shanghai).
func (s *VoyageCabinTypePriceService) ApplyVersionAndRefreshCurrent(ctx context.Context, version *domain.VoyageCabinTypePriceVersion) error {
	if version.EffectiveAt.IsZero() {
		version.EffectiveAt = time.Now().In(shanghaiLocation)
	}
	if err := s.repo.CreateVersion(ctx, version); err != nil {
		return err
	}

	now := time.Now().In(shanghaiLocation)
	latest, err := s.repo.GetLatestVersionAt(ctx, version.VoyageID, version.CabinTypeID, now)
	if err != nil {
		return err
	}
	if latest == nil {
		return nil
	}

	return s.repo.UpsertCurrent(ctx, &domain.VoyageCabinTypeCurrent{
		VoyageID:             latest.VoyageID,
		CabinTypeID:          latest.CabinTypeID,
		InventoryTotal:       latest.InventoryTotal,
		SettlementPriceCents: latest.SettlementPriceCents,
		SalePriceCents:       latest.SalePriceCents,
		EffectiveAt:          latest.EffectiveAt,
		VersionID:            latest.ID,
		UpdatedAt:            now,
	})
}

func (s *VoyageCabinTypePriceService) GetCurrent(ctx context.Context, voyageID, cabinTypeID int64) (*domain.VoyageCabinTypeCurrent, error) {
	now := time.Now().In(shanghaiLocation)
	current, err := s.repo.GetCurrentAt(ctx, voyageID, cabinTypeID, now)
	if err == nil {
		return current, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, err
}

func (s *VoyageCabinTypePriceService) ListCurrentByVoyages(ctx context.Context, voyageIDs []int64) ([]domain.VoyageCabinTypeCurrent, error) {
	now := time.Now().In(shanghaiLocation)
	return s.repo.ListCurrentByVoyagesAt(ctx, voyageIDs, now)
}

func (s *VoyageCabinTypePriceService) ListVersions(ctx context.Context, voyageID, cabinTypeID int64, page, pageSize int) ([]domain.VoyageCabinTypePriceVersion, int64, error) {
	return s.repo.ListVersions(ctx, voyageID, cabinTypeID, page, pageSize)
}
