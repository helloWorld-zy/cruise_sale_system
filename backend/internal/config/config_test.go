package config

import (
	"os"
	"path/filepath"
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

	if cfg.Server.Port != "8080" {
		t.Errorf("Expected port 8080, got %s", cfg.Server.Port)
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

func requireFile(t *testing.T, dir, name string, content []byte) {
	err := os.WriteFile(filepath.Join(dir, name), content, 0644)
	if err != nil {
		t.Fatal(err)
	}
}
