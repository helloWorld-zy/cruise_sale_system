package domain

import "time"

// Voyage 表示一个具体的航次（某艘邮轮的一次出发）。
// 航次是舱房 SKU 和价格日历的核心关联维度。
type Voyage struct {
	ID         int64     `gorm:"primaryKey" json:"id"`            // 主键 ID
	RouteID    int64     `gorm:"-" json:"route_id,omitempty"`     // 已下线字段，仅用于兼容旧测试/旧请求
	CruiseID   int64     `gorm:"index" json:"cruise_id"`          // 执行邮轮 ID
	Code       string    `gorm:"size:50;uniqueIndex" json:"code"` // 航次编码（全局唯一）
	BriefInfo  string    `gorm:"size:300" json:"brief_info"`      // 航次简短信息（手动输入）
	DepartDate time.Time `json:"depart_date"`                     // 出发日期
	ReturnDate time.Time `json:"return_date"`                     // 返航日期
	Status     int16     `gorm:"default:1" json:"status"`         // 状态：1=开放预订，0=关闭
	CreatedAt  time.Time `json:"created_at"`                      // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`                      // 更新时间

	Itineraries   []VoyageItinerary `gorm:"foreignKey:VoyageID" json:"itineraries,omitempty"` // 航次行程明细
	ItineraryDays int               `gorm:"-" json:"itinerary_days,omitempty"`                // 行程天数（列表辅助字段）
	FirstStopCity string            `gorm:"-" json:"first_stop_city,omitempty"`               // 首日首站城市（列表辅助字段）
}

// VoyageItinerary 表示航次中某天某站的计划信息。
type VoyageItinerary struct {
	ID                int64     `gorm:"primaryKey" json:"id"`
	VoyageID          int64     `gorm:"index" json:"voyage_id"`
	DayNo             int       `json:"day_no"`
	StopIndex         int       `json:"stop_index"`
	City              string    `gorm:"size:120" json:"city"`
	Summary           string    `gorm:"type:text" json:"summary"`
	ETATime           *string   `gorm:"column:eta_time" json:"eta_time,omitempty"`
	ETDTime           *string   `gorm:"column:etd_time" json:"etd_time,omitempty"`
	HasBreakfast      bool      `json:"has_breakfast"`
	HasLunch          bool      `json:"has_lunch"`
	HasDinner         bool      `json:"has_dinner"`
	HasAccommodation  bool      `json:"has_accommodation"`
	AccommodationText string    `gorm:"size:300" json:"accommodation_text"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
