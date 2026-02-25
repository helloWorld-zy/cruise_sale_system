package service

import (
	"testing"
	"time"
)

type fakeCodeStore struct {
	ok        bool
	saveCalls int
}

func (f *fakeCodeStore) Save(phone, code string) error {
	f.saveCalls++
	return nil
}
func (f *fakeCodeStore) Verify(phone, code string) bool { return f.ok }

func TestUserAuthVerifySMS(t *testing.T) {
	now := time.Now()
	svc := NewUserAuthServiceWithPolicy(&fakeCodeStore{ok: true}, UserAuthPolicy{
		CodeTTL:        time.Minute,
		ResendInterval: time.Second,
		MaxAttempts:    2,
		LockDuration:   time.Minute,
		Now:            func() time.Time { return now },
	})
	if err := svc.SendSMS("13800000000", "1234"); err != nil {
		t.Fatal(err)
	}
	if !svc.VerifySMS("13800000000", "1234") {
		t.Fatal("expected sms verify")
	}
}

func TestUserAuthSendSMSRateLimit(t *testing.T) {
	now := time.Now()
	store := &fakeCodeStore{ok: true}
	svc := NewUserAuthServiceWithPolicy(store, UserAuthPolicy{
		ResendInterval: time.Minute,
		Now:            func() time.Time { return now },
	})

	if err := svc.SendSMS("13800000000", "1234"); err != nil {
		t.Fatal(err)
	}
	if err := svc.SendSMS("13800000000", "1235"); err == nil {
		t.Fatal("expected rate limit error")
	}
	if store.saveCalls != 1 {
		t.Fatalf("expected one save call, got %d", store.saveCalls)
	}
}

func TestUserAuthVerifyExpiryAndLockout(t *testing.T) {
	now := time.Now()
	store := &fakeCodeStore{ok: false}
	svc := NewUserAuthServiceWithPolicy(store, UserAuthPolicy{
		CodeTTL:      10 * time.Second,
		MaxAttempts:  2,
		LockDuration: time.Minute,
		Now:          func() time.Time { return now },
	})

	if err := svc.SendSMS("13800000000", "1234"); err != nil {
		t.Fatal(err)
	}
	if svc.VerifySMS("13800000000", "bad") {
		t.Fatal("expected first bad code to fail")
	}
	if svc.VerifySMS("13800000000", "bad") {
		t.Fatal("expected second bad code to fail")
	}
	if svc.VerifySMS("13800000000", "bad") {
		t.Fatal("expected lockout to reject verification")
	}

	now = now.Add(2 * time.Minute)
	if svc.VerifySMS("13800000000", "bad") {
		t.Fatal("expected expired code to fail")
	}
}
