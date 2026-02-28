package domain

import "time"

// CabinHold 表示用户对舱位库存的临时占用记录。
// 占座在指定时间后自动过期，释放锁定的库存。
type CabinHold struct {
	ID         int64     `gorm:"primaryKey"`                // 主键 ID
	CabinSKUID int64     `gorm:"column:cabin_sku_id;index"` // 占座的舱房 SKU ID
	UserID     int64     `gorm:"index"`                     // 占座用户 ID
	Qty        int       // 占用数量
	ExpiresAt  time.Time `gorm:"index"` // 占座过期时间
	CreatedAt  time.Time // 创建时间
}
