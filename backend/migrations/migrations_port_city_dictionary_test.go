package migrations

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestPortCityDictionaryMigrationFilesExist(t *testing.T) {
	files := []string{
		"000024_seed_port_city_dictionary.up.sql",
		"000024_seed_port_city_dictionary.down.sql",
	}
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			t.Fatalf("expected migration file %s to exist: %v", f, err)
		}
	}
}

func TestPortCityDictionaryMigrationSeedsCoords(t *testing.T) {
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

	upBytes, err := os.ReadFile("000024_seed_port_city_dictionary.up.sql")
	if err != nil {
		t.Fatalf("read up migration failed: %v", err)
	}
	for _, stmt := range sqliteCompatibleStatements(string(upBytes)) {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("execute up statement failed: %v\nstmt=%s", err, stmt)
		}
	}

	type destinationSeed struct {
		Name      string
		Country   string
		Latitude  float64
		Longitude float64
		Keywords  string
	}

	var miami destinationSeed
	if err := db.Raw(`SELECT name, country, latitude, longitude, keywords FROM custom_destinations WHERE name = ? AND country = ?`, "迈阿密", "美国").Scan(&miami).Error; err != nil {
		t.Fatalf("query miami seed failed: %v", err)
	}
	if miami.Name != "迈阿密" || miami.Country != "美国" {
		t.Fatalf("expected Miami seed to exist, got %+v", miami)
	}
	if miami.Latitude == 0 || miami.Longitude == 0 {
		t.Fatalf("expected Miami seed coords, got %+v", miami)
	}
	if miami.Keywords == "" {
		t.Fatalf("expected Miami seed keywords, got %+v", miami)
	}

	var buenos destinationSeed
	if err := db.Raw(`SELECT name, country, latitude, longitude FROM custom_destinations WHERE name = ? AND country = ?`, "布宜诺斯艾利斯", "阿根廷").Scan(&buenos).Error; err != nil {
		t.Fatalf("query Buenos Aires seed failed: %v", err)
	}
	if buenos.Latitude == 0 || buenos.Longitude == 0 {
		t.Fatalf("expected Buenos Aires seed coords, got %+v", buenos)
	}

	downBytes, err := os.ReadFile("000024_seed_port_city_dictionary.down.sql")
	if err != nil {
		t.Fatalf("read down migration failed: %v", err)
	}
	for _, stmt := range sqliteCompatibleStatements(string(downBytes)) {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("execute down statement failed: %v\nstmt=%s", err, stmt)
		}
	}

	var remaining int64
	if err := db.Raw(`SELECT COUNT(*) FROM custom_destinations WHERE description = ?`, "system_port_city_dictionary_seed_v1").Scan(&remaining).Error; err != nil {
		t.Fatalf("count remaining seeds failed: %v", err)
	}
	if remaining != 0 {
		t.Fatalf("expected seeded rows removed on down migration, got %d", remaining)
	}
}
