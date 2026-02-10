package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type CruiseStatus string

const (
	CruiseStatusActive      CruiseStatus = "active"
	CruiseStatusMaintenance CruiseStatus = "maintenance"
	CruiseStatusRetired     CruiseStatus = "retired"
)

type Cruise struct {
	BaseModel
	NameEn      string         `gorm:"not null" json:"name_en"`
	NameCn      string         `gorm:"not null" json:"name_cn"`
	Code        string         `gorm:"uniqueIndex;not null" json:"code"`
	CompanyID   uuid.UUID      `gorm:"type:uuid;index" json:"company_id"`
	Tonnage     int            `json:"tonnage"`
	Capacity    int            `json:"capacity"`
	Decks       int            `json:"decks"`
	Status      CruiseStatus   `gorm:"type:varchar(20);default:'active'" json:"status"`
	Gallery     datatypes.JSON `json:"gallery"` // Array of image URLs
	Description string         `gorm:"type:text" json:"description"`
}

type CabinType struct {
	BaseModel
	CruiseID   uuid.UUID      `gorm:"type:uuid;index" json:"cruise_id"`
	Name       string         `gorm:"not null" json:"name"`
	Code       string         `json:"code"`
	BaseArea   float64        `gorm:"type:decimal(10,2)" json:"base_area"`
	Capacity   int            `json:"capacity"`
	Facilities datatypes.JSON `json:"facilities"` // List of amenities
	Images     datatypes.JSON `json:"images"`
}
