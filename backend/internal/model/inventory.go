package model

import (
	"time"

	"github.com/google/uuid"
)

type VoyageStatus string

const (
	VoyageStatusScheduled VoyageStatus = "scheduled"
	VoyageStatusSailing   VoyageStatus = "sailing"
	VoyageStatusCompleted VoyageStatus = "completed"
	VoyageStatusCancelled VoyageStatus = "cancelled"
)

type Voyage struct {
	BaseModel
	CruiseID      uuid.UUID    `gorm:"type:uuid;index" json:"cruise_id"`
	RouteID       uuid.UUID    `gorm:"type:uuid;index" json:"route_id"`
	DepartureDate time.Time    `json:"departure_date"`
	ReturnDate    time.Time    `json:"return_date"`
	Status        VoyageStatus `gorm:"type:varchar(20);default:'scheduled'" json:"status"`
}

type CabinLocationType string

const (
	LocationForward CabinLocationType = "forward"
	LocationMid     CabinLocationType = "mid"
	LocationAft     CabinLocationType = "aft"
)

type Cabin struct {
	BaseModel
	CruiseID     uuid.UUID         `gorm:"type:uuid;index" json:"cruise_id"`
	CabinTypeID  uuid.UUID         `gorm:"type:uuid;index" json:"cabin_type_id"`
	Number       string            `gorm:"not null" json:"number"`
	Deck         int               `json:"deck"`
	LocationType CabinLocationType `gorm:"type:varchar(20)" json:"location_type"`
}

type Inventory struct {
	BaseModel
	VoyageID     uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_voyage_cabin_type" json:"voyage_id"`
	CabinTypeID  uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_voyage_cabin_type" json:"cabin_type_id"`
	TotalQty     int       `json:"total_qty"`
	AvailableQty int       `json:"available_qty"`
	ReservedQty  int       `json:"reserved_qty"`
	SoldQty      int       `json:"sold_qty"`
}

type PriceRule struct {
	BaseModel
	VoyageID    uuid.UUID `gorm:"type:uuid;index" json:"voyage_id"`
	CabinTypeID uuid.UUID `gorm:"type:uuid;index" json:"cabin_type_id"`
	Date        time.Time `gorm:"type:date" json:"date"`
	PriceAdult  float64   `gorm:"type:decimal(10,2)" json:"price_adult"`
	PriceChild  float64   `gorm:"type:decimal(10,2)" json:"price_child"`
	Currency    string    `gorm:"default:'CNY'" json:"currency"`
}
