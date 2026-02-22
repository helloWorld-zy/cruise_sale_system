package domain

import "time"

type Image struct {
	ID         int64  `gorm:"primaryKey"`
	EntityType string `gorm:"size:50;index;not null"`
	EntityID   int64  `gorm:"index;not null"`
	URL        string `gorm:"size:500;not null"`
	SortOrder  int    `gorm:"default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
