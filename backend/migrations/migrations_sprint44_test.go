package migrations

import (
	"os"
	"testing"
)

func TestSprint44MigrationFilesExist(t *testing.T) {
	files := []string{
		"000017_voyage_unification.up.sql",
		"000017_voyage_unification.down.sql",
	}

	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			t.Fatalf("expected %s", f)
		}
	}
}
