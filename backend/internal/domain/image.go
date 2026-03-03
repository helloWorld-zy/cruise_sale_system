package domain

import "time"

// Image 表示系统中的图片资源，通过 EntityType 和 EntityID 实现多态关联。
// 可关联邮轮、舱房、设施等多种实体类型。
type Image struct {
	ID         int64     `gorm:"primaryKey" json:"id"`                                        // 主键 ID
	EntityType string    `gorm:"size:50;index:idx_images_entity;not null" json:"entity_type"` // 关联实体类型（如 "cruise"、"cabin"）
	EntityID   int64     `gorm:"index:idx_images_entity;not null" json:"entity_id"`           // 关联实体 ID
	URL        string    `gorm:"size:500;not null" json:"url"`                                // 图片 URL 地址
	SortOrder  int       `gorm:"default:0" json:"sort_order"`                                 // 排序权重，值越大越靠前
	IsPrimary  bool      `gorm:"default:false" json:"is_primary"`                             // 是否为主图
	CreatedAt  time.Time `json:"created_at"`                                                  // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`                                                  // 更新时间
}
