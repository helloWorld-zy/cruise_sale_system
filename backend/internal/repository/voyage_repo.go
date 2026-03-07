package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// VoyageRepository 实现 domain.VoyageRepository 接口，提供航次实体的数据库操作。
type VoyageRepository struct{ db *gorm.DB }

// NewVoyageRepository 创建航次仓储实例。
func NewVoyageRepository(db *gorm.DB) *VoyageRepository { return &VoyageRepository{db: db} }

// Create 插入一条新的航次记录。
func (r *VoyageRepository) Create(ctx context.Context, v *domain.Voyage) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		itineraries := v.Itineraries
		v.Itineraries = nil
		if err := tx.Create(v).Error; err != nil {
			return err
		}
		if len(itineraries) == 0 {
			return nil
		}
		for i := range itineraries {
			itineraries[i].VoyageID = v.ID
		}
		return tx.Create(&itineraries).Error
	})
}

// Update 保存航次的所有字段修改。
func (r *VoyageRepository) Update(ctx context.Context, v *domain.Voyage) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.Voyage{}).Where("id = ?", v.ID).
			Updates(map[string]any{
				"cruise_id":   v.CruiseID,
				"code":        v.Code,
				"image_url":   v.ImageURL,
				"brief_info":  v.BriefInfo,
				"depart_date": v.DepartDate,
				"return_date": v.ReturnDate,
				"status":      v.Status,
			}).Error; err != nil {
			return err
		}
		if err := tx.Where("voyage_id = ?", v.ID).Delete(&domain.VoyageItinerary{}).Error; err != nil {
			return err
		}
		if len(v.Itineraries) == 0 {
			return nil
		}
		items := make([]domain.VoyageItinerary, len(v.Itineraries))
		copy(items, v.Itineraries)
		for i := range items {
			items[i].ID = 0
			items[i].VoyageID = v.ID
		}
		return tx.Create(&items).Error
	})
}

// GetByID 根据主键查询航次记录。
func (r *VoyageRepository) GetByID(ctx context.Context, id int64) (*domain.Voyage, error) {
	var out domain.Voyage
	if err := r.db.WithContext(ctx).
		Preload("Itineraries", func(db *gorm.DB) *gorm.DB {
			return db.Order("day_no asc, stop_index asc")
		}).
		First(&out, id).Error; err != nil {
		return nil, err
	}
	if len(out.Itineraries) > 0 {
		out.ItineraryDays = out.Itineraries[len(out.Itineraries)-1].DayNo
		out.FirstStopCity = out.Itineraries[0].City
	}
	return &out, nil
}

// List 查询所有航次，按出发日期升序排列。
func (r *VoyageRepository) List(ctx context.Context) ([]domain.Voyage, error) {
	var out []domain.Voyage
	err := r.db.WithContext(ctx).
		Preload("Itineraries", func(db *gorm.DB) *gorm.DB {
			return db.Order("day_no asc, stop_index asc")
		}).
		Order("depart_date asc").
		Find(&out).Error
	if err != nil {
		return nil, err
	}
	for i := range out {
		if len(out[i].Itineraries) == 0 {
			continue
		}
		out[i].ItineraryDays = out[i].Itineraries[len(out[i].Itineraries)-1].DayNo
		out[i].FirstStopCity = out[i].Itineraries[0].City
	}
	return out, nil
}

// Delete 删除指定的航次记录。
func (r *VoyageRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Voyage{}, id).Error
}

// 编译时接口实现检查
var _ domain.VoyageRepository = (*VoyageRepository)(nil)
