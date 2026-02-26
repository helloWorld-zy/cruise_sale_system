package service

import "testing"

type fakeRefundRepo struct{ called bool }

func (f *fakeRefundRepo) Create(_ interface{}) error { f.called = true; return nil }

func TestRefundServiceCreate(t *testing.T) {
	svc := NewRefundService(&fakeRefundRepo{})
	_ = svc.Create(1, 100, "cancel")
}
