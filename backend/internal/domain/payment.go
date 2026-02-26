package domain

import "time"

type Payment struct {
	ID          int64  `gorm:"primaryKey"`
	OrderID     int64  `gorm:"index;column:order_id"`
	Provider    string `gorm:"size:20"`
	TradeNo     string `gorm:"size:100"`
	AmountCents int64
	Status      string `gorm:"size:20"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
