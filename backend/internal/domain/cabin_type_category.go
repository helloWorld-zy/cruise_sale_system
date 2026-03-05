package domain

import "time"

// CabinTypeCategory 表示舱型大类字典（如内舱、海景、阳台、套房）。
type CabinTypeCategory struct {
	ID        int64      `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"size:64;not null" json:"name"`
	Code      string     `gorm:"size:32;not null;uniqueIndex" json:"code"`
	Status    int16      `gorm:"default:1" json:"status"`
	SortOrder int        `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
