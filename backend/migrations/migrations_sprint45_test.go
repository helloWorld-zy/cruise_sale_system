package migrations

import (
	"os"
	"testing"
)

func TestSprint45MigrationFilesExist(t *testing.T) {
	files := []string{
		"000018_cabin_type_pricing_unification.up.sql",
		"000018_cabin_type_pricing_unification.down.sql",
	}

	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			t.Fatalf("expected %s", f)
		}
	}
}
