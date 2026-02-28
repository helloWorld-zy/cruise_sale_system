package domain

import "time"

// Refund 表示退款记录实体。
type Refund struct {
	ID          int64     `gorm:"primaryKey"` // 主键 ID
	PaymentID   int64     `gorm:"index"`      // 关联的支付记录 ID
	AmountCents int64     // 退款金额（单位：分）
	Reason      string    `gorm:"size:200"` // 退款原因
	Status      string    `gorm:"size:20"`  // 退款状态（pending / approved / cancelled）
	CreatedAt   time.Time // 创建时间
	UpdatedAt   time.Time // 更新时间
}
