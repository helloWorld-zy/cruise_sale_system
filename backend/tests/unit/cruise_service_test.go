package unit

import (
	"context"
	"testing"

	"cruise_booking_system/internal/core"
	"cruise_booking_system/internal/data"
	"cruise_booking_system/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCruiseRepo is a mock implementation of core.CruiseRepo
type MockCruiseRepo struct {
	mock.Mock
}

func (m *MockCruiseRepo) List(ctx context.Context, filter data.CruiseFilter) ([]model.Cruise, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]model.Cruise), args.Error(1)
}

func (m *MockCruiseRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Cruise, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Cruise), args.Error(1)
}

func (m *MockCruiseRepo) GetCabinTypes(ctx context.Context, cruiseID uuid.UUID) ([]model.CabinType, error) {
	args := m.Called(ctx, cruiseID)
	return args.Get(0).([]model.CabinType), args.Error(1)
}

func TestListCruises(t *testing.T) {
	mockRepo := new(MockCruiseRepo)
	service := core.NewCruiseService(mockRepo)

	expectedCruises := []model.Cruise{
		{NameEn: "Test Cruise"},
	}

	mockRepo.On("List", mock.Anything, mock.Anything).Return(expectedCruises, nil)

	cruises, err := service.ListCruises(context.Background(), "", "")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(cruises))
	assert.Equal(t, "Test Cruise", cruises[0].NameEn)
	mockRepo.AssertExpectations(t)
}

func TestGetCruiseDetail(t *testing.T) {
	mockRepo := new(MockCruiseRepo)
	service := core.NewCruiseService(mockRepo)

	id := uuid.New()
	cruise := &model.Cruise{NameEn: "Detail Cruise"}
	cabinTypes := []model.CabinType{{Name: "Balcony"}}

	mockRepo.On("GetByID", mock.Anything, id).Return(cruise, nil)
	mockRepo.On("GetCabinTypes", mock.Anything, id).Return(cabinTypes, nil)

	detail, err := service.GetCruiseDetail(context.Background(), id.String())

	assert.NoError(t, err)
	assert.Equal(t, "Detail Cruise", detail.NameEn)
	assert.Equal(t, 1, len(detail.CabinTypes))
	mockRepo.AssertExpectations(t)
}
