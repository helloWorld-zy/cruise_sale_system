package unit

import (
	"context"
	"testing"

	"cruise_booking_system/internal/core"
	"cruise_booking_system/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockOrderRepo
type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) Create(ctx context.Context, tx *gorm.DB, order *model.Order) error {
	return m.Called(ctx, tx, order).Error(0)
}
func (m *MockOrderRepo) ListByUserID(ctx context.Context, userID string) ([]model.Order, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Order), args.Error(1)
}
func (m *MockOrderRepo) GetByID(ctx context.Context, id string) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}
func (m *MockOrderRepo) UpdateStatus(ctx context.Context, tx *gorm.DB, id string, status model.OrderStatus) error {
	return m.Called(ctx, tx, id, status).Error(0)
}

// MockInventoryRepo
type MockInventoryRepo struct {
	mock.Mock
}

func (m *MockInventoryRepo) Reserve(ctx context.Context, tx *gorm.DB, voyageID, cabinTypeID uuid.UUID, qty int) error {
	return m.Called(ctx, tx, voyageID, cabinTypeID, qty).Error(0)
}
func (m *MockInventoryRepo) Release(ctx context.Context, tx *gorm.DB, voyageID, cabinTypeID uuid.UUID, qty int) error {
	return m.Called(ctx, tx, voyageID, cabinTypeID, qty).Error(0)
}

func TestCancelOrder_Unauthorized(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	svc := core.NewOrderService(mockOrderRepo, nil, nil, nil)

	userID := uuid.New()
	otherUserID := uuid.New()
	orderID := "order-1"

	mockOrderRepo.On("GetByID", mock.Anything, orderID).Return(&model.Order{
		BaseModel: model.BaseModel{ID: uuid.New()},
		UserID:    otherUserID,
		Status:    model.OrderStatusPending,
	}, nil)

	err := svc.CancelOrder(context.Background(), userID, orderID)

	assert.Error(t, err)
	assert.Equal(t, "unauthorized", err.Error())
}

func TestCancelOrder_InvalidState(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	svc := core.NewOrderService(mockOrderRepo, nil, nil, nil)

	userID := uuid.New()
	orderID := "order-1"

	mockOrderRepo.On("GetByID", mock.Anything, orderID).Return(&model.Order{
		BaseModel: model.BaseModel{ID: uuid.New()},
		UserID:    userID,
		Status:    model.OrderStatusCompleted,
	}, nil)

	err := svc.CancelOrder(context.Background(), userID, orderID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot cancel order")
}
