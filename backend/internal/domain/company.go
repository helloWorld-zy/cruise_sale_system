package domain

import "time"

type CruiseCompany struct {
	ID          int64  `gorm:"primaryKey"`
	Name        string `gorm:"size:100;not null"`
	EnglishName string `gorm:"size:100"`
	Description string `gorm:"type:text"`
	LogoURL     string `gorm:"size:500"`
	Status      int16  `gorm:"default:1"`
	SortOrder   int    `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}
