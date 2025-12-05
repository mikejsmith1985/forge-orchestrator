// Package server provides HTTP handlers for the Forge Orchestrator API.
// This test file verifies the api_handlers, specifically the /api/execute endpoint.
package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mikejsmith1985/forge-orchestrator/internal/data"
	"github.com/mikejsmith1985/forge-orchestrator/internal/execution"
)

// TestHandleExecuteSuccess verifies that the /api/execute endpoint works correctly
// with the real Executor interface (LocalRunner) as required by Contract 5's
// NO MOCKS rule.
func TestHandleExecuteSuccess(t *testing.T) {
	// Create a temporary database for the server.
	tempDB := "test_api_execute.db"
	defer os.Remove(tempDB)

	db, err := sql.Open("sqlite", tempDB)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize schema.
	if _, err := db.Exec(data.SQLiteSchema); err != nil {
		t.Fatalf("Failed to init schema: %v", err)
	}

	// Create the server.
	srv := NewServer(db)

	// Create the request payload.
	reqBody := ExecuteRequest{
		Command: "echo hello world",
	}
	body, _ := json.Marshal(reqBody)

	// Create the HTTP request.
	req := httptest.NewRequest(http.MethodPost, "/api/execute", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder.
	rr := httptest.NewRecorder()

	// Call the handler directly.
	srv.handleExecute(rr, req)

	// Verify status code.
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Parse the response.
	var resp ExecuteResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify the response contains "Execution Request Received" (Contract 5 requirement).
	if resp.Message != "Execution Request Received" {
		t.Errorf("Expected message 'Execution Request Received', got %q", resp.Message)
	}

	// Verify the command was actually executed (no mocks!).
	expectedOutput := "hello world\n"
	if resp.Stdout != expectedOutput {
		t.Errorf("Expected stdout %q, got %q", expectedOutput, resp.Stdout)
	}

	// Verify success.
	if !resp.Success {
		t.Errorf("Expected Success to be true, got false")
	}

	// Verify exit code is 0.
	if resp.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", resp.ExitCode)
	}
}

// TestHandleExecuteWithFailedCommand verifies that failed commands are handled correctly.
func TestHandleExecuteWithFailedCommand(t *testing.T) {
	tempDB := "test_api_execute_failed.db"
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

	reqBody := ExecuteRequest{
		Command: "exit 1",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/execute", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	srv.handleExecute(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var resp ExecuteResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify exit code is 1.
	if resp.ExitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", resp.ExitCode)
	}

	// Verify success is false.
	if resp.Success {
		t.Errorf("Expected Success to be false for exit 1")
	}
}

// TestHandleExecuteWithEmptyCommand verifies that empty commands are rejected.
func TestHandleExecuteWithEmptyCommand(t *testing.T) {
	tempDB := "test_api_execute_empty.db"
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

	reqBody := ExecuteRequest{
		Command: "",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/execute", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	srv.handleExecute(rr, req)

	// Should return 400 Bad Request.
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}

	var resp ExecuteResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Success {
		t.Error("Expected Success to be false for empty command")
	}
}

// TestHandleExecuteMethodNotAllowed verifies that non-POST requests are rejected.
func TestHandleExecuteMethodNotAllowed(t *testing.T) {
	tempDB := "test_api_execute_method.db"
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

	// Try GET instead of POST.
	req := httptest.NewRequest(http.MethodGet, "/api/execute", nil)
	rr := httptest.NewRecorder()
	srv.handleExecute(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", rr.Code)
	}
}

// TestExecutorInterfaceIsUsed verifies that the Executor interface is properly used.
// This validates Contract 5's requirement that the endpoint calls the Executor interface.
func TestExecutorInterfaceIsUsed(t *testing.T) {
	// Verify that the executor variable is of type execution.Executor (the interface).
	var _ execution.Executor = executor

	// Verify it's initialized (not nil).
	if executor == nil {
		t.Fatal("Executor should be initialized, not nil")
	}

	// Verify we can execute a command through the interface.
	result := executor.Execute(execution.ExecutionContext{
		Command: "echo test",
	})

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}
}

// TestRouterHasExecuteEndpoint verifies that the router includes the /api/execute endpoint.
func TestRouterHasExecuteEndpoint(t *testing.T) {
	tempDB := "test_router_execute.db"
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
	router := srv.RegisterRoutes()

	// Create a request to /api/execute.
	reqBody := ExecuteRequest{
		Command: "echo router test",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/execute", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Should get a 200 OK (not 404).
	if rr.Code == http.StatusNotFound {
		t.Error("Expected /api/execute to be registered, got 404")
	}

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}
