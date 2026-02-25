package service_test

import (
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/service"
)

type blackboxCodeStore struct{ ok bool }

func (b *blackboxCodeStore) Save(phone, code string) error  { return nil }
func (b *blackboxCodeStore) Verify(phone, code string) bool { return b.ok }

func TestUserAuthServiceBlackbox(t *testing.T) {
	now := time.Now()
	store := &blackboxCodeStore{ok: true}
	svc := service.NewUserAuthServiceWithPolicy(store, service.UserAuthPolicy{
		CodeTTL:        time.Minute,
		ResendInterval: time.Millisecond,
		MaxAttempts:    2,
		LockDuration:   time.Minute,
		Now:            func() time.Time { return now },
	})

	if err := svc.SendSMS("13800000000", "1234"); err != nil {
		t.Fatalf("SendSMS failed: %v", err)
	}
	if !svc.VerifySMS("13800000000", "1234") {
		t.Fatal("VerifySMS should be true")
	}
	if got := svc.WechatLogin("wx-open-id"); got != "wx-open-id" {
		t.Fatalf("WechatLogin mismatch: %s", got)
	}
}
