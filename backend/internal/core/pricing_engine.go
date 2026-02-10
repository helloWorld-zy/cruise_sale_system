package core

import (
	"context"
	"cruise_booking_system/internal/model"
	"time"
)

type PricingEngine struct {
	inventoryRepo InventoryRepo
	// ruleRepo RuleRepo
}

func NewPricingEngine(invRepo InventoryRepo) *PricingEngine {
	return &PricingEngine{
		inventoryRepo: invRepo,
	}
}

func (e *PricingEngine) CalculatePrice(ctx context.Context, voyageID, cabinTypeID string, basePrice float64) (float64, error) {
	// Logic:
	// 1. Get Inventory
	// 2. Calculate Load Factor
	// 3. Apply Rules (e.g. if load > 80% -> +10%)
	
	// Mock logic for MVP
	// inv, _ := e.inventoryRepo.Get(ctx, voyageID, cabinTypeID)
	// if inv.AvailableQty < 5 { return basePrice * 1.2, nil }
	
	return basePrice, nil
}

// Scheduled job would call this
func (e *PricingEngine) RunDailyAdjustment(ctx context.Context) error {
	// Iterate all active voyages...
	// Update PriceRule...
	return nil
}
