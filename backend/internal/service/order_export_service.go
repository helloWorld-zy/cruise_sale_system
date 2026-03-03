package service

import (
	"context"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

type OrderFilter struct {
	Status    string
	Phone     string
	RouteID   int64
	VoyageID  int64
	StartDate *time.Time
	EndDate   *time.Time
	BookingNo string
}

type OrderRepositoryFilter interface {
	ListWithFilter(ctx context.Context, filter OrderFilter) ([]domain.Booking, error)
}

type OrderExportService struct {
	repo OrderRepositoryFilter
}

func NewOrderExportService(repo OrderRepositoryFilter) *OrderExportService {
	return &OrderExportService{repo: repo}
}

func (s *OrderExportService) ExportToExcel(ctx context.Context, filter OrderFilter) ([]byte, error) {
	orders, err := s.repo.ListWithFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	return generateExcelBytes(orders), nil
}

func generateExcelBytes(orders []domain.Booking) []byte {
	return []byte("Excel content placeholder")
}
