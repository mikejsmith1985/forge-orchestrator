package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mikejsmith1985/forge-orchestrator/internal/security"
	"github.com/zalando/go-keyring"
)

func TestKeyringHandlers(t *testing.T) {
	// Initialize mock keyring
	keyring.MockInit()

	s := &Server{}
	// We don't need a full DB or Gateway for these tests as they only touch the keyring

	t.Run("SetAPIKey", func(t *testing.T) {
		payload := map[string]string{
			"provider": "TestProvider",
			"key":      "sk-test-key",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/keys", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		s.handleSetAPIKey(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Verify it was stored
		key, err := security.GetAPIKey("TestProvider")
		if err != nil || key != "sk-test-key" {
			t.Errorf("Failed to retrieve key from keyring: %v", err)
		}
	})

	t.Run("GetAPIKeyStatus", func(t *testing.T) {
		// Ensure we have a key set
		security.SetAPIKey("Anthropic", "sk-anthropic")

		req := httptest.NewRequest("GET", "/api/keys/status", nil)
		w := httptest.NewRecorder()

		s.handleGetAPIKeyStatus(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var status map[string]bool
		json.NewDecoder(w.Body).Decode(&status)

		if !status["Anthropic"] {
			t.Error("Expected Anthropic to be true")
		}
		if status["OpenAI"] {
			// We didn't set OpenAI, so it should be false (or not present, but our logic sets it to false)
			t.Error("Expected OpenAI to be false")
		}
	})

	t.Run("DeleteAPIKey", func(t *testing.T) {
		security.SetAPIKey("ToDelete", "sk-delete")

		req := httptest.NewRequest("DELETE", "/api/keys/ToDelete", nil)
		req.SetPathValue("provider", "ToDelete") // Manually set path value for Go 1.22+ mux
		w := httptest.NewRecorder()

		s.handleDeleteAPIKey(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", w.Code)
		}

		// Verify it's gone
		_, err := security.GetAPIKey("ToDelete")
		if err != keyring.ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})
}
