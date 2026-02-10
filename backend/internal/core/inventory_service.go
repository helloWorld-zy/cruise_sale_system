package core

import (
	"context"
	"fmt"
	"time"

	"cruise_booking_system/internal/data"
)

type InventoryService struct {
	repo *data.InventoryRepository
}

func NewInventoryService(repo *data.InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) LockCabin(ctx context.Context, cabinID string) (bool, error) {
	// Redis lock for specific cabin number
	key := fmt.Sprintf("lock:cabin:%s", cabinID)
	// 15 min lock
	success, err := data.RDB.SetNX(ctx, key, "locked", 15*time.Minute).Result()
	return success, err
}

func (s *InventoryService) UnlockCabin(ctx context.Context, cabinID string) error {
	key := fmt.Sprintf("lock:cabin:%s", cabinID)
	return data.RDB.Del(ctx, key).Err()
}
