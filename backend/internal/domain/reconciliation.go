package domain

import "time"

// Reconciliation 表示每日财务对账记录。
type Reconciliation struct {
	ID                 int64     `gorm:"primaryKey"`
	Date               time.Time `gorm:"uniqueIndex;not null"` // 对账日期
	TotalPayments      int64     // 总支付笔数
	TotalPaymentAmount int64     // 总支付金额（分）
	TotalRefundAmount  int64     // 总退款金额（分）
	DiscrepancyCount   int64     // 差异笔数
	Status             string    `gorm:"size:20"` // pending / matched / mismatched
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
