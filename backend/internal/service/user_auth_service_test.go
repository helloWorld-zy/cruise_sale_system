package service

import "testing"

type fakeCodeStore struct{ ok bool }

func (f *fakeCodeStore) Save(phone, code string) error  { return nil }
func (f *fakeCodeStore) Verify(phone, code string) bool { return f.ok }

func TestUserAuthVerifySMS(t *testing.T) {
	svc := NewUserAuthService(&fakeCodeStore{ok: true})
	if !svc.VerifySMS("13800000000", "1234") {
		t.Fatal("expected sms verify")
	}
}
