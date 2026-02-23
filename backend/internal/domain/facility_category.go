package domain

import "time"

// FacilityCategory 表示设施分类（如"餐饮"、"娱乐"、"运动"等）。
// 用于对邮轮设施进行归类管理。
type FacilityCategory struct {
	ID        int64     `gorm:"primaryKey"`        // 主键 ID
	Name      string    `gorm:"size:100;not null"` // 分类名称
	Icon      string    `gorm:"size:255"`          // 分类图标 URL 或图标名称
	SortOrder int       `gorm:"default:0"`         // 排序权重，值越大越靠前
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
}
