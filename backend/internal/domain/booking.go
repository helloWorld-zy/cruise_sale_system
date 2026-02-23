package domain

import "time"

type Booking struct {
	ID         int64 `gorm:"primaryKey"`
	UserID     int64 `gorm:"index"`
	VoyageID   int64
	CabinSKUID int64
	Status     string `gorm:"size:20"`
	TotalCents int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
