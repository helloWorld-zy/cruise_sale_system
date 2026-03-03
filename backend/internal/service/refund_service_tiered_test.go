package service

import (
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestTieredRefundCalculation(t *testing.T) {
	rules := []domain.RefundRule{
		{MinDays: 31, MaxDays: 9999, RefundRate: 100},
		{MinDays: 15, MaxDays: 31, RefundRate: 80},
		{MinDays: 7, MaxDays: 15, RefundRate: 50},
		{MinDays: 0, MaxDays: 7, RefundRate: 0},
	}
	svc := &RefundService{}

	amount := svc.CalcRefundAmount(10000, 20, rules)
	if amount != 8000 {
		t.Fatalf("expected 8000, got %d", amount)
	}
}

func TestTieredRefund_FullRefund(t *testing.T) {
	rules := []domain.RefundRule{
		{MinDays: 31, MaxDays: 9999, RefundRate: 100},
		{MinDays: 15, MaxDays: 31, RefundRate: 80},
		{MinDays: 7, MaxDays: 15, RefundRate: 50},
		{MinDays: 0, MaxDays: 7, RefundRate: 0},
	}
	svc := &RefundService{}

	amount := svc.CalcRefundAmount(10000, 31, rules)
	assert.Equal(t, int64(10000), amount)
}

func TestTieredRefund_50Percent(t *testing.T) {
	rules := []domain.RefundRule{
		{MinDays: 31, MaxDays: 9999, RefundRate: 100},
		{MinDays: 15, MaxDays: 31, RefundRate: 80},
		{MinDays: 7, MaxDays: 15, RefundRate: 50},
		{MinDays: 0, MaxDays: 7, RefundRate: 0},
	}
	svc := &RefundService{}

	amount := svc.CalcRefundAmount(10000, 10, rules)
	assert.Equal(t, int64(5000), amount)
}

func TestTieredRefund_NoRefund(t *testing.T) {
	rules := []domain.RefundRule{
		{MinDays: 31, MaxDays: 9999, RefundRate: 100},
		{MinDays: 15, MaxDays: 31, RefundRate: 80},
		{MinDays: 7, MaxDays: 15, RefundRate: 50},
		{MinDays: 0, MaxDays: 7, RefundRate: 0},
	}
	svc := &RefundService{}

	amount := svc.CalcRefundAmount(10000, 3, rules)
	assert.Equal(t, int64(0), amount)
}

func TestTieredRefund_Boundary(t *testing.T) {
	rules := []domain.RefundRule{
		{MinDays: 31, MaxDays: 9999, RefundRate: 100},
		{MinDays: 15, MaxDays: 31, RefundRate: 80},
		{MinDays: 7, MaxDays: 15, RefundRate: 50},
		{MinDays: 0, MaxDays: 7, RefundRate: 0},
	}
	svc := &RefundService{}

	// 15天 = 15-30天范围 → 80%
	amount := svc.CalcRefundAmount(10000, 15, rules)
	assert.Equal(t, int64(8000), amount)

	// 14天 = 7-14天范围 → 50%
	amount = svc.CalcRefundAmount(10000, 14, rules)
	assert.Equal(t, int64(5000), amount)

	// 6天 = 0-6天范围 → 0%
	amount = svc.CalcRefundAmount(10000, 6, rules)
	assert.Equal(t, int64(0), amount)

	// 30天 = 15-30天范围 → 80%
	amount = svc.CalcRefundAmount(10000, 30, rules)
	assert.Equal(t, int64(8000), amount)
}

func TestTieredRefund_BoundaryExactDays30_7_0(t *testing.T) {
	rules := []domain.RefundRule{
		{MinDays: 31, MaxDays: 9999, RefundRate: 100},
		{MinDays: 15, MaxDays: 31, RefundRate: 80},
		{MinDays: 7, MaxDays: 15, RefundRate: 50},
		{MinDays: 0, MaxDays: 7, RefundRate: 0},
	}
	svc := &RefundService{}

	assert.Equal(t, int64(8000), svc.CalcRefundAmount(10000, 30, rules))
	assert.Equal(t, int64(5000), svc.CalcRefundAmount(10000, 7, rules))
	assert.Equal(t, int64(0), svc.CalcRefundAmount(10000, 0, rules))
}

func TestTieredRefund_IntegerArithmeticNoPrecisionLoss(t *testing.T) {
	rules := []domain.RefundRule{{MinDays: 15, MaxDays: 31, RefundRate: 80}}
	svc := &RefundService{}

	amount := svc.CalcRefundAmount(9999, 20, rules)
	assert.Equal(t, int64(7999), amount)
}
