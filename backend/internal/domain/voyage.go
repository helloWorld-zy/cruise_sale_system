package domain

import "time"

// Voyage 表示一个具体的航次（某艘邮轮在某条航线上的一次出发）。
// 航次是舱房 SKU 和价格日历的核心关联维度。
type Voyage struct {
	ID         int64     `gorm:"primaryKey"`          // 主键 ID
	RouteID    int64     `gorm:"index"`               // 所属航线 ID
	CruiseID   int64     `gorm:"index"`               // 执行邮轮 ID
	Code       string    `gorm:"size:50;uniqueIndex"` // 航次编码（全局唯一）
	DepartDate time.Time // 出发日期
	ReturnDate time.Time // 返航日期
	Status     int16     `gorm:"default:1"` // 状态：1=开放预订，0=关闭
	CreatedAt  time.Time // 创建时间
	UpdatedAt  time.Time // 更新时间
}
