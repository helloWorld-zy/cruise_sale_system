package domain

import "time"

// CabinType 表示邮轮上的舱房类型（如内舱、海景舱、阳台舱、套房等）。
// 每种舱房类型隶属于一艘邮轮，包含容量、面积、甲板等基本属性。
type CabinType struct {
	ID           int64      `gorm:"primaryKey" json:"id"`              // 主键 ID
	CruiseID     int64      `gorm:"index;not null" json:"cruise_id"`   // 所属邮轮 ID
	Name         string     `gorm:"size:100;not null" json:"name"`     // 舱房类型名称（中文）
	EnglishName  string     `gorm:"size:100" json:"english_name"`      // 舱房类型英文名称
	Code         string     `gorm:"size:50" json:"code"`               // 舱房类型代码
	AreaMin      float64    `json:"area_min"`                          // 最小面积（平方米）
	AreaMax      float64    `json:"area_max"`                          // 最大面积（平方米）
	Area         float64    `json:"area"`                              // 兼容旧字段：舱房面积（平方米）
	Capacity     int        `gorm:"default:2" json:"capacity"`         // 标准入住人数
	MaxCapacity  int        `json:"max_capacity"`                      // 最大入住人数
	BedType      string     `gorm:"size:200" json:"bed_type"`          // 床型说明
	Tags         string     `gorm:"size:500" json:"tags"`              // 特色标签（逗号分隔）
	Amenities    string     `gorm:"type:text" json:"amenities"`        // 设施清单（逗号分隔）
	FloorPlanURL string     `gorm:"size:500" json:"floor_plan_url"`    // 平面图 URL
	Deck         string     `gorm:"size:50" json:"deck"`               // 所在甲板层
	Description  string     `gorm:"type:text" json:"description"`      // 舱房类型描述
	Status       int16      `gorm:"default:1" json:"status"`           // 状态：1=启用，0=停用
	SortOrder    int        `gorm:"default:0" json:"sort_order"`       // 排序权重，值越大越靠前
	CreatedAt    time.Time  `json:"created_at"`                        // 创建时间
	UpdatedAt    time.Time  `json:"updated_at"`                        // 更新时间
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"` // 软删除时间
}
