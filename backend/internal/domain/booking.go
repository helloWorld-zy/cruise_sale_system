package domain

import "time"

// Booking 表示用户对航次舱位的预订订单。
type Booking struct {
	ID         int64 `gorm:"primaryKey"`
	UserID     int64 `gorm:"index"`
	VoyageID   int64
	CabinSKUID int64  `gorm:"column:cabin_sku_id"`
	Status     string `gorm:"size:20"`
	TotalCents int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
