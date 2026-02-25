package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain_Exit(t *testing.T) {
	// Override osExit to avoid exiting the test panic
	exited := false
	osExit = func(code int) {
		exited = true
	}
	defer func() { osExit = os.Exit }()

	// Provide an invalid config path to trigger run error
	os.Setenv("PORT", "invalid-port")
	defer os.Unsetenv("PORT")

	main()
	assert.True(t, exited)
}

func TestRunApp_Success(t *testing.T) {
	os.Setenv("CRUISE_DATABASE_HOST", "sqlite") // trigger sqlite connection!
	defer os.Unsetenv("CRUISE_DATABASE_HOST")
	os.Setenv("CRUISE_SERVER_PORT", "invalid-port") // fail at r.Run but pass all setups
	defer os.Unsetenv("CRUISE_SERVER_PORT")

	err := RunApp("../../")
	assert.Error(t, err) // fails at the very end when binding to invalid-port
}

func TestRunApp_CasbinError(t *testing.T) {
	os.Setenv("CRUISE_DATABASE_HOST", "sqlite")
	defer os.Unsetenv("CRUISE_DATABASE_HOST")

	// Create a temp dir with ONLY config.yaml, so Casbin fails to find rbac files
	tmpDir := t.TempDir()
	configData, _ := os.ReadFile("../../config.yaml")
	os.WriteFile(tmpDir+"/config.yaml", configData, 0644)

	err := RunApp(tmpDir)
	assert.Error(t, err) // fails at Casbin initialization
}
