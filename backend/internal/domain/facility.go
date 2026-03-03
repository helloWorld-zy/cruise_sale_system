package domain

import "time"

// Facility 表示邮轮上的一项娱乐或服务设施（如泳池、健身房、餐厅等）。
// 每个设施归属于一个设施分类和一艘邮轮。
type Facility struct {
	ID             int64             `gorm:"primaryKey" json:"id"`                            // 主键 ID
	CategoryID     int64             `gorm:"index;not null" json:"category_id"`               // 所属设施分类 ID
	Category       *FacilityCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"` // 关联的设施分类对象
	CruiseID       int64             `gorm:"index;not null" json:"cruise_id"`                 // 所属邮轮 ID
	Name           string            `gorm:"size:100;not null" json:"name"`                   // 设施名称（中文）
	EnglishName    string            `gorm:"size:100" json:"english_name"`                    // 设施英文名称
	Location       string            `gorm:"size:100" json:"location"`                        // 设施位置（如"6层甲板"）
	OpenHours      string            `gorm:"size:100" json:"open_hours"`                      // 开放时间（如 08:00-22:00）
	ExtraCharge    bool              `json:"extra_charge"`                                    // 是否额外收费
	ChargePriceTip string            `gorm:"size:200" json:"charge_price_tip"`                // 收费说明或参考价格提示
	TargetAudience string            `gorm:"size:200" json:"target_audience"`                 // 适合人群（逗号分隔）
	Description    string            `gorm:"type:text" json:"description"`                    // 设施描述
	Status         int16             `gorm:"default:1" json:"status"`                         // 状态：1=开放，0=关闭
	SortOrder      int               `gorm:"default:0" json:"sort_order"`                     // 排序权重，值越大越靠前
	CreatedAt      time.Time         `json:"created_at"`                                      // 创建时间
	UpdatedAt      time.Time         `json:"updated_at"`                                      // 更新时间
	DeletedAt      *time.Time        `gorm:"index" json:"deleted_at,omitempty"`               // 软删除时间
}
