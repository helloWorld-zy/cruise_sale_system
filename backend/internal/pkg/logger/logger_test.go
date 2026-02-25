package logger

import "testing"

func TestNew(t *testing.T) {
	// Valid logLevel tests the happy path
	l1 := New("debug", t.TempDir()+"/debug.log")
	if l1 == nil {
		t.Fatal("Logger should not be nil")
	}

	// Invalid logLevel hits UnmarshalText error branch and falls back to info
	l2 := New("invalid_level_xyz", t.TempDir()+"/info.log")
	if l2 == nil {
		t.Fatal("Logger should not be nil even with invalid level")
	}
}
