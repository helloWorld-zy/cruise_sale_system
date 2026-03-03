package service

import (
	"context"
	"fmt"
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
}

func NewReconciliationService(paymentRepo PaymentReconciliationReader) *ReconciliationService {
	return &ReconciliationService{paymentRepo: paymentRepo}
}

func (s *ReconciliationService) GenerateDailyReport(ctx context.Context, date time.Time) (*domain.Reconciliation, error) {
	totalPayments, err := s.paymentRepo.CountByDate(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("count payments: %w", err)
	}

	totalPaymentAmount, err := s.paymentRepo.SumByDate(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("sum payments: %w", err)
	}

	totalRefundAmount, err := s.paymentRepo.SumRefundsByDate(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("sum refunds: %w", err)
	}

	report := &domain.Reconciliation{
		Date:               date,
		TotalPayments:      totalPayments,
		TotalPaymentAmount: totalPaymentAmount,
		TotalRefundAmount:  totalRefundAmount,
		DiscrepancyCount:   0,
		Status:             "matched",
	}

	return report, nil
}
