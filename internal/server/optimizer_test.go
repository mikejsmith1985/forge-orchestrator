package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestServer(t *testing.T) *Server {
	db, err := sql.Open("sqlite3", ":memory:")
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

	// Note: Testing PathValue in unit tests requires Go 1.22+ and using the actual ServeMux
	// or mocking the request context. Since we are using standard http.HandlerFunc here,
	// PathValue won't be populated unless we route through the mux.

	mux := server.RegisterRoutes()
	req, _ := http.NewRequest("POST", "/api/ledger/optimizations/1/apply", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify ledger update
	var count int
	err := server.db.QueryRow("SELECT COUNT(*) FROM token_ledger WHERE status = 'OPTIMIZED'").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query ledger: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 optimized entry, got %d", count)
	}
}
