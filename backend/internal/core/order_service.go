package core

import (
	"context"
	"cruise_booking_system/internal/data"
	"cruise_booking_system/internal/model"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo     *data.OrderRepository
	inventoryRepo *data.InventoryRepository
	paymentSvc    *PaymentService
	inventorySvc  *InventoryService
}

func NewOrderService(
	orderRepo *data.OrderRepository,
	invRepo *data.InventoryRepository,
	paymentSvc *PaymentService,
	invSvc *InventoryService,
) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		inventoryRepo: invRepo,
		paymentSvc:    paymentSvc,
		inventorySvc:  invSvc,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID uuid.UUID, req *model.Order) (*model.Order, error) {
	err := data.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Reserve Inventory (for each item)
		// Assuming 1 item per order for MVP simplicity or loop
		for _, item := range req.Items {
			if item.CabinID != nil {
				// Lock specific cabin
				locked, err := s.inventorySvc.LockCabin(ctx, item.CabinID.String())
				if err != nil {
					return err
				}
				if !locked {
					return errors.New("cabin already locked")
				}
			}
			
			// Reserve Quota
			err := s.inventoryRepo.Reserve(ctx, tx, req.VoyageID, item.CabinTypeID, 1)
			if err != nil {
				return err
			}
		}

		// 2. Create Order
		req.UserID = userID
		req.OrderNo = uuid.New().String() // Simple unique string
		req.Status = model.OrderStatusPending
		req.ExpiresAt = time.Now().Add(15 * time.Minute)
		
		if err := s.orderRepo.Create(ctx, tx, req); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return req, nil
}
