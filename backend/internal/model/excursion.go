package model

import (
	"github.com/google/uuid"
)

type Excursion struct {
	BaseModel
	CruiseID    uuid.UUID `gorm:"type:uuid;index" json:"cruise_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Price       float64   `gorm:"type:decimal(10,2)" json:"price"`
	Duration    string    `json:"duration"` // e.g. "4h"
}
