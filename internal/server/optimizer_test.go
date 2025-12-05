package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestServer(t *testing.T) *Server {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	schema := `
	CREATE TABLE token_ledger (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		flow_id TEXT NOT NULL,
		model_used TEXT NOT NULL,
		agent_role TEXT NOT NULL,
		prompt_hash TEXT NOT NULL,
		input_tokens INTEGER NOT NULL,
		output_tokens INTEGER NOT NULL,
		total_cost_usd REAL NOT NULL,
		latency_ms INTEGER NOT NULL,
		status TEXT NOT NULL,
		error_message TEXT
	);

	CREATE TABLE IF NOT EXISTS optimization_suggestions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		estimated_savings REAL NOT NULL,
		savings_unit TEXT NOT NULL,
		target_flow_id TEXT,
		target_command_id INTEGER,
		apply_action TEXT NOT NULL,
		status TEXT DEFAULT 'pending',
		applied_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS forge_flows (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		data TEXT NOT NULL,
		status TEXT DEFAULT 'draft',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	return &Server{db: db}
}

func TestHandleGetOptimizations(t *testing.T) {
	server := setupTestServer(t)
	defer server.db.Close()

	// Insert data that should trigger an optimization
	_, err := server.db.Exec(`
		INSERT INTO token_ledger (flow_id, model_used, agent_role, prompt_hash, input_tokens, output_tokens, total_cost_usd, latency_ms, status) 
		VALUES ('flow_expensive', 'gpt-4', 'coder', 'hash1', 1000, 500, 0.09, 1000, 'SUCCESS')
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}
	_, err = server.db.Exec(`
		INSERT INTO token_ledger (flow_id, model_used, agent_role, prompt_hash, input_tokens, output_tokens, total_cost_usd, latency_ms, status) 
		VALUES ('flow_expensive', 'gpt-4', 'coder', 'hash2', 1000, 500, 0.09, 1000, 'SUCCESS')
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	req, _ := http.NewRequest("GET", "/api/ledger/optimizations", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(server.handleGetOptimizations)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var suggestions []map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&suggestions); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(suggestions) == 0 {
		t.Error("Expected suggestions, got none")
	}
}

func TestHandleApplyOptimization(t *testing.T) {
	server := setupTestServer(t)
	defer server.db.Close()

	// Create a test flow that we can apply optimization to
	_, err := server.db.Exec(`
		INSERT INTO forge_flows (name, data, status) 
		VALUES ('Test Flow', '{"nodes":[{"id":"1","type":"input","data":{"label":"Start","provider":"gpt-4"}}],"edges":[]}', 'active')
	`)
	if err != nil {
		t.Fatalf("Failed to create test flow: %v", err)
	}

	// Create a suggestion in the database
	_, err = server.db.Exec(`
		INSERT INTO optimization_suggestions (type, title, description, estimated_savings, savings_unit, target_flow_id, apply_action, status)
		VALUES ('model_switch', 'Switch model', 'Test', 0.05, 'USD', '1', '{"action":"switch_model","from_model":"gpt-4","to_model":"gpt-3.5-turbo","flow_id":"1"}', 'pending')
	`)
	if err != nil {
		t.Fatalf("Failed to create suggestion: %v", err)
	}

	mux := server.RegisterRoutes()
	req, _ := http.NewRequest("POST", "/api/ledger/optimizations/1/apply", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v. Body: %s", status, http.StatusOK, rr.Body.String())
	}

	// Verify suggestion is marked as applied
	var suggestionStatus string
	err = server.db.QueryRow("SELECT status FROM optimization_suggestions WHERE id = 1").Scan(&suggestionStatus)
	if err != nil {
		t.Fatalf("Failed to query suggestion: %v", err)
	}

	if suggestionStatus != "applied" {
		t.Errorf("Expected suggestion status 'applied', got '%s'", suggestionStatus)
	}
}
