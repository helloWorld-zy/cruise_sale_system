package domain

import "time"

// Route 表示邮轮航线，定义了一条从起点到终点的旅行路线。
// 一条航线可以关联多个航次（Voyage）。
type Route struct {
	ID            int64     `gorm:"primaryKey" json:"id"`            // 主键 ID
	Code          string    `gorm:"size:50;uniqueIndex" json:"code"` // 航线编码（全局唯一）
	Name          string    `gorm:"size:200" json:"name"`            // 航线名称（如"上海-日本冲绳 5 日游"）
	DeparturePort string    `gorm:"size:100" json:"departure_port"`  // 出发港口
	ArrivalPort   string    `gorm:"size:100" json:"arrival_port"`    // 到达港口
	Stops         string    `gorm:"type:text" json:"stops"`          // 途经停靠港（逗号分隔）
	Description   string    `gorm:"type:text" json:"description"`    // 航线描述
	Status        int16     `gorm:"default:1" json:"status"`         // 状态：1=启用，0=停用
	SortOrder     int       `gorm:"default:0" json:"sort_order"`     // 排序权重
	CreatedAt     time.Time `json:"created_at"`                      // 创建时间
	UpdatedAt     time.Time `json:"updated_at"`                      // 更新时间
}
