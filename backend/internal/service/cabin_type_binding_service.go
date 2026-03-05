package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

// CabinTypeBindingService 提供舱型与邮轮绑定关系的业务能力。
type CabinTypeBindingService struct {
	repo domain.CabinTypeBindingRepository
}

func NewCabinTypeBindingService(repo domain.CabinTypeBindingRepository) *CabinTypeBindingService {
	return &CabinTypeBindingService{repo: repo}
}

func (s *CabinTypeBindingService) ReplaceCruiseBindings(ctx context.Context, cabinTypeID int64, cruiseIDs []int64) error {
	return s.repo.ReplaceCruiseBindings(ctx, cabinTypeID, cruiseIDs)
}

func (s *CabinTypeBindingService) ListCruiseIDsByCabinType(ctx context.Context, cabinTypeID int64) ([]int64, error) {
	return s.repo.ListCruiseIDsByCabinType(ctx, cabinTypeID)
}

func (s *CabinTypeBindingService) ListCabinTypeIDsByCruise(ctx context.Context, cruiseID int64) ([]int64, error) {
	return s.repo.ListCabinTypeIDsByCruise(ctx, cruiseID)
}

func (s *CabinTypeBindingService) HasCabinTypesByCruise(ctx context.Context, cruiseID int64) (bool, error) {
	return s.repo.HasCabinTypesByCruise(ctx, cruiseID)
}
