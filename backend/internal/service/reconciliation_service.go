package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

// PaymentReconciliationReader 定义支付对账数据查询接口。
type PaymentReconciliationReader interface {
	SumByDate(ctx context.Context, date time.Time) (int64, error)        // 按日期统计支付金额
	CountByDate(ctx context.Context, date time.Time) (int64, error)      // 按日期统计支付笔数
	SumRefundsByDate(ctx context.Context, date time.Time) (int64, error) // 按日期统计退款金额
}

// ReconciliationService 提供每日对账报表生成服务。
type ReconciliationService struct {
	paymentRepo PaymentReconciliationReader // 支付数据仓储
	mu          sync.Mutex                  // 互斥锁，防止重复生成
	generated   map[string]struct{}         // 已生成的报表日期缓存
}

var (
	// ErrReconciliationReportAlreadyGenerated 表示指定日期的对账报表已生成。
	ErrReconciliationReportAlreadyGenerated = errors.New("reconciliation report already generated for date")
)

// NewReconciliationService 创建对账服务实例。
func NewReconciliationService(paymentRepo PaymentReconciliationReader) *ReconciliationService {
	return &ReconciliationService{
		paymentRepo: paymentRepo,
		generated:   make(map[string]struct{}),
	}
}

// GenerateDailyReport 生成指定日期的每日对账报表。
// 流程：检查是否已生成 → 查询支付/退款数据 → 生成对账报表。
// 注意：此方法仅在内存中记录已生成状态，重启服务后会重置。
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
