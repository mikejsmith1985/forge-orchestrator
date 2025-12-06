package server

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/mikejsmith1985/forge-orchestrator/internal/config"
)

// WelcomeState tracks whether the welcome modal has been shown for a version.
type WelcomeState struct {
	LastShownVersion string `json:"lastShownVersion"`
}

var (
	welcomeState   *WelcomeState
	welcomeStateMu sync.RWMutex
)

// WelcomeResponse represents the response for the welcome status endpoint.
type WelcomeResponse struct {
	Shown          bool   `json:"shown"`
	CurrentVersion string `json:"currentVersion"`
	LastVersion    string `json:"lastVersion"`
}

// getWelcomeStatePath returns the path to the welcome state file.
func getWelcomeStatePath() (string, error) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "welcome_state.json"), nil
}

// loadWelcomeState loads the welcome state from disk.
func loadWelcomeState() (*WelcomeState, error) {
	welcomeStateMu.Lock()
	defer welcomeStateMu.Unlock()

	if welcomeState != nil {
		return welcomeState, nil
	}

	statePath, err := getWelcomeStatePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			welcomeState = &WelcomeState{}
			return welcomeState, nil
		}
		return nil, err
	}

	var state WelcomeState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	welcomeState = &state
	return welcomeState, nil
}

// saveWelcomeState saves the welcome state to disk.
func saveWelcomeState(state *WelcomeState) error {
	welcomeStateMu.Lock()
	defer welcomeStateMu.Unlock()

	statePath, err := getWelcomeStatePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(statePath, data, 0644); err != nil {
		return err
	}

	welcomeState = state
	return nil
}

// handleGetWelcome returns the welcome status.
func (s *Server) handleGetWelcome(w http.ResponseWriter, r *http.Request) {
	state, err := loadWelcomeState()
	if err != nil {
		http.Error(w, "Failed to load welcome state: "+err.Error(), http.StatusInternalServerError)
		return
	}

	currentVersion := s.getVersion()
	
	// Show welcome if:
	// 1. Never shown before (lastShownVersion is empty)
	// 2. Current version is different from last shown version
	shown := state.LastShownVersion != "" && state.LastShownVersion == currentVersion

	response := WelcomeResponse{
		Shown:          shown,
		CurrentVersion: currentVersion,
		LastVersion:    state.LastShownVersion,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleMarkWelcomeShown marks the welcome modal as shown for the current version.
func (s *Server) handleMarkWelcomeShown(w http.ResponseWriter, r *http.Request) {
	currentVersion := s.getVersion()

	state := &WelcomeState{
		LastShownVersion: currentVersion,
	}

	if err := saveWelcomeState(state); err != nil {
		http.Error(w, "Failed to save welcome state: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// getVersion returns the current application version.
func (s *Server) getVersion() string {
	// This should be set during build time or read from a version file
	// For now, return a default version
	return "1.1.0"
}

// ResetWelcomeState is used for testing to reset the in-memory state.
// It also removes the welcome state file if it exists.
func ResetWelcomeState() {
	welcomeStateMu.Lock()
	defer welcomeStateMu.Unlock()
	welcomeState = nil

	// Also remove the file for clean tests
	if statePath, err := getWelcomeStatePath(); err == nil {
		os.Remove(statePath)
	}
}
