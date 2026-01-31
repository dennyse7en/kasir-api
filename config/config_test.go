package config

import (
	"os"
	"testing"
)

func TestLoadConfigMissingFile(t *testing.T) {
	// Set environment variable
	expectedPort := ":9999"
	os.Setenv("SERVER_ADDRESS", expectedPort)
	defer os.Unsetenv("SERVER_ADDRESS")

	// Use os.TempDir() as a path that definitely exists but won't contain app.env (unless by weird coincidence)
	// This simulates the Railway environment where the file is missing.
	cfg, err := LoadConfig(os.TempDir())

	if err != nil {
		t.Fatalf("Expected no error when config file is missing, got: %v", err)
	}

	if cfg.ServerAddress != expectedPort {
		t.Errorf("Expected ServerAddress %s, got %s", expectedPort, cfg.ServerAddress)
	}
}
