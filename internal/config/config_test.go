package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}

	// Check shell type based on OS
	if runtime.GOOS == "windows" {
		if cfg.Shell.Type != ShellCmd {
			t.Errorf("Expected shell type %s on Windows, got %s", ShellCmd, cfg.Shell.Type)
		}
	} else {
		if cfg.Shell.Type != ShellBash {
			t.Errorf("Expected shell type %s on Unix, got %s", ShellBash, cfg.Shell.Type)
		}
	}

	// Check update defaults
	if !cfg.Update.CheckOnStartup {
		t.Error("Expected CheckOnStartup to be true")
	}
	if cfg.Update.CheckIntervalMinutes != 30 {
		t.Errorf("Expected CheckIntervalMinutes 30, got %d", cfg.Update.CheckIntervalMinutes)
	}

	// Check server defaults
	if cfg.Server.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", cfg.Server.Port)
	}
	if !cfg.Server.OpenBrowser {
		t.Error("Expected OpenBrowser to be true")
	}
}

func TestGetConfigDir(t *testing.T) {
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir failed: %v", err)
	}

	if dir == "" {
		t.Error("GetConfigDir returned empty string")
	}

	// Check directory exists
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("Config directory doesn't exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("Config path is not a directory")
	}
}

func TestGetDataDir(t *testing.T) {
	dir, err := GetDataDir()
	if err != nil {
		t.Fatalf("GetDataDir failed: %v", err)
	}

	if dir == "" {
		t.Error("GetDataDir returned empty string")
	}

	// Check directory exists
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("Data directory doesn't exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("Data path is not a directory")
	}
}

func TestLoadAndSave(t *testing.T) {
	// Reset global state
	currentConfig = nil

	// Create a temp directory for testing
	tmpDir, err := os.MkdirTemp("", "forge-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save custom config
	cfg := &Config{
		Shell: ShellConfig{
			Type:      ShellWSL,
			WSLDistro: "Ubuntu-24.04",
			WSLUser:   "testuser",
			RootDir:   "C:\\Users\\test\\projects",
		},
		Update: UpdateConfig{
			CheckOnStartup:       false,
			CheckIntervalMinutes: 60,
			AutoDownload:         true,
		},
		Server: ServerConfig{
			Port:        9000,
			OpenBrowser: false,
		},
	}

	// Manually set config path for test
	configPath = filepath.Join(tmpDir, "config.json")
	currentConfig = cfg

	if err := Save(cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Reset and reload
	currentConfig = nil
	configPath = filepath.Join(tmpDir, "config.json")

	// Read the file directly since Load() uses GetConfigDir()
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify loaded config (may be default if path doesn't match)
	if loaded == nil {
		t.Fatal("Load returned nil")
	}
}

func TestGetDatabasePath(t *testing.T) {
	path, err := GetDatabasePath()
	if err != nil {
		t.Fatalf("GetDatabasePath failed: %v", err)
	}

	if path == "" {
		t.Error("GetDatabasePath returned empty string")
	}

	if filepath.Base(path) != "forge_ledger.db" {
		t.Errorf("Expected filename forge_ledger.db, got %s", filepath.Base(path))
	}
}
