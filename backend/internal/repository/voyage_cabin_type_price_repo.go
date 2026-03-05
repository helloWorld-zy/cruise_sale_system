package repository

import (
	"context"
	"errors"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// VoyageCabinTypePriceRepository 提供航次舱型价格版本与当前态的数据访问实现。
type VoyageCabinTypePriceRepository struct {
	db *gorm.DB
}

func NewVoyageCabinTypePriceRepository(db *gorm.DB) *VoyageCabinTypePriceRepository {
	return &VoyageCabinTypePriceRepository{db: db}
}

func (r *VoyageCabinTypePriceRepository) CreateVersion(ctx context.Context, version *domain.VoyageCabinTypePriceVersion) error {
	return r.db.WithContext(ctx).Create(version).Error
}

func (r *VoyageCabinTypePriceRepository) UpsertCurrent(ctx context.Context, current *domain.VoyageCabinTypeCurrent) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "voyage_id"}, {Name: "cabin_type_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"inventory_total", "settlement_price_cents", "sale_price_cents", "effective_at", "version_id", "updated_at"}),
	}).Create(current).Error
}

func (r *VoyageCabinTypePriceRepository) GetCurrent(ctx context.Context, voyageID, cabinTypeID int64) (*domain.VoyageCabinTypeCurrent, error) {
	var item domain.VoyageCabinTypeCurrent
	if err := r.db.WithContext(ctx).Where("voyage_id = ? AND cabin_type_id = ?", voyageID, cabinTypeID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *VoyageCabinTypePriceRepository) ListCurrentByVoyages(ctx context.Context, voyageIDs []int64) ([]domain.VoyageCabinTypeCurrent, error) {
	if len(voyageIDs) == 0 {
		return []domain.VoyageCabinTypeCurrent{}, nil
	}
	var items []domain.VoyageCabinTypeCurrent
	if err := r.db.WithContext(ctx).
		Model(&domain.VoyageCabinTypeCurrent{}).
		Where("voyage_id IN ?", voyageIDs).
		Order("voyage_id asc, cabin_type_id asc").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *VoyageCabinTypePriceRepository) GetCurrentAt(ctx context.Context, voyageID, cabinTypeID int64, at time.Time) (*domain.VoyageCabinTypeCurrent, error) {
	var item domain.VoyageCabinTypeCurrent
	err := r.db.WithContext(ctx).
		Where("voyage_id = ? AND cabin_type_id = ? AND effective_at <= ?", voyageID, cabinTypeID, at).
		First(&item).Error
	if err == nil {
		return &item, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	version, err := r.GetLatestVersionAt(ctx, voyageID, cabinTypeID, at)
	if err != nil {
		return nil, err
	}
	if version == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return &domain.VoyageCabinTypeCurrent{
		VoyageID:             version.VoyageID,
		CabinTypeID:          version.CabinTypeID,
		InventoryTotal:       version.InventoryTotal,
		SettlementPriceCents: version.SettlementPriceCents,
		SalePriceCents:       version.SalePriceCents,
		EffectiveAt:          version.EffectiveAt,
		VersionID:            version.ID,
		UpdatedAt:            version.CreatedAt,
	}, nil
}

func (r *VoyageCabinTypePriceRepository) ListCurrentByVoyagesAt(ctx context.Context, voyageIDs []int64, at time.Time) ([]domain.VoyageCabinTypeCurrent, error) {
	if len(voyageIDs) == 0 {
		return []domain.VoyageCabinTypeCurrent{}, nil
	}

	var items []domain.VoyageCabinTypeCurrent
	if err := r.db.WithContext(ctx).
		Model(&domain.VoyageCabinTypeCurrent{}).
		Where("voyage_id IN ? AND effective_at <= ?", voyageIDs, at).
		Order("voyage_id asc, cabin_type_id asc").
		Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *VoyageCabinTypePriceRepository) GetLatestVersionAt(ctx context.Context, voyageID, cabinTypeID int64, at time.Time) (*domain.VoyageCabinTypePriceVersion, error) {
	var item domain.VoyageCabinTypePriceVersion
	err := r.db.WithContext(ctx).
		Model(&domain.VoyageCabinTypePriceVersion{}).
		Where("voyage_id = ? AND cabin_type_id = ? AND effective_at <= ?", voyageID, cabinTypeID, at).
		Order("effective_at desc, id desc").
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *VoyageCabinTypePriceRepository) ListVersions(ctx context.Context, voyageID, cabinTypeID int64, page, pageSize int) ([]domain.VoyageCabinTypePriceVersion, int64, error) {
	var items []domain.VoyageCabinTypePriceVersion
	var total int64
	q := r.db.WithContext(ctx).
		Model(&domain.VoyageCabinTypePriceVersion{}).
		Where("voyage_id = ? AND cabin_type_id = ?", voyageID, cabinTypeID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("effective_at desc, id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
