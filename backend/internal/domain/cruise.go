package domain

import "time"

type Cruise struct {
	ID                int64          `gorm:"primaryKey"`
	CompanyID         int64          `gorm:"index;not null"`
	Company           *CruiseCompany `gorm:"foreignKey:CompanyID"`
	Name              string         `gorm:"size:100;not null"`
	EnglishName       string         `gorm:"size:100"`
	BuildYear         int
	Tonnage           float64
	PassengerCapacity int
	RoomCount         int
	Description       string `gorm:"type:text"`
	Status            int16  `gorm:"default:1"`
	SortOrder         int    `gorm:"default:0"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time `gorm:"index"`
}
