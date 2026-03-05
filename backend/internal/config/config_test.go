package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	content := []byte(`
server:
  port: 8080
  mode: debug
database:
  driver: sqlite
  dsn: test.db
jwt:
  secret: testsecret
  expire: 24h
`)
	requireFile(t, tmpDir, "config.yaml", content)
	requireFile(t, tmpDir, "config.test.yaml", content)

	// Test mapping
	os.Setenv("ENV", "test")
	defer os.Unsetenv("ENV")

	// Call Load with the temp path
	cfg := Load(tmpDir)

	if strings.TrimPrefix(cfg.Server.Port, ":") != "8080" {
		t.Errorf("Expected port 8080 (optional ':' prefix), got %s", cfg.Server.Port)
	}
}

func TestLoadPanicsOnMissingConfig(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic on missing config")
		}
	}()
	Load(t.TempDir() + "/doesnotexist")
}

func TestLoadPanicsOnBadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	requireFile(t, tmpDir, "config.yaml", []byte("\tinvalid_yaml: \n- item\n"))
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic on bad config")
		}
	}()
	Load(tmpDir)
}

func TestLoadPanicsOnUnmarshalError(t *testing.T) {
	tmpDir := t.TempDir()
	requireFile(t, tmpDir, "config.yaml", []byte("database:\n  port: not-a-number\n"))
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic on unmarshal error")
		}
	}()
	Load(tmpDir)
}

func TestLoadDatabaseEnvFallbacksFromPostgresVars(t *testing.T) {
	tmpDir := t.TempDir()
	requireFile(t, tmpDir, "config.yaml", []byte(`
server:
  port: ":8080"
database:
  host: ""
  port: 0
  user: ""
  password: ""
  dbname: ""
  sslmode: "disable"
`))

	t.Setenv("POSTGRES_HOST", "127.0.0.1")
	t.Setenv("POSTGRES_PORT", "15432")
	t.Setenv("POSTGRES_USER", "cruise")
	t.Setenv("POSTGRES_PASSWORD", "postgres")
	t.Setenv("POSTGRES_DB", "cruise_booking")

	cfg := Load(tmpDir)

	if cfg.Database.Host != "127.0.0.1" {
		t.Fatalf("expected host from POSTGRES_HOST, got %q", cfg.Database.Host)
	}
	if cfg.Database.Port != 15432 {
		t.Fatalf("expected port 15432 from POSTGRES_PORT, got %d", cfg.Database.Port)
	}
	if cfg.Database.User != "cruise" {
		t.Fatalf("expected user from POSTGRES_USER, got %q", cfg.Database.User)
	}
	if cfg.Database.Password != "postgres" {
		t.Fatalf("expected password from POSTGRES_PASSWORD, got %q", cfg.Database.Password)
	}
	if cfg.Database.DBName != "cruise_booking" {
		t.Fatalf("expected dbname from POSTGRES_DB, got %q", cfg.Database.DBName)
	}
}

func TestLoadDatabasePortIgnoresInvalidPostgresPort(t *testing.T) {
	tmpDir := t.TempDir()
	requireFile(t, tmpDir, "config.yaml", []byte(`
server:
  port: ":8080"
database:
  host: "localhost"
  port: 15432
  user: "cruise"
  password: "postgres"
  dbname: "cruise_booking"
`))

	t.Setenv("CRUISE_DATABASE_PORT", "")
	t.Setenv("POSTGRES_PORT", "not-a-number")

	cfg := Load(tmpDir)
	if cfg.Database.Port != 15432 {
		t.Fatalf("expected config file port to remain 15432, got %d", cfg.Database.Port)
	}
}

func TestLoadUploadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	requireFile(t, tmpDir, "config.yaml", []byte(`
server:
  port: ":8080"
upload:
  storagedir: "tmp/uploads"
  publicpath: "/assets/uploads"
  maxfilesize: 2097152
`))

	cfg := Load(tmpDir)
	if cfg.Upload.StorageDir != "tmp/uploads" {
		t.Fatalf("expected upload storagedir tmp/uploads, got %q", cfg.Upload.StorageDir)
	}
	if cfg.Upload.PublicPath != "/assets/uploads" {
		t.Fatalf("expected upload publicpath /assets/uploads, got %q", cfg.Upload.PublicPath)
	}
	if cfg.Upload.MaxFileSize != 2097152 {
		t.Fatalf("expected upload maxfilesize 2097152, got %d", cfg.Upload.MaxFileSize)
	}
}

func TestLoadUploadConfigDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	requireFile(t, tmpDir, "config.yaml", []byte(`
server:
  port: ":8080"
upload:
  storagedir: ""
  publicpath: ""
  maxfilesize: 0
`))

	cfg := Load(tmpDir)
	if cfg.Upload.StorageDir != "uploads" {
		t.Fatalf("expected default upload storagedir uploads, got %q", cfg.Upload.StorageDir)
	}
	if cfg.Upload.PublicPath != "/uploads" {
		t.Fatalf("expected default upload publicpath /uploads, got %q", cfg.Upload.PublicPath)
	}
	if cfg.Upload.MaxFileSize != 10*1024*1024 {
		t.Fatalf("expected default upload maxfilesize 10485760, got %d", cfg.Upload.MaxFileSize)
	}
}

func requireFile(t *testing.T, dir, name string, content []byte) {
	err := os.WriteFile(filepath.Join(dir, name), content, 0644)
	if err != nil {
		t.Fatal(err)
	}
}
