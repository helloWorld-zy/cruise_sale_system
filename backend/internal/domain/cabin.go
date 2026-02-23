package domain

import "time"

type CabinSKU struct {
	ID          int64  `gorm:"primaryKey"`
	VoyageID    int64  `gorm:"index"`
	CabinTypeID int64  `gorm:"index"`
	Code        string `gorm:"size:80;uniqueIndex"`
	Deck        string `gorm:"size:20"`
	Area        float64
	MaxGuests   int
	Status      int16 `gorm:"default:1"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CabinPrice struct {
	ID         int64     `gorm:"primaryKey"`
	CabinSKUID int64     `gorm:"index"`
	Date       time.Time `gorm:"index"`
	Occupancy  int
	PriceCents int64 `gorm:"column:price_cents"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CabinInventory struct {
	ID         int64 `gorm:"primaryKey"`
	CabinSKUID int64 `gorm:"uniqueIndex"`
	Total      int
	Locked     int
	Sold       int
	UpdatedAt  time.Time
}

type InventoryLog struct {
	ID         int64 `gorm:"primaryKey"`
	CabinSKUID int64 `gorm:"index"`
	Change     int
	Reason     string `gorm:"size:200"`
	CreatedAt  time.Time
}
