package domain

import "time"

type CabinType struct {
	ID          int64  `gorm:"primaryKey"`
	CruiseID    int64  `gorm:"index;not null"`
	Name        string `gorm:"size:100;not null"`
	EnglishName string `gorm:"size:100"`
	Capacity    int    `gorm:"default:2"`
	Area        float64
	Deck        string `gorm:"size:50"`
	Description string `gorm:"type:text"`
	Status      int16  `gorm:"default:1"`
	SortOrder   int    `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}
