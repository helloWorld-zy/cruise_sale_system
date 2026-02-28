package service

import (
	"context"
	"errors"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// BookingRepo 定义预订写入与事务边界能力。
type BookingRepo interface {
	Create(ctx context.Context, b *domain.Booking) error
	InTx(fn func(tx *gorm.DB, create func(b *domain.Booking) error) error) error
}

// PriceService 定义舱位价格查询能力。
type PriceService interface {
	FindPrice(ctx context.Context, skuID int64, date time.Time, occupancy int) (int64, bool, error)
}

// HoldService 定义库存占用能力。
type HoldService interface {
	HoldWithTx(tx *gorm.DB, skuID int64, userID int64, qty int) bool
}

// BookingService 负责预订创建流程编排。
type BookingService struct {
	repo  BookingRepo
	price PriceService
	hold  HoldService
}

// NewBookingService 创建预订服务实例。
func NewBookingService(repo BookingRepo, price PriceService, hold HoldService) *BookingService {
	return &BookingService{repo: repo, price: price, hold: hold}
}

// Create 创建预订并在事务内完成库存占用与金额计算，返回已创建订单。
func (s *BookingService) Create(ctx context.Context, userID, voyageID, skuID int64, guests int) (*domain.Booking, error) {
	if s.repo == nil || s.price == nil || s.hold == nil {
		return nil, errors.New("booking dependencies not ready")
	}

	var created domain.Booking

	err := s.repo.InTx(func(tx *gorm.DB, create func(b *domain.Booking) error) error {
		if !s.hold.HoldWithTx(tx, skuID, userID, 1) {
			return errors.New("cannot hold inventory")
		}

		price, _, err := s.price.FindPrice(ctx, skuID, time.Now(), guests)
		if err != nil {
			return err
		}

		created = domain.Booking{UserID: userID, VoyageID: voyageID, CabinSKUID: skuID, Status: "created", TotalCents: price}
		return create(&created)
	})
	if err != nil {
		return nil, err
	}

	return &created, nil
}
