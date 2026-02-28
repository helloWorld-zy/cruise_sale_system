package domain

import "time"

// Booking 表示用户对航次舱位的预订订单。
type Booking struct {
	ID         int64     `gorm:"primaryKey" json:"id"`                    // 主键 ID
	UserID     int64     `gorm:"index" json:"user_id"`                    // 下单用户 ID
	VoyageID   int64     `json:"voyage_id"`                               // 所属航次 ID
	CabinSKUID int64     `gorm:"column:cabin_sku_id" json:"cabin_sku_id"` // 预订的舱房 SKU ID
	Status     string    `gorm:"size:20" json:"status"`                   // 订单状态（created/paid/cancelled）
	TotalCents int64     `json:"total_cents"`                             // 订单总金额（单位：分）
	CreatedAt  time.Time `json:"created_at"`                              // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`                              // 更新时间
}
