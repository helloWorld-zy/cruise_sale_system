package domain

import "time"

// FacilityCategory 表示设施分类（如"餐饮"、"娱乐"、"运动"等）。
// 用于对邮轮设施进行归类管理。
type FacilityCategory struct {
	ID        int64      `gorm:"primaryKey" json:"id"`              // 主键 ID
	Name      string     `gorm:"size:100;not null" json:"name"`     // 分类名称
	Icon      string     `gorm:"size:255" json:"icon"`              // 分类图标 URL 或图标名称
	Status    int16      `gorm:"default:1" json:"status"`           // 状态：1=启用，0=停用
	SortOrder int        `gorm:"default:0" json:"sort_order"`       // 排序权重，值越大越靠前
	CreatedAt time.Time  `json:"created_at"`                        // 创建时间
	UpdatedAt time.Time  `json:"updated_at"`                        // 更新时间
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"` // 软删除时间戳
}
