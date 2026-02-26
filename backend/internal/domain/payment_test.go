package domain

import "testing"

func TestPaymentModelFields(t *testing.T) {
	p := Payment{OrderID: 1, AmountCents: 100}
	if p.AmountCents == 0 {
		t.Fatal("expected amount")
	}
}
