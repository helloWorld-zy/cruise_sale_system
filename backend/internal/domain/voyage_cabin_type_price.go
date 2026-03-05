package domain

import "time"

// VoyageCabinTypePriceVersion 记录航次舱型价格与库存的历史版本。
type VoyageCabinTypePriceVersion struct {
	ID                   int64     `gorm:"primaryKey" json:"id"`
	VoyageID             int64     `gorm:"index;not null" json:"voyage_id"`
	CabinTypeID          int64     `gorm:"index;not null" json:"cabin_type_id"`
	InventoryTotal       int       `json:"inventory_total"`
	SettlementPriceCents int64     `json:"settlement_price_cents"`
	SalePriceCents       int64     `json:"sale_price_cents"`
	EffectiveAt          time.Time `json:"effective_at"`
	CreatedBy            *int64    `json:"created_by,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
}

// VoyageCabinTypeCurrent 表示航次舱型当前对外生效价格快照。
type VoyageCabinTypeCurrent struct {
	VoyageID             int64     `gorm:"primaryKey" json:"voyage_id"`
	CabinTypeID          int64     `gorm:"primaryKey" json:"cabin_type_id"`
	InventoryTotal       int       `json:"inventory_total"`
	SettlementPriceCents int64     `json:"settlement_price_cents"`
	SalePriceCents       int64     `json:"sale_price_cents"`
	EffectiveAt          time.Time `json:"effective_at"`
	VersionID            int64     `json:"version_id"`
	UpdatedAt            time.Time `json:"updated_at"`
}
