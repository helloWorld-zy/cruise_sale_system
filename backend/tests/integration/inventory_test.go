package integration

import (
	"testing"
)

func setup() {
}

func TestInventoryLock(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
}