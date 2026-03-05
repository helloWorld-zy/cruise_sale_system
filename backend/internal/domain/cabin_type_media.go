package domain

import "time"

const (
	CabinTypeMediaTypeImage     = "image"
	CabinTypeMediaTypeFloorPlan = "floor_plan"
)

// CabinTypeMedia 表示舱型图片或平面图等媒体。
type CabinTypeMedia struct {
	ID          int64      `gorm:"primaryKey" json:"id"`
	CabinTypeID int64      `gorm:"index;not null" json:"cabin_type_id"`
	MediaType   string     `gorm:"size:20;not null" json:"media_type"`
	URL         string     `gorm:"type:text;not null" json:"url"`
	Title       string     `gorm:"size:120;not null" json:"title"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	IsPrimary   bool       `gorm:"default:false" json:"is_primary"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
