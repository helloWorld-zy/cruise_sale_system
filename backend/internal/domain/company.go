package domain

import "time"

// CruiseCompany 表示邮轮运营公司（如皇家加勒比、地中海邮轮等）。
// 一个公司可以拥有多艘邮轮，删除公司前需确保无关联邮轮。
type CruiseCompany struct {
	ID          int64      `gorm:"primaryKey"`        // 主键 ID
	Name        string     `gorm:"size:100;not null"` // 公司名称（中文）
	EnglishName string     `gorm:"size:100"`          // 公司英文名称
	Description string     `gorm:"type:text"`         // 公司简介
	LogoURL     string     `gorm:"size:500"`          // 公司 Logo 图片地址
	Status      int16      `gorm:"default:1"`         // 状态：1=启用，0=停用
	SortOrder   int        `gorm:"default:0"`         // 排序权重，值越大越靠前
	CreatedAt   time.Time  // 创建时间
	UpdatedAt   time.Time  // 更新时间
	DeletedAt   *time.Time `gorm:"index"` // 软删除时间
}
