package domain

// RefundRule 定义阶梯退款规则，用于根据提前退款的天数计算退款金额。
// 例如：提前7天退款100%，提前3天退款50%等。
type RefundRule struct {
	ID         int64 `gorm:"primaryKey"` // 主键 ID
	MinDays    int   // 提前天数下限（包含），如 7 表示提前7天及以上
	MaxDays    int   // 提前天数上限（不包含），如 3 表示不足3天
	RefundRate int   // 退款百分比（0-100），如 100 表示全额退款，50 表示退款50%
}
