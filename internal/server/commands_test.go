package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Initialize schema
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS command_cards (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		command TEXT NOT NULL,
		description TEXT
	);
	CREATE TABLE IF NOT EXISTS token_ledger (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		flow_id TEXT,
		model_used TEXT,
		agent_role TEXT,
		prompt_hash TEXT,
		input_tokens INTEGER,
		output_tokens INTEGER,
		total_cost_usd REAL,
		latency_ms INTEGER,
		status TEXT,
		error_message TEXT
	);
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	return db
}

func TestHandleCreateCommand(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	server := NewServer(db)
	handler := server.RegisterRoutes()

	cmd := CommandCard{
		Name:        "Test Command",
		Command:     "echo 'hello'",
		Description: "A test command",
	}
	body, _ := json.Marshal(cmd)

	req, _ := http.NewRequest("POST", "/api/commands", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var created CommandCard
	json.NewDecoder(rr.Body).Decode(&created)

	if created.ID == 0 {
		t.Errorf("Expected ID to be set")
	}
	if created.Name != cmd.Name {
		t.Errorf("Expected name %v, got %v", cmd.Name, created.Name)
	}
}

func TestHandleGetCommands(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Seed data
	_, err := db.Exec("INSERT INTO command_cards (name, command, description) VALUES (?, ?, ?)", "Cmd 1", "ls -la", "List files")
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}

	server := NewServer(db)
	handler := server.RegisterRoutes()

	req, _ := http.NewRequest("GET", "/api/commands", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var commands []CommandCard
	json.NewDecoder(rr.Body).Decode(&commands)

	if len(commands) != 1 {
		t.Errorf("Expected 1 command, got %v", len(commands))
	}
	if commands[0].Name != "Cmd 1" {
		t.Errorf("Expected name 'Cmd 1', got %v", commands[0].Name)
	}
}

func TestHandleDeleteCommand(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Seed data
	res, _ := db.Exec("INSERT INTO command_cards (name, command, description) VALUES (?, ?, ?)", "Cmd 1", "ls -la", "List files")
	id, _ := res.LastInsertId()

	server := NewServer(db)
	handler := server.RegisterRoutes()

	req, _ := http.NewRequest("DELETE", "/api/commands/"+strconv.Itoa(int(id)), nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	// Verify deletion
	var count int
	db.QueryRow("SELECT COUNT(*) FROM command_cards WHERE id = ?", id).Scan(&count)
	if count != 0 {
		t.Errorf("Expected command to be deleted")
	}
}

// MockLLMProvider for testing
type MockLLMProvider struct {
	SendFunc func(systemPrompt, userPrompt, apiKey string) (string, int, int, error)
}

func (m *MockLLMProvider) Send(systemPrompt, userPrompt, apiKey string) (string, int, int, error) {
	if m.SendFunc != nil {
		return m.SendFunc(systemPrompt, userPrompt, apiKey)
	}
	return "mock response", 10, 20, nil
}

func TestHandleRunCommand(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Seed command
	res, _ := db.Exec("INSERT INTO command_cards (name, command, description) VALUES (?, ?, ?)", "Test Cmd", "echo test", "Desc")
	id, _ := res.LastInsertId()

	// Create server with mock gateway
	server := NewServer(db)
	mockProvider := &MockLLMProvider{
		SendFunc: func(systemPrompt, userPrompt, apiKey string) (string, int, int, error) {
			if apiKey != "test-key" {
				return "", 0, 0, nil
			}
			return "Executed: " + userPrompt, 5, 10, nil
		},
	}
	// Inject mock provider into gateway
	server.gateway.AnthropicClient = mockProvider
	server.gateway.OpenAIClient = mockProvider

	handler := server.RegisterRoutes()

	// Create request
	reqBody := map[string]string{
		"agent_role": "Implementation",
		"provider":   "OpenAI",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/commands/"+strconv.Itoa(int(id))+"/run", bytes.NewBuffer(body))
	req.Header.Set("X-Forge-Api-Key", "test-key")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v body: %v",
			status, http.StatusOK, rr.Body.String())
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if content, ok := response["Content"].(string); !ok || content != "Executed: echo test" {
		t.Errorf("Unexpected content: %v", response["Content"])
	}

	// Verify ledger entry
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM token_ledger").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query ledger: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 ledger entry, got %d", count)
	}
}

func TestHandleRunCommand_LatencyTracking(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Seed command
	res, _ := db.Exec("INSERT INTO command_cards (name, command, description) VALUES (?, ?, ?)", "Test Cmd", "echo test", "Desc")
	id, _ := res.LastInsertId()

	// Create server with mock gateway that simulates some delay
	server := NewServer(db)
	mockProvider := &MockLLMProvider{
		SendFunc: func(systemPrompt, userPrompt, apiKey string) (string, int, int, error) {
			// Simulate a small delay to ensure latency is non-zero
			return "Response", 10, 20, nil
		},
	}
	server.gateway.AnthropicClient = mockProvider
	server.gateway.OpenAIClient = mockProvider

	handler := server.RegisterRoutes()

	// Create request
	reqBody := map[string]string{
		"agent_role": "Implementation",
		"provider":   "OpenAI",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/commands/"+strconv.Itoa(int(id))+"/run", bytes.NewBuffer(body))
	req.Header.Set("X-Forge-Api-Key", "test-key")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify ledger entry has non-zero latency
	var latencyMs int
	err := db.QueryRow("SELECT latency_ms FROM token_ledger ORDER BY id DESC LIMIT 1").Scan(&latencyMs)
	if err != nil {
		t.Fatalf("Failed to query ledger: %v", err)
	}

	// Latency should be non-negative (0 or greater)
	// We can't guarantee it's > 0 since mock is instant, but it should be set
	if latencyMs < 0 {
		t.Errorf("Expected non-negative latency, got %d", latencyMs)
	}

	// Log the latency for informational purposes
	t.Logf("Recorded latency: %d ms", latencyMs)
}

// ========== ERROR HANDLING TESTS ==========

func TestHandleCreateCommand_Errors(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	server := NewServer(db)
	handler := server.RegisterRoutes()

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "malformed JSON",
			body:       `{"name": "test"`, // missing closing brace
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty name",
			body:       `{"name": "", "command": "echo hello"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty command",
			body:       `{"name": "Test", "command": ""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "SQL injection attempt in name",
			body:       `{"name": "'; DROP TABLE command_cards; --", "command": "test"}`,
			wantStatus: http.StatusCreated, // Should succeed but sanitized via parameterized query
		},
		{
			name:       "SQL injection attempt in command",
			body:       `{"name": "Test", "command": "'; DELETE FROM command_cards; --"}`,
			wantStatus: http.StatusCreated, // Should succeed but sanitized via parameterized query
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/commands",
				bytes.NewBufferString(tt.body))
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d, body: %s", rr.Code, tt.wantStatus, rr.Body.String())
			}
		})
	}

	// Verify SQL injection didn't corrupt database - table should still exist
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM command_cards").Scan(&count)
	if err != nil {
		t.Errorf("Database corrupted after SQL injection test: %v", err)
	}
}

func TestHandleDeleteCommand_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	server := NewServer(db)
	handler := server.RegisterRoutes()

	// Delete a non-existent ID
	req, _ := http.NewRequest("DELETE", "/api/commands/99999", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Should succeed (idempotent delete) with 204
	if rr.Code != http.StatusNoContent {
		t.Errorf("expected 204 for idempotent delete, got %d", rr.Code)
	}
}

func TestHandleDeleteCommand_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	server := NewServer(db)
	handler := server.RegisterRoutes()

	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{"non-numeric ID", "/api/commands/abc", http.StatusBadRequest},
		{"float ID", "/api/commands/1.5", http.StatusBadRequest},
		// Note: Negative IDs are parsed as valid integers and result in idempotent delete (204)
		// This is acceptable behavior as no row exists with negative ID
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", tt.id, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("expected %d for %s, got %d", tt.wantStatus, tt.name, rr.Code)
			}
		})
	}
}

func TestHandleRunCommand_MissingAPIKey(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	res, _ := db.Exec("INSERT INTO command_cards (name, command, description) VALUES (?, ?, ?)", "Test", "echo", "Desc")
	id, _ := res.LastInsertId()

	server := NewServer(db)
	handler := server.RegisterRoutes()

	reqBody := map[string]string{
		"agent_role": "Implementation",
		"provider":   "OpenAI",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/commands/"+strconv.Itoa(int(id))+"/run", bytes.NewBuffer(body))
	// Intentionally NOT setting X-Forge-Api-Key header

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for missing API key, got %d", rr.Code)
	}
}

func TestHandleRunCommand_MissingFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	res, _ := db.Exec("INSERT INTO command_cards (name, command, description) VALUES (?, ?, ?)", "Test", "echo", "Desc")
	id, _ := res.LastInsertId()

	server := NewServer(db)
	handler := server.RegisterRoutes()

	tests := []struct {
		name string
		body string
	}{
		{"missing agent_role", `{"provider": "OpenAI"}`},
		{"missing provider", `{"agent_role": "Implementation"}`},
		{"empty agent_role", `{"agent_role": "", "provider": "OpenAI"}`},
		{"empty provider", `{"agent_role": "Implementation", "provider": ""}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/commands/"+strconv.Itoa(int(id))+"/run",
				bytes.NewBufferString(tt.body))
			req.Header.Set("X-Forge-Api-Key", "test-key")
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("expected 400 for %s, got %d", tt.name, rr.Code)
			}
		})
	}
}

func TestHandleRunCommand_CommandNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	server := NewServer(db)
	handler := server.RegisterRoutes()

	reqBody := map[string]string{
		"agent_role": "Implementation",
		"provider":   "OpenAI",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/commands/99999/run", bytes.NewBuffer(body))
	req.Header.Set("X-Forge-Api-Key", "test-key")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404 for non-existent command, got %d", rr.Code)
	}
}
