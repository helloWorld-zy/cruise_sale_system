package domain

import "time"

// Cruise 表示一艘邮轮的基本信息。
// 每艘邮轮隶属于一个邮轮公司，可包含多种舱房类型和设施。
type Cruise struct {
	ID                int64          `gorm:"primaryKey"`           // 主键 ID
	CompanyID         int64          `gorm:"index;not null"`       // 所属公司 ID
	Company           *CruiseCompany `gorm:"foreignKey:CompanyID"` // 关联的邮轮公司对象
	Name              string         `gorm:"size:100;not null"`    // 邮轮名称（中文）
	EnglishName       string         `gorm:"size:100"`             // 邮轮英文名称
	BuildYear         int            // 建造年份
	Tonnage           float64        // 吨位（总吨）
	PassengerCapacity int            // 最大载客量
	RoomCount         int            // 舱房总数
	Description       string         `gorm:"type:text"` // 邮轮描述
	Status            int16          `gorm:"default:1"` // 状态：1=上架，0=下架
	SortOrder         int            `gorm:"default:0"` // 排序权重，值越大越靠前
	CreatedAt         time.Time      // 创建时间
	UpdatedAt         time.Time      // 更新时间
	DeletedAt         *time.Time     `gorm:"index"` // 软删除时间
}
