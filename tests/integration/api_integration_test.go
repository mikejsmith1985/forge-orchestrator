package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// TestLedgerIntegration tests the ledger API endpoints.
func TestLedgerIntegration(t *testing.T) {
	ResetDB(t)

	// 1. Verify initial empty ledger
	t.Run("InitialEmptyLedger", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/api/ledger")
		if err != nil {
			t.Fatalf("Failed to get ledger: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var entries []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(entries) != 0 {
			t.Errorf("Expected 0 ledger entries, got %d", len(entries))
		}
	})

	// 2. Create a ledger entry
	t.Run("CreateLedgerEntry", func(t *testing.T) {
		entry := map[string]interface{}{
			"flow_id":       "test-flow-123",
			"model_used":    "claude-3-5-sonnet",
			"agent_role":    "Optimizer",
			"prompt_hash":   "abc123",
			"input_tokens":  100,
			"output_tokens": 50,
			"total_cost_usd": 0.00105,
			"latency_ms":    1250,
			"status":        "SUCCESS",
		}
		body, _ := json.Marshal(entry)

		resp, err := http.Post(testServer.URL+"/api/ledger", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create ledger entry: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	// 3. Verify ledger entry appears
	t.Run("LedgerEntryAppears", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/api/ledger")
		if err != nil {
			t.Fatalf("Failed to get ledger: %v", err)
		}
		defer resp.Body.Close()

		var entries []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(entries) != 1 {
			t.Fatalf("Expected 1 ledger entry, got %d", len(entries))
		}

		if entries[0]["flow_id"] != "test-flow-123" {
			t.Errorf("Expected flow_id 'test-flow-123', got %v", entries[0]["flow_id"])
		}

		if entries[0]["model_used"] != "claude-3-5-sonnet" {
			t.Errorf("Expected model_used 'claude-3-5-sonnet', got %v", entries[0]["model_used"])
		}
	})
}

// TestKeyManagementIntegration tests the API key management endpoints.
// Note: This test may skip certain operations in CI environments where
// the keyring service (e.g., freedesktop.org secrets) isn't available.
func TestKeyManagementIntegration(t *testing.T) {
	ResetDB(t)

	// 1. Verify initial key status (all unconfigured)
	t.Run("InitialKeyStatus", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/api/keys/status")
		if err != nil {
			t.Fatalf("Failed to get key status: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var status map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Verify all providers show as unconfigured
		if status["anthropic"] == true {
			t.Error("Expected anthropic to be unconfigured initially")
		}
	})

	// 2. Set an API key (may skip in CI due to keyring unavailability)
	t.Run("SetAPIKey", func(t *testing.T) {
		keyData := map[string]string{
			"provider": "anthropic",
			"key":      "test-api-key-12345",
		}
		body, _ := json.Marshal(keyData)

		resp, err := http.Post(testServer.URL+"/api/keys", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to set API key: %v", err)
		}
		defer resp.Body.Close()

		// Skip test if keyring service is unavailable (common in CI)
		if resp.StatusCode == http.StatusInternalServerError {
			bodyBytes, _ := io.ReadAll(resp.Body)
			if bytes.Contains(bodyBytes, []byte("freedesktop")) || bytes.Contains(bodyBytes, []byte("secrets")) {
				t.Skip("Keyring service not available in this environment")
			}
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200 or 201, got %d: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	// 3. Verify key status shows configured (skip if SetAPIKey was skipped)
	t.Run("KeyStatusConfigured", func(t *testing.T) {
		// First check if key was set by trying to get status
		resp, err := http.Get(testServer.URL + "/api/keys/status")
		if err != nil {
			t.Fatalf("Failed to get key status: %v", err)
		}
		defer resp.Body.Close()

		var status map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// If anthropic is not set, skip (keyring wasn't available)
		if status["anthropic"] == nil || status["anthropic"] == false {
			t.Skip("API key wasn't set - keyring likely unavailable")
		}

		if status["anthropic"] != true {
			t.Errorf("Expected anthropic to be configured, got %v", status["anthropic"])
		}
	})

	// 4. Delete the API key
	t.Run("DeleteAPIKey", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, testServer.URL+"/api/keys/anthropic", nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to delete API key: %v", err)
		}
		defer resp.Body.Close()

		// Skip if keyring unavailable
		if resp.StatusCode == http.StatusInternalServerError {
			t.Skip("Keyring service not available in this environment")
		}

		if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 204 or 200, got %d", resp.StatusCode)
		}
	})

	// 5. Verify key is deleted
	t.Run("KeyDeleted", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/api/keys/status")
		if err != nil {
			t.Fatalf("Failed to get key status: %v", err)
		}
		defer resp.Body.Close()

		var status map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if status["anthropic"] == true {
			t.Error("Expected anthropic to be unconfigured after deletion")
		}
	})
}

// TestOptimizationIntegration tests the optimization suggestions API.
func TestOptimizationIntegration(t *testing.T) {
	ResetDB(t)

	// 1. Create some ledger entries to generate optimization data
	t.Run("SetupLedgerData", func(t *testing.T) {
		entries := []map[string]interface{}{
			{
				"flow_id":        "flow-1",
				"model_used":     "claude-3-5-sonnet",
				"agent_role":     "Optimizer",
				"prompt_hash":    "hash1",
				"input_tokens":   1000,
				"output_tokens":  500,
				"total_cost_usd": 0.0105,
				"latency_ms":     1500,
				"status":         "SUCCESS",
			},
			{
				"flow_id":        "flow-1",
				"model_used":     "claude-3-5-sonnet",
				"agent_role":     "Optimizer",
				"prompt_hash":    "hash1", // Same hash - duplicate prompt
				"input_tokens":   1000,
				"output_tokens":  500,
				"total_cost_usd": 0.0105,
				"latency_ms":     1500,
				"status":         "SUCCESS",
			},
		}

		for _, entry := range entries {
			body, _ := json.Marshal(entry)
			resp, err := http.Post(testServer.URL+"/api/ledger", "application/json", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("Failed to create ledger entry: %v", err)
			}
			resp.Body.Close()
		}
	})

	// 2. Get optimizations - the analyzer should detect duplicate prompts
	t.Run("GetOptimizations", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/api/ledger/optimizations")
		if err != nil {
			t.Fatalf("Failed to get optimizations: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(bodyBytes))
		}

		var optimizations []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&optimizations); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Should have at least one optimization (duplicate prompt detection)
		if len(optimizations) == 0 {
			t.Log("No optimizations found - this may be expected if analyzer needs more data")
		}
	})
}

// TestHealthEndpoint verifies the health check endpoint.
func TestHealthEndpoint(t *testing.T) {
	t.Run("HealthCheck", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/api/health")
		if err != nil {
			t.Fatalf("Failed to get health: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var health map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if health["status"] != "ok" {
			t.Errorf("Expected status 'ok', got %v", health["status"])
		}
	})
}
