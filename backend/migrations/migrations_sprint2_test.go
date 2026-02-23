package migrations

import (
	"os"
	"testing"
)

func TestSprint2MigrationFilesExist(t *testing.T) {
	files := []string{
		"000002_route_voyage_cabin.up.sql",
		"000002_route_voyage_cabin.down.sql",
	}
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			t.Fatalf("expected %s: %v", f, err)
		}
	}
}
