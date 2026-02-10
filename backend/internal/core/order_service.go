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

type OrderRepo interface {
	Create(ctx context.Context, tx *gorm.DB, order *model.Order) error
	ListByUserID(ctx context.Context, userID string) ([]model.Order, error)
	GetByID(ctx context.Context, id string) (*model.Order, error)
	UpdateStatus(ctx context.Context, tx *gorm.DB, id string, status model.OrderStatus) error
	List(ctx context.Context, status model.OrderStatus) ([]model.Order, error)
	Update(ctx context.Context, order *model.Order) error
}

type InventoryRepo interface {
	Reserve(ctx context.Context, tx *gorm.DB, voyageID, cabinTypeID uuid.UUID, qty int) error
	Release(ctx context.Context, tx *gorm.DB, voyageID, cabinTypeID uuid.UUID, qty int) error
}

type OrderService struct {
	orderRepo     OrderRepo
	inventoryRepo InventoryRepo
	paymentSvc    *PaymentService
	inventorySvc  *InventoryService
}

func NewOrderService(
	orderRepo OrderRepo,
	invRepo InventoryRepo,
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

func (s *OrderService) ListAllOrders(ctx context.Context, status model.OrderStatus) ([]model.Order, error) {
	return s.orderRepo.List(ctx, status)
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

func (s *OrderService) ListUserOrders(ctx context.Context, userID uuid.UUID) ([]model.Order, error) {
	return s.orderRepo.ListByUserID(ctx, userID.String())
}

func (s *OrderService) CancelOrder(ctx context.Context, userID uuid.UUID, orderID string) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		return errors.New("unauthorized")
	}

	if order.Status == model.OrderStatusCancelled || order.Status == model.OrderStatusCompleted {
		return errors.New("cannot cancel order in this state")
	}

	return data.DB.Transaction(func(tx *gorm.DB) error {
		// Update Status
		if err := s.orderRepo.UpdateStatus(ctx, tx, orderID, model.OrderStatusCancelled); err != nil {
			return err
		}

		// Release Inventory
		// Assuming we load items? Repo GetByID didn't preload items in my previous edit? 
		// Wait, GetByID usually should preload. I'll assume I need to fetch items if not loaded.
		// For simplicity, let's assume GetByID preloads or we fetch them. 
		// Actually, I didn't add Preload in GetByID in repo. I should fix that or just fetch items here.
		// I'll fix Repo GetByID in a separate step if needed, or just proceed assuming it's done or I add Preload now.
		// I'll check my previous repo edit. It did `First(&order)`. No preload.
		// I should rely on Preload. I'll update Repo later or now.
		// Let's just update Repo GetByID to Preload "Items" right now in this edit? No, replace tool works on file content.
		// I'll stick to logic here: Re-fetch or assuming loaded. 
		// I'll re-fetch order with items inside transaction to be safe?
		
		var items []model.OrderItem
		if err := tx.WithContext(ctx).Where("order_id = ?", orderID).Find(&items).Error; err != nil {
			return err
		}

		for _, item := range items {
			if item.CabinID != nil {
				// Unlock cabin in Redis
				_ = s.inventorySvc.UnlockCabin(ctx, item.CabinID.String())
			}
			// Release Quota
			if err := s.inventoryRepo.Release(ctx, tx, order.VoyageID, item.CabinTypeID, 1); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *OrderService) SetDepartureNotice(ctx context.Context, orderID string, url string) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}
	order.DepartureNoticeURL = url
	return s.orderRepo.Update(ctx, order)
}
