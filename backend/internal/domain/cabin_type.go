package domain

import "time"

// CabinType 表示邮轮上的舱房类型（如内舱、海景舱、阳台舱、套房等）。
// 每种舱房类型隶属于一艘邮轮，包含容量、面积、甲板等基本属性。
type CabinType struct {
	ID          int64      `gorm:"primaryKey"`        // 主键 ID
	CruiseID    int64      `gorm:"index;not null"`    // 所属邮轮 ID
	Name        string     `gorm:"size:100;not null"` // 舱房类型名称（中文）
	EnglishName string     `gorm:"size:100"`          // 舱房类型英文名称
	Capacity    int        `gorm:"default:2"`         // 默认容纳人数
	Area        float64    // 舱房面积（平方米）
	Deck        string     `gorm:"size:50"`   // 所在甲板层
	Description string     `gorm:"type:text"` // 舱房类型描述
	Status      int16      `gorm:"default:1"` // 状态：1=启用，0=停用
	SortOrder   int        `gorm:"default:0"` // 排序权重，值越大越靠前
	CreatedAt   time.Time  // 创建时间
	UpdatedAt   time.Time  // 更新时间
	DeletedAt   *time.Time `gorm:"index"` // 软删除时间
}
