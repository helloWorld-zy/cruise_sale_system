package domain

import "time"

// CustomDestination 表示自定义目的地（私属岛屿等搜索不到的特殊目的地）。
type CustomDestination struct {
	ID          int64      `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"size:100;not null" json:"name"`                // 目的地名称
	Country     string     `gorm:"size:100;not null;default:''" json:"country"`  // 所属国家/地区
	Latitude    *float64   `json:"latitude,omitempty"`                           // 纬度
	Longitude   *float64   `json:"longitude,omitempty"`                          // 经度
	Keywords    string     `gorm:"type:text;not null;default:''" json:"keywords"` // 搜索关键词（逗号分隔）
	Description string     `gorm:"type:text;not null;default:''" json:"description"`
	Status      int16      `gorm:"default:1" json:"status"`                      // 1=启用, 0=停用
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
