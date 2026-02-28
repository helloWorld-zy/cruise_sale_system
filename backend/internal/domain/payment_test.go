package domain

import "testing"

// TestPaymentModelFields 测试 Payment 模型字段。
func TestPaymentModelFields(t *testing.T) {
	p := Payment{OrderID: 1, AmountCents: 100}
	if p.AmountCents == 0 {
		t.Fatal("expected amount")
	}
}
