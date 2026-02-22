package database

import "testing"

func TestDSN(t *testing.T) {
	dsn := BuildDSN(Config{
		Host:     "localhost",
		Port:     5432,
		User:     "user",
		Password: "pass",
		DBName:   "db",
		SSLMode:  "disable",
	})

	if dsn == "" {
		t.Fatal("expected DSN to be non-empty")
	}
}
