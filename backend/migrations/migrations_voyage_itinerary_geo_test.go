package migrations

import (
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestVoyageItineraryGeoMigrationFilesExist(t *testing.T) {
	files := []string{
		"000022_voyage_itinerary_geo.up.sql",
		"000022_voyage_itinerary_geo.down.sql",
	}
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			t.Fatalf("expected migration file %s to exist: %v", f, err)
		}
	}
}

func TestVoyageItineraryGeoMigrationExecuteUpDown(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	if err := db.Exec(`CREATE TABLE voyage_itineraries (id INTEGER PRIMARY KEY AUTOINCREMENT, city TEXT NOT NULL);`).Error; err != nil {
		t.Fatalf("create voyage_itineraries failed: %v", err)
	}

	upBytes, err := os.ReadFile("000022_voyage_itinerary_geo.up.sql")
	if err != nil {
		t.Fatalf("read up migration failed: %v", err)
	}
	for _, stmt := range sqliteCompatibleStatements(string(upBytes)) {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("execute up statement failed: %v\nstmt=%s", err, stmt)
		}
	}
	assertColumnExists(t, db, "voyage_itineraries", "latitude")
	assertColumnExists(t, db, "voyage_itineraries", "longitude")

	downBytes, err := os.ReadFile("000022_voyage_itinerary_geo.down.sql")
	if err != nil {
		t.Fatalf("read down migration failed: %v", err)
	}
	for _, stmt := range sqliteCompatibleStatements(string(downBytes)) {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("execute down statement failed: %v\nstmt=%s", err, stmt)
		}
	}
}
