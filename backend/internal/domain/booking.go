package domain

import "time"

const (
	OrderStatusCreated        = "created"
	OrderStatusPendingPayment = "pending_payment"
	OrderStatusPaid           = "paid"
	OrderStatusConfirmed      = "confirmed"
	OrderStatusPendingTravel  = "pending_travel"
	OrderStatusTraveling      = "traveling"
	OrderStatusCompleted      = "completed"
	OrderStatusCancelled      = "cancelled"
	OrderStatusRefunding      = "refunding"
	OrderStatusRefunded       = "refunded"
)

var validTransitions = map[string][]string{
	OrderStatusCreated:        {OrderStatusPendingPayment, OrderStatusCancelled},
	OrderStatusPendingPayment: {OrderStatusPaid, OrderStatusCancelled},
	OrderStatusPaid:           {OrderStatusConfirmed, OrderStatusRefunding},
	OrderStatusConfirmed:      {OrderStatusPendingTravel, OrderStatusRefunding},
	OrderStatusPendingTravel:  {OrderStatusTraveling, OrderStatusRefunding},
	OrderStatusTraveling:      {OrderStatusCompleted},
	OrderStatusRefunding:      {OrderStatusRefunded},
}

func (b *Booking) CanTransitionTo(targetStatus string) bool {
	allowed, exists := validTransitions[b.Status]
	if !exists {
		return false
	}
	for _, s := range allowed {
		if s == targetStatus {
			return true
		}
	}
	return false
}

// Booking 表示用户对航次舱位的预订订单。
type Booking struct {
	ID         int64     `gorm:"primaryKey" json:"id"`                    // 主键 ID
	UserID     int64     `gorm:"index" json:"user_id"`                    // 下单用户 ID
	VoyageID   int64     `json:"voyage_id"`                               // 所属航次 ID
	CabinSKUID int64     `gorm:"column:cabin_sku_id" json:"cabin_sku_id"` // 预订的舱房 SKU ID
	Status     string    `gorm:"size:30;default:created" json:"status"`   // 订单状态
	TotalCents int64     `json:"total_cents"`                             // 订单总金额（单位：分）
	PaidCents  int64     `json:"paid_cents"`                              // 已支付金额（单位：分）
	CreatedAt  time.Time `json:"created_at"`                              // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`                              // 更新时间
}

// OrderStatusLog 记录订单状态变更日志。
type OrderStatusLog struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	OrderID    int64     `gorm:"index" json:"order_id"`
	FromStatus string    `gorm:"size:30" json:"from_status"`
	ToStatus   string    `gorm:"size:30" json:"to_status"`
	OperatorID int64     `json:"operator_id"`
	Remark     string    `gorm:"type:text" json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
}
