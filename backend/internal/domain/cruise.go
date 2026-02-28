package domain

import "time"

// Cruise 表示一艘邮轮的基本信息。
// 每艘邮轮隶属于一个邮轮公司，可包含多种舱房类型和设施。
type Cruise struct {
	ID                int64          `gorm:"primaryKey" json:"id"`                          // 主键 ID
	CompanyID         int64          `gorm:"index;not null" json:"company_id"`              // 所属公司 ID
	Company           *CruiseCompany `gorm:"foreignKey:CompanyID" json:"company,omitempty"` // 关联的邮轮公司对象
	Name              string         `gorm:"size:100;not null" json:"name"`                 // 邮轮名称（中文）
	EnglishName       string         `gorm:"size:100" json:"english_name"`                  // 邮轮英文名称
	BuildYear         int            `json:"build_year"`                                    // 建造年份
	Tonnage           float64        `json:"tonnage"`                                       // 吨位（总吨）
	PassengerCapacity int            `json:"passenger_capacity"`                            // 最大载客量
	RoomCount         int            `json:"room_count"`                                    // 舱房总数
	Description       string         `gorm:"type:text" json:"description"`                  // 邮轮描述
	Status            int16          `gorm:"default:1" json:"status"`                       // 状态：1=上架，0=下架
	SortOrder         int            `gorm:"default:0" json:"sort_order"`                   // 排序权重，值越大越靠前
	CreatedAt         time.Time      `json:"created_at"`                                    // 创建时间
	UpdatedAt         time.Time      `json:"updated_at"`                                    // 更新时间
	DeletedAt         *time.Time     `gorm:"index" json:"deleted_at,omitempty"`             // 软删除时间
}
