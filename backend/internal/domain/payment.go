package domain

import "time"

// Payment 表示支付记录实体。
type Payment struct {
	ID          int64     `gorm:"primaryKey"`            // 主键 ID
	OrderID     int64     `gorm:"index;column:order_id"` // 订单 ID
	Provider    string    `gorm:"size:20"`               // 支付提供商（如：alipay, wechat）
	TradeNo     string    `gorm:"size:100"`              // 交易流水号
	AmountCents int64     // 支付金额（单位：分）
	Status      string    `gorm:"size:20"` // 支付状态（如：pending, success, failed）
	CreatedAt   time.Time // 创建时间
	UpdatedAt   time.Time // 更新时间
}
