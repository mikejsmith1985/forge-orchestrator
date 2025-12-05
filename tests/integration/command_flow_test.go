package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// TestCommandLifecycle tests the complete command CRUD flow using real API endpoints.
// This is an integration test that verifies the frontend and backend work together.
func TestCommandLifecycle(t *testing.T) {
	ResetDB(t)

	// 1. Verify empty initial state
	t.Run("InitialEmptyState", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/api/commands")
		if err != nil {
			t.Fatalf("Failed to get commands: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var commands []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&commands); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(commands) != 0 {
			t.Errorf("Expected 0 commands, got %d", len(commands))
		}
	})

	// 2. Create a new command
	var createdID string
	t.Run("CreateCommand", func(t *testing.T) {
		cmd := map[string]string{
			"name":        "Integration Test Command",
			"description": "A command created by integration tests",
			"command":     "echo 'hello from integration test'",
		}
		body, _ := json.Marshal(cmd)

		resp, err := http.Post(testServer.URL+"/api/commands", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create command: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, string(bodyBytes))
		}

		var created map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if id, ok := created["id"].(float64); ok {
			createdID = string(rune(int(id)))
			// Store as string for later use
			createdID = "1" // First created command
		}

		if created["name"] != "Integration Test Command" {
			t.Errorf("Expected name 'Integration Test Command', got %v", created["name"])
		}
	})

	// 3. Verify command appears in list
	t.Run("CommandAppearsInList", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/api/commands")
		if err != nil {
			t.Fatalf("Failed to get commands: %v", err)
		}
		defer resp.Body.Close()

		var commands []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&commands); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(commands) != 1 {
			t.Fatalf("Expected 1 command, got %d", len(commands))
		}

		if commands[0]["name"] != "Integration Test Command" {
			t.Errorf("Expected name 'Integration Test Command', got %v", commands[0]["name"])
		}
	})

	// 4. Delete the command
	t.Run("DeleteCommand", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, testServer.URL+"/api/commands/"+createdID, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to delete command: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", resp.StatusCode)
		}
	})

	// 5. Verify command is deleted
	t.Run("CommandDeleted", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/api/commands")
		if err != nil {
			t.Fatalf("Failed to get commands: %v", err)
		}
		defer resp.Body.Close()

		var commands []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&commands); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(commands) != 0 {
			t.Errorf("Expected 0 commands after deletion, got %d", len(commands))
		}
	})
}

// TestFlowLifecycle tests the complete flow CRUD operations.
func TestFlowLifecycle(t *testing.T) {
	ResetDB(t)

	// 1. Create a flow
	var flowID string
	t.Run("CreateFlow", func(t *testing.T) {
		flow := map[string]interface{}{
			"name":        "Integration Test Flow",
			"description": "A flow created by integration tests",
			"data":        `{"nodes":[],"edges":[]}`,
		}
		body, _ := json.Marshal(flow)

		resp, err := http.Post(testServer.URL+"/api/flows", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create flow: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, string(bodyBytes))
		}

		var created map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if id, ok := created["id"].(float64); ok {
			flowID = fmt.Sprintf("%d", int(id))
		} else {
			flowID = "1"
		}
	})

	// 2. Get the flow by ID
	t.Run("GetFlowByID", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/api/flows/" + flowID)
		if err != nil {
			t.Fatalf("Failed to get flow: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var flow map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&flow); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if flow["name"] != "Integration Test Flow" {
			t.Errorf("Expected name 'Integration Test Flow', got %v", flow["name"])
		}
	})

	// 3. Update the flow
	t.Run("UpdateFlow", func(t *testing.T) {
		update := map[string]interface{}{
			"name":        "Updated Integration Flow",
			"description": "Updated description",
			"data":        `{"nodes":[{"id":"1"}],"edges":[]}`,
		}
		body, _ := json.Marshal(update)

		req, _ := http.NewRequest(http.MethodPut, testServer.URL+"/api/flows/"+flowID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to update flow: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	// 4. Verify the update
	t.Run("VerifyUpdate", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/api/flows/" + flowID)
		if err != nil {
			t.Fatalf("Failed to get flow: %v", err)
		}
		defer resp.Body.Close()

		var flow map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&flow); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if flow["name"] != "Updated Integration Flow" {
			t.Errorf("Expected updated name, got %v", flow["name"])
		}
	})

	// 5. Delete the flow
	t.Run("DeleteFlow", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, testServer.URL+"/api/flows/"+flowID, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to delete flow: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", resp.StatusCode)
		}
	})
}
