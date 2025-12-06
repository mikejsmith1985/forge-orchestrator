// Package config provides configuration management for Forge Orchestrator.
// Configuration is stored in a JSON file in the user's config directory.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// ShellType represents the type of shell to use for command execution.
type ShellType string

const (
	ShellBash       ShellType = "bash"
	ShellCmd        ShellType = "cmd"
	ShellPowerShell ShellType = "powershell"
	ShellWSL        ShellType = "wsl"
)

// Config represents the application configuration.
type Config struct {
	// Shell configuration
	Shell ShellConfig `json:"shell"`

	// Update configuration
	Update UpdateConfig `json:"update"`

	// Server configuration
	Server ServerConfig `json:"server"`
}

// ShellConfig contains shell-related settings.
type ShellConfig struct {
	// Type is the shell type to use (bash, cmd, powershell, wsl)
	Type ShellType `json:"type"`

	// WSLDistro is the WSL distribution name (only used when Type is "wsl")
	WSLDistro string `json:"wsl_distro,omitempty"`

	// WSLUser is the WSL username (only used when Type is "wsl")
	WSLUser string `json:"wsl_user,omitempty"`

	// RootDir is the starting directory for the terminal (empty = current working directory)
	RootDir string `json:"root_dir,omitempty"`
}

// UpdateConfig contains update-related settings.
type UpdateConfig struct {
	// CheckOnStartup determines if updates should be checked on app start
	CheckOnStartup bool `json:"check_on_startup"`

	// CheckIntervalMinutes is how often to check for updates (0 = disabled)
	CheckIntervalMinutes int `json:"check_interval_minutes"`

	// AutoDownload determines if updates should be downloaded automatically
	AutoDownload bool `json:"auto_download"`
}

// ServerConfig contains server-related settings.
type ServerConfig struct {
	// Port is the preferred port (0 = auto-select)
	Port int `json:"port"`

	// OpenBrowser determines if browser should open on startup
	OpenBrowser bool `json:"open_browser"`
}

var (
	currentConfig *Config
	configMu      sync.RWMutex
	configPath    string
)

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	shellType := ShellBash
	if runtime.GOOS == "windows" {
		shellType = ShellCmd
	}

	return &Config{
		Shell: ShellConfig{
			Type: shellType,
		},
		Update: UpdateConfig{
			CheckOnStartup:       true,
			CheckIntervalMinutes: 30,
			AutoDownload:         false,
		},
		Server: ServerConfig{
			Port:        8080,
			OpenBrowser: true,
		},
	}
}

// GetConfigDir returns the configuration directory path.
func GetConfigDir() (string, error) {
	var dir string

	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(home, "Library", "Application Support")
	default: // linux and others
		dir = os.Getenv("XDG_CONFIG_HOME")
		if dir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			dir = filepath.Join(home, ".config")
		}
	}

	configDir := filepath.Join(dir, "forge-orchestrator")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return configDir, nil
}

// GetDataDir returns the data directory path (for database, logs, etc.).
func GetDataDir() (string, error) {
	var dir string

	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("LOCALAPPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
		}
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(home, "Library", "Application Support")
	default: // linux and others
		dir = os.Getenv("XDG_DATA_HOME")
		if dir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			dir = filepath.Join(home, ".local", "share")
		}
	}

	dataDir := filepath.Join(dir, "forge-orchestrator")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", err
	}

	return dataDir, nil
}

// Load loads the configuration from disk.
func Load() (*Config, error) {
	configMu.Lock()
	defer configMu.Unlock()

	if currentConfig != nil {
		return currentConfig, nil
	}

	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	configPath = filepath.Join(configDir, "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			currentConfig = DefaultConfig()
			return currentConfig, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	currentConfig = &cfg
	return currentConfig, nil
}

// Save saves the configuration to disk.
func Save(cfg *Config) error {
	configMu.Lock()
	defer configMu.Unlock()

	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	configPath = filepath.Join(configDir, "config.json")

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return err
	}

	currentConfig = cfg
	return nil
}

// Get returns the current configuration (loads if not already loaded).
func Get() (*Config, error) {
	configMu.RLock()
	if currentConfig != nil {
		defer configMu.RUnlock()
		return currentConfig, nil
	}
	configMu.RUnlock()

	return Load()
}

// GetDatabasePath returns the path to the SQLite database.
func GetDatabasePath() (string, error) {
	dataDir, err := GetDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, "forge_ledger.db"), nil
}
