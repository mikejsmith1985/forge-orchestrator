package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

	t.Run("GetAPIKeyStatus_ReturnsCorrectFormat", func(t *testing.T) {
		// Ensure we have a key set
		security.SetAPIKey("Anthropic", "sk-anthropic")

		req := httptest.NewRequest("GET", "/api/keys/status", nil)
		w := httptest.NewRecorder()

		s.handleGetAPIKeyStatus(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Frontend expects: {"keys": [{"provider": "anthropic", "isSet": true}]}
		var response struct {
			Keys []struct {
				Provider string `json:"provider"`
				IsSet    bool   `json:"isSet"`
			} `json:"keys"`
		}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Keys == nil {
			t.Fatal("Expected 'keys' array in response, got nil")
		}

		// Find Anthropic in the keys array
		foundAnthropic := false
		for _, key := range response.Keys {
			if key.Provider == "anthropic" {
				foundAnthropic = true
				if !key.IsSet {
					t.Error("Expected Anthropic isSet to be true")
				}
			}
		}
		if !foundAnthropic {
			t.Error("Expected to find 'anthropic' provider in keys array")
		}
	})

	t.Run("GetAPIKeyStatus_ProviderNamesAreLowercase", func(t *testing.T) {
		// Clear and set up fresh state
		keyring.MockInit()
		security.SetAPIKey("OpenAI", "sk-openai")

		req := httptest.NewRequest("GET", "/api/keys/status", nil)
		w := httptest.NewRecorder()

		s.handleGetAPIKeyStatus(w, req)

		var response struct {
			Keys []struct {
				Provider string `json:"provider"`
				IsSet    bool   `json:"isSet"`
			} `json:"keys"`
		}
		json.NewDecoder(w.Body).Decode(&response)

		for _, key := range response.Keys {
			// Provider names should be lowercase for frontend consistency
			if key.Provider != strings.ToLower(key.Provider) {
				t.Errorf("Provider name should be lowercase, got: %s", key.Provider)
			}
		}
	})

	t.Run("SetAPIKey_ReturnsSuccessMessage", func(t *testing.T) {
		keyring.MockInit()

		payload := map[string]string{
			"provider": "anthropic",
			"key":      "sk-ant-test-key",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/keys", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		s.handleSetAPIKey(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Response should include success confirmation
		var response struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Status != "ok" {
			t.Errorf("Expected status 'ok', got '%s'", response.Status)
		}
		if response.Message == "" {
			t.Error("Expected a success message in response")
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
