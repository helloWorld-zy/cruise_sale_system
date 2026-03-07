package domain

import "time"

// Reconciliation 表示每日财务对账记录。
// 记录每日的支付笔数、支付金额、退款金额以及对账状态。
type Reconciliation struct {
	ID                 int64     `gorm:"primaryKey"`           // 主键 ID
	Date               time.Time `gorm:"uniqueIndex;not null"` // 对账日期
	TotalPayments      int64     // 总支付笔数
	TotalPaymentAmount int64     // 总支付金额（单位：分）
	TotalRefundAmount  int64     // 总退款金额（单位：分）
	DiscrepancyCount   int64     // 差异笔数（支付与退款不一致的数量）
	Status             string    `gorm:"size:20"` // 对账状态：pending（待对账）/ matched（已匹配）/ mismatched（有差异）
	CreatedAt          time.Time // 创建时间
	UpdatedAt          time.Time // 更新时间
}
