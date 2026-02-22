package logger

import "testing"

func TestNewLogger(t *testing.T) {
	l := New("debug", "logs/test.log")
	if l == nil {
		t.Fatal("expected logger to be non-nil")
	}
}
