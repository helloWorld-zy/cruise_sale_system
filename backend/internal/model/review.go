package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ReviewStatus string

const (
	ReviewStatusPending  ReviewStatus = "pending"
	ReviewStatusApproved ReviewStatus = "approved"
	ReviewStatusRejected ReviewStatus = "rejected"
)

type Review struct {
	BaseModel
	UserID   uuid.UUID      `gorm:"type:uuid;index" json:"user_id"`
	CruiseID uuid.UUID      `gorm:"type:uuid;index" json:"cruise_id"`
	OrderID  uuid.UUID      `gorm:"type:uuid;index" json:"order_id"`
	Rating   int            `json:"rating"` // 1-5
	Comment  string         `gorm:"type:text" json:"comment"`
	Media    datatypes.JSON `json:"media"`
	Status   ReviewStatus   `gorm:"type:varchar(20);default:'pending'" json:"status"`
}
