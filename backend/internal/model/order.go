package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusRefunded  OrderStatus = "refunded"
	OrderStatusCompleted OrderStatus = "completed"
)

type Order struct {
	BaseModel
	OrderNo            string      `gorm:"uniqueIndex;not null" json:"order_no"`
	UserID             uuid.UUID   `gorm:"type:uuid;index" json:"user_id"`
	VoyageID           uuid.UUID   `gorm:"type:uuid;index" json:"voyage_id"`
	Status             OrderStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	TotalAmount        float64     `gorm:"type:decimal(10,2)" json:"total_amount"`
	Currency           string      `gorm:"default:'CNY'" json:"currency"`
	DepartureNoticeURL string      `json:"departure_notice_url"`
	ExpiresAt          time.Time   `json:"expires_at"`
	Items              []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	Passengers         []Passenger `gorm:"foreignKey:OrderID" json:"passengers"`
}

type OrderItem struct {
	BaseModel
	OrderID       uuid.UUID `gorm:"type:uuid;index" json:"order_id"`
	CabinID       *uuid.UUID `gorm:"type:uuid" json:"cabin_id"` // Nullable if GTY
	CabinTypeID   uuid.UUID `gorm:"type:uuid" json:"cabin_type_id"`
	PriceSnapshot float64   `gorm:"type:decimal(10,2)" json:"price_snapshot"`
}

type DocType string

const (
	DocTypePassport DocType = "passport"
	DocTypeIDCard   DocType = "id_card"
)

type Passenger struct {
	BaseModel
	OrderID   uuid.UUID `gorm:"type:uuid;index" json:"order_id"`
	NameCn    string    `json:"name_cn"`
	NameEn    string    `json:"name_en"`
	DocType   DocType   `gorm:"type:varchar(20)" json:"doc_type"`
	DocNumber string    `json:"doc_number"`
	Phone     string    `json:"phone"`
}
