package migrations

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestContentTemplateMigrationFilesExist(t *testing.T) {
	files := []string{
		"000021_content_templates.up.sql",
		"000021_content_templates.down.sql",
	}
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			t.Fatalf("expected migration file %s to exist: %v", f, err)
		}
	}
}

func TestContentTemplateMigrationExecuteUpDown(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}

	if err := db.Exec(`CREATE TABLE voyages (id INTEGER PRIMARY KEY AUTOINCREMENT, code TEXT NOT NULL);`).Error; err != nil {
		t.Fatalf("create voyages failed: %v", err)
	}

	upBytes, err := os.ReadFile("000021_content_templates.up.sql")
	if err != nil {
		t.Fatalf("read up migration failed: %v", err)
	}
	for _, stmt := range sqliteCompatibleStatements(string(upBytes)) {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("execute up statement failed: %v\nstmt=%s", err, stmt)
		}
	}

	assertTableExists(t, db, "content_templates")
	assertColumnExists(t, db, "voyages", "fee_note_template_id")
	assertColumnExists(t, db, "voyages", "fee_note_mode")
	assertColumnExists(t, db, "voyages", "fee_note_content_json")
	assertColumnExists(t, db, "voyages", "booking_notice_template_id")
	assertColumnExists(t, db, "voyages", "booking_notice_mode")
	assertColumnExists(t, db, "voyages", "booking_notice_content_json")

	downBytes, err := os.ReadFile("000021_content_templates.down.sql")
	if err != nil {
		t.Fatalf("read down migration failed: %v", err)
	}
	for _, stmt := range sqliteCompatibleStatements(string(downBytes)) {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("execute down statement failed: %v\nstmt=%s", err, stmt)
		}
	}

	assertTableMissing(t, db, "content_templates")
}

func sqliteCompatibleStatements(sql string) []string {
	normalized := strings.ReplaceAll(sql, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "BIGSERIAL", "INTEGER")
	normalized = strings.ReplaceAll(normalized, "TIMESTAMPTZ", "DATETIME")
	normalized = strings.ReplaceAll(normalized, " NOW()", " CURRENT_TIMESTAMP")

	rawStatements := strings.Split(normalized, ";")
	out := make([]string, 0, len(rawStatements)+8)
	for _, raw := range rawStatements {
		stmt := strings.TrimSpace(raw)
		if stmt == "" {
			continue
		}
		if strings.HasPrefix(strings.ToUpper(stmt), "ALTER TABLE ") {
			parts := strings.Split(stmt, "\n")
			head := strings.Fields(parts[0])
			table := ""
			if len(head) >= 3 {
				table = head[2]
			}
			if table == "" {
				out = append(out, stmt+";")
				continue
			}
			for _, part := range parts[1:] {
				line := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(part), ","), ","))
				if line == "" {
					continue
				}
				line = strings.Replace(line, "ADD COLUMN IF NOT EXISTS", "ADD COLUMN", 1)
				line = strings.Replace(line, "DROP COLUMN IF EXISTS", "DROP COLUMN", 1)
				out = append(out, fmt.Sprintf("ALTER TABLE %s %s;", table, line))
			}
			continue
		}
		out = append(out, stmt+";")
	}
	return out
}

func assertTableExists(t *testing.T, db *gorm.DB, table string) {
	t.Helper()
	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count).Error; err != nil {
		t.Fatalf("check table exists failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected table %s to exist", table)
	}
}

func assertTableMissing(t *testing.T, db *gorm.DB, table string) {
	t.Helper()
	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count).Error; err != nil {
		t.Fatalf("check table missing failed: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected table %s to be dropped", table)
	}
}

func assertColumnExists(t *testing.T, db *gorm.DB, table string, column string) {
	t.Helper()
	rows, err := db.Raw(fmt.Sprintf("PRAGMA table_info(%s)", table)).Rows()
	if err != nil {
		t.Fatalf("query pragma failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name string
		var ctype string
		var notnull int
		var dflt any
		var pk int
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			t.Fatalf("scan pragma failed: %v", err)
		}
		if name == column {
			return
		}
	}
	t.Fatalf("expected column %s on table %s", column, table)
}
