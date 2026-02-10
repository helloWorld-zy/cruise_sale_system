package model

import (
	"github.com/google/uuid"
)

type Facility struct {
	BaseModel
	CruiseID   uuid.UUID `gorm:"type:uuid;index" json:"cruise_id"`
	CategoryID uuid.UUID `gorm:"type:uuid;index" json:"category_id"`
	Name       string    `gorm:"not null" json:"name"`
	Location   string    `json:"location"`
	IsPaid     bool      `json:"is_paid"`
}
