package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// CabinTypeMediaRepository 提供舱型媒体的数据访问实现。
type CabinTypeMediaRepository struct {
	db *gorm.DB
}

func NewCabinTypeMediaRepository(db *gorm.DB) *CabinTypeMediaRepository {
	return &CabinTypeMediaRepository{db: db}
}

func (r *CabinTypeMediaRepository) Create(ctx context.Context, media *domain.CabinTypeMedia) error {
	return r.db.WithContext(ctx).Create(media).Error
}

func (r *CabinTypeMediaRepository) Update(ctx context.Context, media *domain.CabinTypeMedia) error {
	return r.db.WithContext(ctx).Save(media).Error
}

func (r *CabinTypeMediaRepository) GetByID(ctx context.Context, id int64) (*domain.CabinTypeMedia, error) {
	var item domain.CabinTypeMedia
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CabinTypeMediaRepository) ListByCabinType(ctx context.Context, cabinTypeID int64) ([]domain.CabinTypeMedia, error) {
	var items []domain.CabinTypeMedia
	if err := r.db.WithContext(ctx).
		Model(&domain.CabinTypeMedia{}).
		Where("cabin_type_id = ?", cabinTypeID).
		Order("is_primary desc, sort_order desc, id desc").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CabinTypeMediaRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.CabinTypeMedia{}, id).Error
}

func (r *CabinTypeMediaRepository) SetPrimary(ctx context.Context, cabinTypeID int64, mediaType string, mediaID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.CabinTypeMedia{}).
			Where("cabin_type_id = ? AND media_type = ?", cabinTypeID, mediaType).
			Updates(map[string]any{"is_primary": false}).Error; err != nil {
			return err
		}
		return tx.Model(&domain.CabinTypeMedia{}).
			Where("id = ? AND cabin_type_id = ? AND media_type = ?", mediaID, cabinTypeID, mediaType).
			Updates(map[string]any{"is_primary": true}).Error
	})
}
