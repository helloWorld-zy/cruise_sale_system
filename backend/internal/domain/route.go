package domain

import "time"

type Route struct {
	ID          int64  `gorm:"primaryKey"`
	Code        string `gorm:"size:50;uniqueIndex"`
	Name        string `gorm:"size:200"`
	Description string `gorm:"type:text"`
	Status      int16  `gorm:"default:1"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
