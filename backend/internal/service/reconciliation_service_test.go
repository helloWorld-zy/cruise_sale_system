package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type fakePaymentRepoReconc struct {
	payments    map[string]int64
	refunds     map[string]int64
	countByDate map[string]int64
}

func newFakePaymentRepoReconc() *fakePaymentRepoReconc {
	return &fakePaymentRepoReconc{
		payments:    make(map[string]int64),
		refunds:     make(map[string]int64),
		countByDate: make(map[string]int64),
	}
}

func (r *fakePaymentRepoReconc) SumByDate(ctx context.Context, date time.Time) (int64, error) {
	return r.payments[date.Format("2006-01-02")], nil
}

func (r *fakePaymentRepoReconc) CountByDate(ctx context.Context, date time.Time) (int64, error) {
	return r.countByDate[date.Format("2006-01-02")], nil
}

func (r *fakePaymentRepoReconc) SumRefundsByDate(ctx context.Context, date time.Time) (int64, error) {
	return r.refunds[date.Format("2006-01-02")], nil
}

func TestReconciliationGenerate(t *testing.T) {
	date := time.Now().AddDate(0, 0, -1)
	dateStr := date.Format("2006-01-02")

	paymentRepo := newFakePaymentRepoReconc()
	paymentRepo.payments[dateStr] = 50000
	paymentRepo.countByDate[dateStr] = 10

	svc := NewReconciliationService(paymentRepo)
	report, err := svc.GenerateDailyReport(context.Background(), date)

	if err != nil {
		t.Fatal(err)
	}
	if report.TotalPayments == 0 {
		t.Fatal("expected payments in report")
	}
	assert.Equal(t, int64(10), report.TotalPayments)
	assert.Equal(t, int64(50000), report.TotalPaymentAmount)
}

func TestReconciliationService_GenerateDailyReport(t *testing.T) {
	date := time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC)
	dateStr := "2026-02-28"

	paymentRepo := newFakePaymentRepoReconc()
	paymentRepo.payments[dateStr] = 100000
	paymentRepo.countByDate[dateStr] = 5
	paymentRepo.refunds[dateStr] = 10000

	svc := NewReconciliationService(paymentRepo)
	report, err := svc.GenerateDailyReport(context.Background(), date)

	assert.NoError(t, err)
	assert.Equal(t, date, report.Date)
	assert.Equal(t, int64(5), report.TotalPayments)
	assert.Equal(t, int64(100000), report.TotalPaymentAmount)
	assert.Equal(t, int64(10000), report.TotalRefundAmount)
}

func TestReconciliationService_GenerateDailyReport_Empty(t *testing.T) {
	date := time.Now()

	paymentRepo := newFakePaymentRepoReconc()

	svc := NewReconciliationService(paymentRepo)
	report, err := svc.GenerateDailyReport(context.Background(), date)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), report.TotalPayments)
	assert.Equal(t, int64(0), report.TotalPaymentAmount)
}

func TestReconciliationService_GenerateDailyReportRejectsDuplicateDate(t *testing.T) {
	date := time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)
	dateStr := date.Format("2006-01-02")

	paymentRepo := newFakePaymentRepoReconc()
	paymentRepo.payments[dateStr] = 1000
	paymentRepo.countByDate[dateStr] = 1

	svc := NewReconciliationService(paymentRepo)
	_, err := svc.GenerateDailyReport(context.Background(), date)
	assert.NoError(t, err)

	_, err = svc.GenerateDailyReport(context.Background(), date)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrReconciliationReportAlreadyGenerated)
}

func TestReconciliationService_GenerateDailyReportConcurrentIdempotent(t *testing.T) {
	date := time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)
	dateStr := date.Format("2006-01-02")

	paymentRepo := newFakePaymentRepoReconc()
	paymentRepo.payments[dateStr] = 1500
	paymentRepo.countByDate[dateStr] = 2

	svc := NewReconciliationService(paymentRepo)

	const callers = 8
	var wg sync.WaitGroup
	wg.Add(callers)

	var mu sync.Mutex
	success := 0
	duplicate := 0

	for i := 0; i < callers; i++ {
		go func() {
			defer wg.Done()
			_, err := svc.GenerateDailyReport(context.Background(), date)

			mu.Lock()
			defer mu.Unlock()
			if err == nil {
				success++
				return
			}
			if assert.ErrorIs(t, err, ErrReconciliationReportAlreadyGenerated) {
				duplicate++
			}
		}()
	}

	wg.Wait()
	assert.Equal(t, 1, success)
	assert.Equal(t, callers-1, duplicate)
}
