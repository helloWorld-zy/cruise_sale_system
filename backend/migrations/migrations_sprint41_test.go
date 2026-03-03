package migrations

import (
	"os"
	"strings"
	"testing"
)

func TestSprint41MigrationFilesExist(t *testing.T) {
	files := []string{
		"000008_sprint41_extend.up.sql",
		"000008_sprint41_extend.down.sql",
	}
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			t.Fatalf("missing %s", f)
		}
	}
}

func TestSprint41MigrationUpContainsKeyStatements(t *testing.T) {
	data, err := os.ReadFile("000008_sprint41_extend.up.sql")
	if err != nil {
		t.Fatalf("read up migration failed: %v", err)
	}
	sql := strings.ToUpper(string(data))

	// 验证关键字段扩展语句存在，避免仅检查文件存在导致的假通过。
	mustContain := []string{
		"ALTER TABLE CRUISES ADD COLUMN IF NOT EXISTS CODE",
		"ALTER TABLE CRUISES ADD COLUMN IF NOT EXISTS CREW_COUNT",
		"ALTER TABLE CABIN_TYPES ADD COLUMN IF NOT EXISTS FLOOR_PLAN_URL",
		"ALTER TABLE CABIN_SKUS ADD COLUMN IF NOT EXISTS POSITION",
		"ALTER TABLE CABIN_SKUS ADD COLUMN IF NOT EXISTS HAS_BALCONY",
		"ALTER TABLE CABIN_INVENTORIES ADD COLUMN IF NOT EXISTS ALERT_THRESHOLD",
		"ALTER TABLE CABIN_PRICES ADD COLUMN IF NOT EXISTS PRICE_TYPE",
		"ALTER TABLE FACILITIES ADD COLUMN IF NOT EXISTS OPEN_HOURS",
		"ALTER TABLE FACILITY_CATEGORIES ADD COLUMN IF NOT EXISTS STATUS",
		"ALTER TABLE IMAGES ADD COLUMN IF NOT EXISTS IS_PRIMARY",
		"CREATE INDEX IF NOT EXISTS IDX_IMAGES_ENTITY",
	}
	for _, frag := range mustContain {
		if !strings.Contains(sql, frag) {
			t.Fatalf("missing sql fragment in up migration: %s", frag)
		}
	}
}

func TestSprint41MigrationDownContainsKeyStatements(t *testing.T) {
	data, err := os.ReadFile("000008_sprint41_extend.down.sql")
	if err != nil {
		t.Fatalf("read down migration failed: %v", err)
	}
	sql := strings.ToUpper(string(data))

	// 验证关键回滚语句存在，确保 up/down 具备对称回滚能力。
	mustContain := []string{
		"DROP INDEX IF EXISTS IDX_IMAGES_ENTITY",
		"ALTER TABLE IMAGES DROP COLUMN IF EXISTS IS_PRIMARY",
		"ALTER TABLE FACILITY_CATEGORIES DROP COLUMN IF EXISTS STATUS",
		"ALTER TABLE CABIN_TYPES DROP COLUMN IF EXISTS FLOOR_PLAN_URL",
		"ALTER TABLE CABIN_SKUS DROP COLUMN IF EXISTS AMENITIES",
		"ALTER TABLE CABIN_INVENTORIES DROP COLUMN IF EXISTS ALERT_THRESHOLD",
		"ALTER TABLE CABIN_PRICES DROP COLUMN IF EXISTS PRICE_TYPE",
		"ALTER TABLE CRUISES DROP COLUMN IF EXISTS DECK_COUNT",
	}
	for _, frag := range mustContain {
		if !strings.Contains(sql, frag) {
			t.Fatalf("missing sql fragment in down migration: %s", frag)
		}
	}
}
