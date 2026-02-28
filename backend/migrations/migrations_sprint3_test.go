package migrations

import (
	"os"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestSprint3MigrationFilesExist(t *testing.T) {
	files := []string{
		"000004_user_booking.up.sql",
		"000004_user_booking.down.sql",
		"000006_cabin_hold_unique.up.sql",
		"000006_cabin_hold_unique.down.sql",
	}
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			t.Fatalf("expected %s", f)
		}
	}
}

func TestSprint3HoldUniqueMigrationHasIndexStatements(t *testing.T) {
	up, err := os.ReadFile("000006_cabin_hold_unique.up.sql")
	if err != nil {
		t.Fatalf("read up migration failed: %v", err)
	}
	down, err := os.ReadFile("000006_cabin_hold_unique.down.sql")
	if err != nil {
		t.Fatalf("read down migration failed: %v", err)
	}

	upSQL := strings.ToUpper(string(up))
	downSQL := strings.ToUpper(string(down))

	if !strings.Contains(upSQL, "CREATE UNIQUE INDEX") || !strings.Contains(upSQL, "UQ_CABIN_HOLDS_SKU_USER") {
		t.Fatal("up migration must create unique index uq_cabin_holds_sku_user")
	}
	if !strings.Contains(downSQL, "DROP INDEX") || !strings.Contains(downSQL, "UQ_CABIN_HOLDS_SKU_USER") {
		t.Fatal("down migration must drop unique index uq_cabin_holds_sku_user")
	}
}

func TestSprint3HoldUniqueMigrationExecuteUpDown(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}

	createTable := `
CREATE TABLE cabin_holds (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  cabin_sku_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  qty INTEGER NOT NULL,
  expires_at DATETIME NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`
	if err := db.Exec(createTable).Error; err != nil {
		t.Fatalf("create table failed: %v", err)
	}

	seed := `
INSERT INTO cabin_holds (cabin_sku_id, user_id, qty, expires_at) VALUES
(101, 1001, 1, CURRENT_TIMESTAMP),
(101, 1001, 1, CURRENT_TIMESTAMP),
(102, 1002, 1, CURRENT_TIMESTAMP);`
	if err := db.Exec(seed).Error; err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	upBytes, err := os.ReadFile("000006_cabin_hold_unique.up.sql")
	if err != nil {
		t.Fatalf("read up failed: %v", err)
	}
	if err := db.Exec(string(upBytes)).Error; err != nil {
		t.Fatalf("execute up migration failed: %v", err)
	}

	var afterUpCount int64
	if err := db.Raw("SELECT COUNT(*) FROM cabin_holds WHERE cabin_sku_id = ? AND user_id = ?", 101, 1001).Scan(&afterUpCount).Error; err != nil {
		t.Fatalf("count after up failed: %v", err)
	}
	if afterUpCount != 1 {
		t.Fatalf("expected dedup to keep 1 row, got %d", afterUpCount)
	}

	dupInsertErr := db.Exec("INSERT INTO cabin_holds (cabin_sku_id, user_id, qty, expires_at) VALUES (?, ?, 1, CURRENT_TIMESTAMP)", 101, 1001).Error
	if dupInsertErr == nil {
		t.Fatal("expected duplicate insert to fail after unique index created")
	}

	downBytes, err := os.ReadFile("000006_cabin_hold_unique.down.sql")
	if err != nil {
		t.Fatalf("read down failed: %v", err)
	}
	if err := db.Exec(string(downBytes)).Error; err != nil {
		t.Fatalf("execute down migration failed: %v", err)
	}

	if err := db.Exec("INSERT INTO cabin_holds (cabin_sku_id, user_id, qty, expires_at) VALUES (?, ?, 1, CURRENT_TIMESTAMP)", 101, 1001).Error; err != nil {
		t.Fatalf("expected duplicate insert to succeed after dropping unique index, got: %v", err)
	}
}
