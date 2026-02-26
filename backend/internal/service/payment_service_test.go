package service

import "testing"

type fakeGateway struct{ called bool }

func (f *fakeGateway) CreatePay(_ int64, _ int64) (string, error) {
	f.called = true
	return "pay_url", nil
}

func TestPaymentServiceCreate(t *testing.T) {
	svc := NewPaymentService(&fakeGateway{})
	url, _ := svc.Create(1, 100)
	if url == "" {
		t.Fatal("expected pay url")
	}
}
