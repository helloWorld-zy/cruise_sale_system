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

func TestUserAuthAlipayLogin(t *testing.T) {
	svc := NewUserAuthService(&fakeCodeStore{ok: true})
	signedUID := "alipay_uid_001"
	sig := svc.signAlipayUID(signedUID)

	token, err := svc.AlipayLogin("alipay_uid_001", signedUID, sig)
	if err != nil {
		t.Fatalf("expected alipay login success, got error: %v", err)
	}
	if token != signedUID {
		t.Fatalf("expected verified uid %s, got %s", signedUID, token)
	}
}

func TestUserAuthAlipayLoginRejectsInvalidSignature(t *testing.T) {
	svc := NewUserAuthService(&fakeCodeStore{ok: true})

	_, err := svc.AlipayLogin("alipay_uid_001", "alipay_uid_001", "bad-signature")
	if err == nil {
		t.Fatal("expected invalid signature error")
	}
}

func TestUserAuthAlipayLoginRejectsForgedClientUID(t *testing.T) {
	svc := NewUserAuthService(&fakeCodeStore{ok: true})
	providerUID := "alipay_uid_real"
	sig := svc.signAlipayUID(providerUID)

	_, err := svc.AlipayLogin("alipay_uid_forged", providerUID, sig)
	if err == nil {
		t.Fatal("expected forged uid to be rejected")
	}
}

func TestUserAuthBindAccount(t *testing.T) {
	now := time.Now()
	store := &fakeCodeStore{ok: true}
	svc := NewUserAuthServiceWithPolicy(store, UserAuthPolicy{
		CodeTTL:        time.Minute,
		ResendInterval: time.Second,
		MaxAttempts:    5,
		LockDuration:   time.Minute,
		Now:            func() time.Time { return now },
	})
	// 先发 SMS 验证码
	if err := svc.SendSMS("13800000000", "1234"); err != nil {
		t.Fatal(err)
	}
	if err := svc.AuthorizeBinding(1, "13800000000", "1234"); err != nil {
		t.Fatalf("expected authorize success: %v", err)
	}
	err := svc.BindAccount(1, "alipay", "alipay_uid_001")
	if err != nil {
		t.Fatal("expected bind success")
	}
}

func TestUserAuthBindAccountRequiresConfirmation(t *testing.T) {
	svc := NewUserAuthService(&fakeCodeStore{ok: true})

	err := svc.BindAccount(1, "alipay", "alipay_uid_001")
	if err == nil {
		t.Fatal("expected bind to fail without confirmation")
	}
}

func TestUserAuthBindAccountRejectsDuplicateIdentifier(t *testing.T) {
	now := time.Now()
	store := &fakeCodeStore{ok: true}
	svc := NewUserAuthServiceWithPolicy(store, UserAuthPolicy{
		CodeTTL:        time.Minute,
		ResendInterval: time.Second,
		MaxAttempts:    5,
		LockDuration:   time.Minute,
		Now:            func() time.Time { return now },
	})

	if err := svc.SendSMS("13800000001", "0001"); err != nil {
		t.Fatal(err)
	}
	if err := svc.AuthorizeBinding(1, "13800000001", "0001"); err != nil {
		t.Fatalf("authorize user1 failed: %v", err)
	}
	if err := svc.BindAccount(1, "alipay", "alipay_uid_001"); err != nil {
		t.Fatalf("user1 bind failed: %v", err)
	}

	if err := svc.SendSMS("13800000002", "0002"); err != nil {
		t.Fatal(err)
	}
	if err := svc.AuthorizeBinding(2, "13800000002", "0002"); err != nil {
		t.Fatalf("authorize user2 failed: %v", err)
	}
	err := svc.BindAccount(2, "alipay", "alipay_uid_001")
	if err == nil {
		t.Fatal("expected duplicate binding to be rejected")
	}
}
