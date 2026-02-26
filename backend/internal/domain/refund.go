package domain

import "time"

type Refund struct {
    ID int64 `gorm:"primaryKey"`
    PaymentID int64 `gorm:"index"`
    AmountCents int64
    Reason string `gorm:"size:200"`
    Status string `gorm:"size:20"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
