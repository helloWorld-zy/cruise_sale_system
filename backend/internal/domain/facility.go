package domain

import "time"

// Facility 表示邮轮上的一项娱乐或服务设施（如泳池、健身房、餐厅等）。
// 每个设施归属于一个设施分类和一艘邮轮。
type Facility struct {
	ID          int64             `gorm:"primaryKey"`            // 主键 ID
	CategoryID  int64             `gorm:"index;not null"`        // 所属设施分类 ID
	Category    *FacilityCategory `gorm:"foreignKey:CategoryID"` // 关联的设施分类对象
	CruiseID    int64             `gorm:"index;not null"`        // 所属邮轮 ID
	Name        string            `gorm:"size:100;not null"`     // 设施名称（中文）
	EnglishName string            `gorm:"size:100"`              // 设施英文名称
	Location    string            `gorm:"size:100"`              // 设施位置（如"6层甲板"）
	Description string            `gorm:"type:text"`             // 设施描述
	Status      int16             `gorm:"default:1"`             // 状态：1=开放，0=关闭
	SortOrder   int               `gorm:"default:0"`             // 排序权重，值越大越靠前
	CreatedAt   time.Time         // 创建时间
	UpdatedAt   time.Time         // 更新时间
	DeletedAt   *time.Time        `gorm:"index"` // 软删除时间
}
