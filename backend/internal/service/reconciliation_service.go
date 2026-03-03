package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

type PaymentReconciliationReader interface {
	SumByDate(ctx context.Context, date time.Time) (int64, error)
	CountByDate(ctx context.Context, date time.Time) (int64, error)
	SumRefundsByDate(ctx context.Context, date time.Time) (int64, error)
}

type ReconciliationService struct {
	paymentRepo PaymentReconciliationReader
	mu          sync.Mutex
	generated   map[string]struct{}
}

var (
	// ErrReconciliationReportAlreadyGenerated indicates a daily report has already been generated for the date.
	ErrReconciliationReportAlreadyGenerated = errors.New("reconciliation report already generated for date")
)

func NewReconciliationService(paymentRepo PaymentReconciliationReader) *ReconciliationService {
	return &ReconciliationService{
		paymentRepo: paymentRepo,
		generated:   make(map[string]struct{}),
	}
}

func (s *ReconciliationService) GenerateDailyReport(ctx context.Context, date time.Time) (*domain.Reconciliation, error) {
	normalizedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dateKey := normalizedDate.Format("2006-01-02")

	s.mu.Lock()
	_, exists := s.generated[dateKey]
	s.mu.Unlock()
	if exists {
		return nil, ErrReconciliationReportAlreadyGenerated
	}

	totalPayments, err := s.paymentRepo.CountByDate(ctx, normalizedDate)
	if err != nil {
		return nil, fmt.Errorf("count payments: %w", err)
	}

	totalPaymentAmount, err := s.paymentRepo.SumByDate(ctx, normalizedDate)
	if err != nil {
		return nil, fmt.Errorf("sum payments: %w", err)
	}

	totalRefundAmount, err := s.paymentRepo.SumRefundsByDate(ctx, normalizedDate)
	if err != nil {
		return nil, fmt.Errorf("sum refunds: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, duplicated := s.generated[dateKey]; duplicated {
		return nil, ErrReconciliationReportAlreadyGenerated
	}
	s.generated[dateKey] = struct{}{}

	report := &domain.Reconciliation{
		Date:               normalizedDate,
		TotalPayments:      totalPayments,
		TotalPaymentAmount: totalPaymentAmount,
		TotalRefundAmount:  totalRefundAmount,
		DiscrepancyCount:   0,
		Status:             "matched",
	}

	return report, nil
}
