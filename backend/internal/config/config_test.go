package config

import "testing"

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("CRUISE_SERVER_PORT", ":9000")

	cfg := Load("../../")

	if cfg.Server.Port != ":9000" {
		t.Fatalf("expected port :9000, got %s", cfg.Server.Port)
	}
	if cfg.Database.Host == "" {
		t.Fatalf("expected database host to be set")
	}
}
