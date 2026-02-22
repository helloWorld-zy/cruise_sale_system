package domain

import "time"

type FacilityCategory struct {
	ID        int64  `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`
	Icon      string `gorm:"size:255"`
	SortOrder int    `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
