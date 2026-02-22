package domain

import "time"

type Role struct {
	ID          int64  `gorm:"primaryKey"`
	Name        string `gorm:"size:50;uniqueIndex"`
	DisplayName string `gorm:"size:100"`
	Description string `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
