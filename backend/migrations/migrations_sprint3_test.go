package migrations

import (
	"os"
	"testing"
)

func TestSprint3MigrationFilesExist(t *testing.T) {
	files := []string{
		"000003_user_booking.up.sql",
		"000003_user_booking.down.sql",
	}
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			t.Fatalf("expected %s", f)
		}
	}
}
