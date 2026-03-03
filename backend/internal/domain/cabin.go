package domain

import "time"

// CabinSKU 表示具体的舱房库存单元（SKU），是可售卖的最小舱房单位。
// 每个 SKU 关联一个航次（Voyage）和一种舱房类型（CabinType）。
type CabinSKU struct {
	ID          int64     `gorm:"primaryKey" json:"id"`            // 主键 ID
	VoyageID    int64     `gorm:"index" json:"voyage_id"`          // 所属航次 ID
	CabinTypeID int64     `gorm:"index" json:"cabin_type_id"`      // 所属舱房类型 ID
	Code        string    `gorm:"size:80;uniqueIndex" json:"code"` // 舱房编号，全局唯一
	Deck        string    `gorm:"size:20" json:"deck"`             // 所在甲板层
	Area        float64   `json:"area"`                            // 舱房面积（平方米）
	MaxGuests   int       `json:"max_guests"`                      // 最大入住人数
	Position    string    `gorm:"size:20" json:"position"`         // 位置：fore/mid/aft
	Orientation string    `gorm:"size:20" json:"orientation"`      // 朝向：port/starboard
	HasWindow   bool      `json:"has_window"`                      // 是否有窗
	HasBalcony  bool      `json:"has_balcony"`                     // 是否有阳台
	BedType     string    `gorm:"size:100" json:"bed_type"`        // 床型说明
	Amenities   string    `gorm:"type:text" json:"amenities"`      // 附属设施（逗号分隔）
	Grade       string    `gorm:"size:50" json:"grade"`            // 舱位等级：standard/premium/deluxe
	Status      int16     `gorm:"default:1" json:"status"`         // 状态：1=上架，0=下架
	CreatedAt   time.Time `json:"created_at"`                      // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`                      // 更新时间
}

// CabinPrice 表示舱房的日历价格，按日期和入住人数维度定价"分"为单位。
// 价格以存储，避免浮点数精度问题。
type CabinPrice struct {
	ID                    int64     `gorm:"primaryKey" json:"id"`                          // 主键 ID
	CabinSKUID            int64     `gorm:"column:cabin_sku_id;index" json:"cabin_sku_id"` // 关联的舱房 SKU ID
	Date                  time.Time `gorm:"index" json:"date"`                             // 价格生效日期
	Occupancy             int       `json:"occupancy"`                                     // 入住人数
	PriceCents            int64     `gorm:"column:price_cents" json:"price_cents"`         // 基础价格（分）
	ChildPriceCents       int64     `json:"child_price_cents"`                             // 儿童价格（分）
	SingleSupplementCents int64     `json:"single_supplement_cents"`                       // 单人补差价（分）
	PriceType             string    `gorm:"size:20;default:base" json:"price_type"`        // 价格类型：base/child/single_supplement/holiday/early_bird
	CreatedAt             time.Time `json:"created_at"`                                    // 创建时间
	UpdatedAt             time.Time `json:"updated_at"`                                    // 更新时间
}

// CabinInventory 表示舱房的库存信息，记录总量、锁定量和已售量。
// 可用库存 = Total - Locked - Sold。
type CabinInventory struct {
	ID             int64     `gorm:"primaryKey" json:"id"`                                // 主键 ID
	CabinSKUID     int64     `gorm:"column:cabin_sku_id;uniqueIndex" json:"cabin_sku_id"` // 关联的舱房 SKU ID（唯一索引）
	Total          int       `json:"total"`                                               // 库存总量
	Locked         int       `json:"locked"`                                              // 锁定量（待支付订单占用）
	Sold           int       `json:"sold"`                                                // 已售数量
	AlertThreshold int       `gorm:"default:0" json:"alert_threshold"`                    // 库存预警阈值
	UpdatedAt      time.Time `json:"updated_at"`                                          // 最后更新时间
}

// InventoryAlert 表示库存预警项。
type InventoryAlert struct {
	CabinSKUID     int64 `json:"cabin_sku_id"`    // 舱位 SKU ID
	Available      int   `json:"available"`       // 可用库存
	AlertThreshold int   `json:"alert_threshold"` // 预警阈值
}

// InventoryLog 记录库存变动的审计日志，用于追溯库存调整的原因。
type InventoryLog struct {
	ID         int64     `gorm:"primaryKey"`                // 主键 ID
	CabinSKUID int64     `gorm:"column:cabin_sku_id;index"` // 关联的舱房 SKU ID
	Change     int       // 变动量（正数为增加，负数为减少）
	Reason     string    `gorm:"size:200"` // 变动原因说明
	CreatedAt  time.Time // 变动时间
}
