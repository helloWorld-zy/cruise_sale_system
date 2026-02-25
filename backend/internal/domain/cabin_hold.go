package domain

import "time"

// CabinHold 表示用户对舱位库存的临时占用记录。
type CabinHold struct {
	ID         int64 `gorm:"primaryKey"`
	CabinSKUID int64 `gorm:"column:cabin_sku_id;index"`
	UserID     int64 `gorm:"index"`
	Qty        int
	ExpiresAt  time.Time `gorm:"index"`
	CreatedAt  time.Time
}
