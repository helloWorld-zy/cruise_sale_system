package repository

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

type voyagePriceAgg struct {
	VoyageID      int64
	MinPriceCents int64
}

type voyageSoldAgg struct {
	VoyageID  int64
	SoldCount int64
}

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
				"cruise_id":                   v.CruiseID,
				"code":                        v.Code,
				"image_url":                   v.ImageURL,
				"brief_info":                  v.BriefInfo,
				"depart_date":                 v.DepartDate,
				"return_date":                 v.ReturnDate,
				"status":                      v.Status,
				"fee_note_template_id":        v.FeeNoteTemplateID,
				"fee_note_mode":               v.FeeNoteMode,
				"fee_note_content_json":       v.FeeNoteContentJSON,
				"booking_notice_template_id":  v.BookingNoticeTemplateID,
				"booking_notice_mode":         v.BookingNoticeMode,
				"booking_notice_content_json": v.BookingNoticeContentJSON,
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
		Preload("Cruise.Company").
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
	out.FeeNoteContent = decodeFeeNoteJSON(out.FeeNoteContentJSON)
	out.BookingNoticeContent = decodeBookingNoticeJSON(out.BookingNoticeContentJSON)
	out.FeeNote = out.FeeNoteContent
	out.BookingNotice = out.BookingNoticeContent
	if err := r.resolveVoyageTemplateContent(ctx, &out); err != nil {
		return nil, err
	}
	items := []domain.Voyage{out}
	if err := r.enrichVoyageMetrics(ctx, items); err != nil {
		return nil, err
	}
	out = items[0]
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
			out[i].FeeNoteContent = decodeFeeNoteJSON(out[i].FeeNoteContentJSON)
			out[i].BookingNoticeContent = decodeBookingNoticeJSON(out[i].BookingNoticeContentJSON)
			continue
		}
		out[i].ItineraryDays = out[i].Itineraries[len(out[i].Itineraries)-1].DayNo
		out[i].FirstStopCity = out[i].Itineraries[0].City
		out[i].FeeNoteContent = decodeFeeNoteJSON(out[i].FeeNoteContentJSON)
		out[i].BookingNoticeContent = decodeBookingNoticeJSON(out[i].BookingNoticeContentJSON)
	}
	return out, nil
}

// ListPublic 查询前台可见航次，支持按邮轮过滤和分页。
func (r *VoyageRepository) ListPublic(ctx context.Context, cruiseID int64, keyword string, page, pageSize int) ([]domain.Voyage, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	query := r.db.WithContext(ctx).
		Model(&domain.Voyage{}).
		Joins("LEFT JOIN cruises ON cruises.id = voyages.cruise_id").
		Where("voyages.status = ?", 1)
	if cruiseID > 0 {
		query = query.Where("voyages.cruise_id = ?", cruiseID)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			"voyages.brief_info LIKE ? OR voyages.code LIKE ? OR cruises.name LIKE ? OR cruises.english_name LIKE ?",
			like,
			like,
			like,
			like,
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var out []domain.Voyage
	err := query.
		Preload("Cruise.Company").
		Preload("Itineraries", func(db *gorm.DB) *gorm.DB {
			return db.Order("day_no asc, stop_index asc")
		}).
		Order("cruises.sort_order desc, cruises.name asc, voyages.depart_date asc, voyages.id asc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&out).Error
	if err != nil {
		return nil, 0, err
	}

	for i := range out {
		if len(out[i].Itineraries) == 0 {
			out[i].FeeNoteContent = decodeFeeNoteJSON(out[i].FeeNoteContentJSON)
			out[i].BookingNoticeContent = decodeBookingNoticeJSON(out[i].BookingNoticeContentJSON)
			continue
		}
		out[i].ItineraryDays = out[i].Itineraries[len(out[i].Itineraries)-1].DayNo
		out[i].FirstStopCity = out[i].Itineraries[0].City
		out[i].FeeNoteContent = decodeFeeNoteJSON(out[i].FeeNoteContentJSON)
		out[i].BookingNoticeContent = decodeBookingNoticeJSON(out[i].BookingNoticeContentJSON)
	}
	if err := r.enrichVoyageMetrics(ctx, out); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (r *VoyageRepository) resolveVoyageTemplateContent(ctx context.Context, v *domain.Voyage) error {
	if v == nil {
		return nil
	}
	if v.FeeNoteMode == domain.VoyageContentModeTemplate && v.FeeNoteTemplateID > 0 {
		var tpl domain.ContentTemplate
		if err := r.db.WithContext(ctx).First(&tpl, v.FeeNoteTemplateID).Error; err == nil {
			v.FeeNote = decodeFeeNoteJSON(tpl.ContentJSON)
		}
	}
	if v.FeeNote == nil {
		v.FeeNote = v.FeeNoteContent
	}
	if v.BookingNoticeMode == domain.VoyageContentModeTemplate && v.BookingNoticeTemplateID > 0 {
		var tpl domain.ContentTemplate
		if err := r.db.WithContext(ctx).First(&tpl, v.BookingNoticeTemplateID).Error; err == nil {
			v.BookingNotice = decodeBookingNoticeJSON(tpl.ContentJSON)
		}
	}
	if v.BookingNotice == nil {
		v.BookingNotice = v.BookingNoticeContent
	}
	return nil
}

func decodeFeeNoteJSON(raw string) *domain.FeeNoteContent {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var out domain.FeeNoteContent
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return &out
}

func decodeBookingNoticeJSON(raw string) *domain.BookingNoticeContent {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var out domain.BookingNoticeContent
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return &out
}

func (r *VoyageRepository) enrichVoyageMetrics(ctx context.Context, voyages []domain.Voyage) error {
	if len(voyages) == 0 {
		return nil
	}
	voyageIDs := make([]int64, 0, len(voyages))
	for _, item := range voyages {
		voyageIDs = append(voyageIDs, item.ID)
	}

	var priceRows []voyagePriceAgg
	if err := r.db.WithContext(ctx).
		Model(&domain.VoyageCabinTypeCurrent{}).
		Select("voyage_id, MIN(sale_price_cents) AS min_price_cents").
		Where("voyage_id IN ?", voyageIDs).
		Group("voyage_id").
		Scan(&priceRows).Error; err != nil {
		return err
	}
	priceMap := make(map[int64]int64, len(priceRows))
	for _, row := range priceRows {
		priceMap[row.VoyageID] = row.MinPriceCents
	}

	var soldRows []voyageSoldAgg
	if err := r.db.WithContext(ctx).
		Model(&domain.Booking{}).
		Select("voyage_id, COUNT(*) AS sold_count").
		Where("voyage_id IN ?", voyageIDs).
		Group("voyage_id").
		Scan(&soldRows).Error; err != nil {
		return err
	}
	soldMap := make(map[int64]int64, len(soldRows))
	for _, row := range soldRows {
		soldMap[row.VoyageID] = row.SoldCount
	}

	for index := range voyages {
		voyages[index].MinPriceCents = priceMap[voyages[index].ID]
		voyages[index].SoldCount = soldMap[voyages[index].ID]
	}
	return nil
}

// Delete 删除指定的航次记录。
func (r *VoyageRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Voyage{}, id).Error
}

// 编译时接口实现检查
var _ domain.VoyageRepository = (*VoyageRepository)(nil)
