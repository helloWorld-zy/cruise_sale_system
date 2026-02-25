package database

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type dummyDialector struct{}

func (d dummyDialector) Name() string { return "dummy" }
func (d dummyDialector) Initialize(db *gorm.DB) error {
	db.ConnPool = nil
	return nil
}
func (d dummyDialector) Migrator(db *gorm.DB) gorm.Migrator                                  { return nil }
func (d dummyDialector) DataTypeOf(field *schema.Field) string                               { return "" }
func (d dummyDialector) DefaultValueOf(field *schema.Field) clause.Expression                { return nil }
func (d dummyDialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {}
func (d dummyDialector) QuoteTo(writer clause.Writer, str string)                            {}
func (d dummyDialector) Explain(sql string, vars ...interface{}) string                      { return "" }

func TestDatabaseConnectDBError(t *testing.T) {
	_, err := ConnectWithDialector(dummyDialector{}, Config{})
	if err == nil {
		t.Error("expected db.DB() error")
	}
}

func TestDatabaseConnect(t *testing.T) {
	// Simple coverage for database.go errors
	// Real connection test needs a DB
	cfg := Config{
		Host:   "localhost",
		Port:   5432,
		User:   "none",
		DBName: "none",
	}
	_, err := Connect(cfg)
	if err == nil {
		t.Errorf("Expected error for invalid connection")
	}

	d := sqlite.Open("file::memory:?cache=shared")
	db, err := ConnectWithDialector(d, cfg)
	if err != nil {
		t.Fatalf("ConnectWithDialector failed: %v", err)
	}
	if db == nil {
		t.Fatal("expected non-nil database instance")
	}
}
