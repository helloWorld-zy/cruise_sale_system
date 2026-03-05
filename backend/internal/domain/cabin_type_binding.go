package domain

import "time"

// CabinTypeCruiseBinding 表示舱型与邮轮的绑定关系。
type CabinTypeCruiseBinding struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	CabinTypeID int64     `gorm:"index;not null" json:"cabin_type_id"`
	CruiseID    int64     `gorm:"index;not null" json:"cruise_id"`
	CreatedAt   time.Time `json:"created_at"`
}
