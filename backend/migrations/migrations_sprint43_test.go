package migrations

import (
	"os"
	"testing"
)

func TestSprint43MigrationFilesExist(t *testing.T) {
	files := []string{
		"000015_company_admin_hardening.up.sql",
		"000015_company_admin_hardening.down.sql",
		"000016_facility_categories_soft_delete.up.sql",
		"000016_facility_categories_soft_delete.down.sql",
	}

	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			t.Fatalf("expected %s", f)
		}
	}
}
