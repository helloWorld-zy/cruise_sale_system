package domain

import "time"

type Facility struct {
	ID          int64             `gorm:"primaryKey"`
	CategoryID  int64             `gorm:"index;not null"`
	Category    *FacilityCategory `gorm:"foreignKey:CategoryID"`
	CruiseID    int64             `gorm:"index;not null"`
	Name        string            `gorm:"size:100;not null"`
	EnglishName string            `gorm:"size:100"`
	Location    string            `gorm:"size:100"`
	Description string            `gorm:"type:text"`
	Status      int16             `gorm:"default:1"`
	SortOrder   int               `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}
