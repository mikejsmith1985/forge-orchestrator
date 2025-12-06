package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mikejsmith1985/forge-orchestrator/internal/data"
)

func TestWelcomeHandlers(t *testing.T) {
	// Create a temporary database
	tempDB := "test_welcome.db"
	defer os.Remove(tempDB)

	db, err := sql.Open("sqlite", tempDB)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(data.SQLiteSchema); err != nil {
		t.Fatalf("Failed to init schema: %v", err)
	}

	srv := NewServer(db)

	// Reset welcome state before each test run
	ResetWelcomeState()

	t.Run("GET /api/welcome returns welcome status", func(t *testing.T) {
		// Reset state for this specific test
		ResetWelcomeState()
		
		req := httptest.NewRequest("GET", "/api/welcome", nil)
		w := httptest.NewRecorder()

		srv.handleGetWelcome(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response struct {
			Shown          bool   `json:"shown"`
			CurrentVersion string `json:"currentVersion"`
			LastVersion    string `json:"lastVersion"`
		}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// First time should not be shown
		if response.Shown {
			t.Error("Expected shown=false for first time user")
		}
	})

	t.Run("POST /api/welcome marks welcome as shown", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/welcome", nil)
		w := httptest.NewRecorder()

		srv.handleMarkWelcomeShown(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Status != "ok" {
			t.Errorf("Expected status 'ok', got '%s'", response.Status)
		}
	})

	t.Run("GET /api/welcome returns shown=true after POST", func(t *testing.T) {
		// First mark as shown
		postReq := httptest.NewRequest("POST", "/api/welcome", nil)
		postW := httptest.NewRecorder()
		srv.handleMarkWelcomeShown(postW, postReq)

		// Then check status
		req := httptest.NewRequest("GET", "/api/welcome", nil)
		w := httptest.NewRecorder()
		srv.handleGetWelcome(w, req)

		var response struct {
			Shown bool `json:"shown"`
		}
		json.NewDecoder(w.Body).Decode(&response)

		if !response.Shown {
			t.Error("Expected shown=true after marking welcome as shown")
		}
	})
}
