package service

import (
	"context"
	"errors"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubPayRepoForRefund 为退款测试实现 domain.PaymentRepository。
type stubPayRepoForRefund struct {
	payment *domain.Payment
	findErr error
}

func (r *stubPayRepoForRefund) Create(_ context.Context, _ *domain.Payment) error { return nil }
func (r *stubPayRepoForRefund) FindByTradeNo(_ context.Context, _ string) (*domain.Payment, error) {
	return r.payment, r.findErr
}
func (r *stubPayRepoForRefund) FindByID(_ context.Context, id int64) (*domain.Payment, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	if r.payment != nil && r.payment.ID == id {
		return r.payment, nil
	}
	return nil, errors.New("not found")
}
func (r *stubPayRepoForRefund) UpdateStatus(_ context.Context, _ int64, _ string) error { return nil }

// stubRefundRepo 实现 domain.RefundRepository。
type stubRefundRepo struct {
	refunds   []*domain.Refund
	sumByID   map[int64]int64
	createErr error
}

func newStubRefundRepo(sums map[int64]int64) *stubRefundRepo {
	if sums == nil {
		sums = map[int64]int64{}
	}
	return &stubRefundRepo{sumByID: sums}
}

func (r *stubRefundRepo) Create(_ context.Context, refund *domain.Refund) error {
	if r.createErr != nil {
		return r.createErr
	}
	r.refunds = append(r.refunds, refund)
	return nil
}

func (r *stubRefundRepo) SumByPaymentID(_ context.Context, paymentID int64) (int64, error) {
	return r.sumByID[paymentID], nil
}

// paidPayment 构建一个用于测试的最小已支付 Payment。
func paidPayment(id, cents int64) *domain.Payment {
	return &domain.Payment{ID: id, OrderID: 1, AmountCents: cents, Status: PaymentStatusPaid}
}

func newRefundTestSvc(payment *domain.Payment, alreadyRefunded int64) (*RefundService, *stubRefundRepo) {
	pr := &stubPayRepoForRefund{payment: payment}
	rr := newStubRefundRepo(map[int64]int64{payment.ID: alreadyRefunded})
	return NewRefundService(pr, rr), rr
}

func TestRefundService_OK(t *testing.T) {
	p := paidPayment(5, 10000)
	svc, rr := newRefundTestSvc(p, 0)

	err := svc.Create(context.Background(), 5, 5000, "partial refund")

	require.NoError(t, err)
	assert.Len(t, rr.refunds, 1)
	assert.Equal(t, int64(5000), rr.refunds[0].AmountCents)
	assert.Equal(t, RefundStatusPending, rr.refunds[0].Status)
}

func TestRefundService_FullRefund(t *testing.T) {
	p := paidPayment(5, 10000)
	svc, rr := newRefundTestSvc(p, 0)

	err := svc.Create(context.Background(), 5, 10000, "full refund")

	require.NoError(t, err)
	assert.Len(t, rr.refunds, 1)
}

func TestRefundService_NegativeAmount(t *testing.T) {
	p := paidPayment(5, 10000)
	svc, _ := newRefundTestSvc(p, 0)

	err := svc.Create(context.Background(), 5, -100, "x")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "positive")
}

func TestRefundService_ZeroAmount(t *testing.T) {
	p := paidPayment(5, 10000)
	svc, _ := newRefundTestSvc(p, 0)

	err := svc.Create(context.Background(), 5, 0, "x")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "positive")
}

func TestRefundService_PaymentNotPaid(t *testing.T) {
	p := &domain.Payment{ID: 5, AmountCents: 10000, Status: PaymentStatusPending}
	svc, _ := newRefundTestSvc(p, 0)

	err := svc.Create(context.Background(), 5, 1000, "x")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not in paid status")
}

func TestRefundService_ExceedsOriginalAmount(t *testing.T) {
	p := paidPayment(5, 10000)
	svc, _ := newRefundTestSvc(p, 0)

	err := svc.Create(context.Background(), 5, 10001, "too much")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds")
}

func TestRefundService_CumulativeExceeds(t *testing.T) {
	p := paidPayment(5, 10000)
	svc, _ := newRefundTestSvc(p, 8000) // 8000 already refunded

	err := svc.Create(context.Background(), 5, 3000, "extra") // 8000+3000 > 10000
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds")
}

func TestRefundService_CumulativeExact(t *testing.T) {
	p := paidPayment(5, 10000)
	svc, rr := newRefundTestSvc(p, 7000) // 7000 already refunded, 3000 remaining

	err := svc.Create(context.Background(), 5, 3000, "exact remainder")
	require.NoError(t, err)
	assert.Len(t, rr.refunds, 1)
}

func TestRefundService_PaymentNotFound(t *testing.T) {
	pr := &stubPayRepoForRefund{findErr: errors.New("not found")}
	rr := newStubRefundRepo(nil)
	svc := NewRefundService(pr, rr)

	err := svc.Create(context.Background(), 99, 1000, "x")
	assert.Error(t, err)
}
