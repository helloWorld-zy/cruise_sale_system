package domain

import "time"

// Route 表示邮轮航线，定义了一条从起点到终点的旅行路线。
// 一条航线可以关联多个航次（Voyage）。
type Route struct {
	ID          int64     `gorm:"primaryKey"`          // 主键 ID
	Code        string    `gorm:"size:50;uniqueIndex"` // 航线编码（全局唯一）
	Name        string    `gorm:"size:200"`            // 航线名称（如"上海-日本冲绳 5 日游"）
	Description string    `gorm:"type:text"`           // 航线描述
	Status      int16     `gorm:"default:1"`           // 状态：1=启用，0=停用
	CreatedAt   time.Time // 创建时间
	UpdatedAt   time.Time // 更新时间
}
