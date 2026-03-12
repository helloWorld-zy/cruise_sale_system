package migrations

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestPortCityDictionaryExtendMigrationFilesExist(t *testing.T) {
	files := []string{
		"000025_extend_port_city_dictionary.up.sql",
		"000025_extend_port_city_dictionary.down.sql",
	}
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			t.Fatalf("expected migration file %s to exist: %v", f, err)
		}
	}
}

func TestPortCityDictionaryExtendMigrationSeedsMorePorts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	if err := db.Exec(`
		CREATE TABLE custom_destinations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			country TEXT NOT NULL DEFAULT '',
			latitude REAL,
			longitude REAL,
			keywords TEXT NOT NULL DEFAULT '',
			description TEXT NOT NULL DEFAULT '',
			status SMALLINT NOT NULL DEFAULT 1,
			sort_order INT NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("create custom_destinations failed: %v", err)
	}

	upBytes, err := os.ReadFile("000025_extend_port_city_dictionary.up.sql")
	if err != nil {
		t.Fatalf("read up migration failed: %v", err)
	}
	for _, stmt := range sqliteCompatibleStatements(string(upBytes)) {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("execute up statement failed: %v\nstmt=%s", err, stmt)
		}
	}

	type seed struct {
		Name, Country       string
		Latitude, Longitude float64
	}
	var vancouver seed
	if err := db.Raw(`SELECT name, country, latitude, longitude FROM custom_destinations WHERE name = ? AND country = ?`, "温哥华", "加拿大").Scan(&vancouver).Error; err != nil {
		t.Fatalf("query Vancouver seed failed: %v", err)
	}
	if vancouver.Latitude == 0 || vancouver.Longitude == 0 {
		t.Fatalf("expected Vancouver coords, got %+v", vancouver)
	}

	var reykjavik seed
	if err := db.Raw(`SELECT name, country, latitude, longitude FROM custom_destinations WHERE name = ? AND country = ?`, "雷克雅未克", "冰岛").Scan(&reykjavik).Error; err != nil {
		t.Fatalf("query Reykjavik seed failed: %v", err)
	}
	if reykjavik.Latitude == 0 || reykjavik.Longitude == 0 {
		t.Fatalf("expected Reykjavik coords, got %+v", reykjavik)
	}
}
