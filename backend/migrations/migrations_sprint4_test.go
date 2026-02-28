package migrations

import (
	"os"
	"testing"
)

func TestSprint4MigrationFilesExist(t *testing.T) {
	files := []string{
		"000005_payment_notify.up.sql",
		"000005_payment_notify.down.sql",
	}
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			t.Fatalf("expected %s", f)
		}
	}
}
