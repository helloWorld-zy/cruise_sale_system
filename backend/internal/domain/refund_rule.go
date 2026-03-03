package domain

// RefundRule 定义阶梯退款规则。
type RefundRule struct {
	ID         int64 `gorm:"primaryKey"`
	MinDays    int   // 提前天数下限（包含）
	MaxDays    int   // 提前天数上限（不包含）
	RefundRate int   // 退款百分比（0-100）
}
