package service

import "testing"

type fakeNotifier struct{ called bool }

func (f *fakeNotifier) Send(_ string, _ string, _ string) error { f.called = true; return nil }

func TestNotifyServiceSend(t *testing.T) {
	svc := NewNotifyService(&fakeNotifier{})
	_ = svc.Send("sms", "template", "payload")
}
