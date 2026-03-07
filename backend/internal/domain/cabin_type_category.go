package domain

import "time"

// CabinTypeCategory 表示舱型大类字典（如内舱、海景、阳台、套房）。
// 用于对舱房类型进行一级分类，方便用户筛选和展示。
type CabinTypeCategory struct {
	ID        int64      `gorm:"primaryKey" json:"id"`                     // 主键 ID
	Name      string     `gorm:"size:64;not null" json:"name"`             // 分类名称（如"内舱"、"海景舱"）
	Code      string     `gorm:"size:32;not null;uniqueIndex" json:"code"` // 分类代码（全局唯一）
	Status    int16      `gorm:"default:1" json:"status"`                  // 状态：1=启用，0=停用
	SortOrder int        `gorm:"default:0" json:"sort_order"`              // 排序权重，值越大越靠前
	CreatedAt time.Time  `json:"created_at"`                               // 创建时间
	UpdatedAt time.Time  `json:"updated_at"`                               // 更新时间
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`        // 软删除时间
}
