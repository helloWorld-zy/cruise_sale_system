package domain

import "time"

type Voyage struct {
	ID         int64  `gorm:"primaryKey"`
	RouteID    int64  `gorm:"index"`
	CruiseID   int64  `gorm:"index"`
	Code       string `gorm:"size:50;uniqueIndex"`
	DepartDate time.Time
	ReturnDate time.Time
	Status     int16 `gorm:"default:1"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
