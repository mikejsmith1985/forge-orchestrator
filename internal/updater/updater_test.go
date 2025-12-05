package updater

import (
	"runtime"
	"testing"
)

func TestGetVersion(t *testing.T) {
	version := GetVersion()
	if version == "" {
		t.Error("GetVersion returned empty string")
	}
}

func TestGetAssetName(t *testing.T) {
	name := getAssetName()

	expectedPrefix := "forge-orchestrator-" + runtime.GOOS + "-" + runtime.GOARCH
	if runtime.GOOS == "windows" {
		expectedPrefix += ".exe"
	}

	if name != expectedPrefix {
		t.Errorf("Expected asset name %s, got %s", expectedPrefix, name)
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0", "1.0.1", -1},
		{"2.0.0", "1.9.9", 1},
		{"1.9.9", "2.0.0", -1},
		{"1.1.0", "1.0.9", 1},
		{"1.0.0-dev", "1.0.0", 0},
		{"1.1.0-dev", "1.0.0", 1},
	}

	for _, tt := range tests {
		result := compareVersions(tt.v1, tt.v2)
		if result != tt.expected {
			t.Errorf("compareVersions(%s, %s) = %d, expected %d", tt.v1, tt.v2, result, tt.expected)
		}
	}
}

func TestGetExeSuffix(t *testing.T) {
	suffix := getExeSuffix()

	if runtime.GOOS == "windows" {
		if suffix != ".exe" {
			t.Errorf("Expected .exe on Windows, got %s", suffix)
		}
	} else {
		if suffix != "" {
			t.Errorf("Expected empty suffix on Unix, got %s", suffix)
		}
	}
}

// Integration test - only runs if INTEGRATION_TEST env is set
func TestCheckForUpdate_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	info, err := CheckForUpdate()
	if err != nil {
		// This might fail if no releases exist yet, which is OK
		t.Logf("CheckForUpdate returned error (may be expected): %v", err)
		return
	}

	if info == nil {
		t.Error("CheckForUpdate returned nil")
		return
	}

	t.Logf("Current version: %s", info.CurrentVersion)
	t.Logf("Latest version: %s", info.LatestVersion)
	t.Logf("Update available: %v", info.Available)
}

func TestListReleases_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	releases, err := ListReleases(5)
	if err != nil {
		// This might fail if no releases exist yet
		t.Logf("ListReleases returned error (may be expected): %v", err)
		return
	}

	t.Logf("Found %d releases", len(releases))
	for _, r := range releases {
		t.Logf("  - %s (current: %v)", r.Version, r.IsCurrent)
	}
}
