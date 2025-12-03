package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikejsmith1985/forge-orchestrator/internal/data"
)

func TestHandleEstimateTokens(t *testing.T) {
	// Server doesn't need DB for this test
	s := &Server{}
	handler := s.RegisterRoutes()

	payload := map[string]string{"text": "hello world"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/tokens/estimate", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]int
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// "hello world" is 11 chars. 11 / 4 = 2 (integer division)
	if count, ok := response["count"]; !ok || count != 2 {
		t.Errorf("handler returned wrong token count: got %v want %v", count, 2)
	}
}

func TestHandleCreateLedgerEntry(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Initialize schema
	_, err = db.Exec(data.SQLiteSchema)
	if err != nil {
		t.Fatal(err)
	}

	s := NewServer(db)
	handler := s.RegisterRoutes()

	entry := map[string]interface{}{
		"flow_id":        "test-flow-123",
		"model_used":     "gpt-4",
		"agent_role":     "developer",
		"prompt_hash":    "abc123hash",
		"input_tokens":   100,
		"output_tokens":  50,
		"total_cost_usd": 0.003,
		"latency_ms":     500,
		"status":         "SUCCESS",
	}
	body, _ := json.Marshal(entry)
	req, _ := http.NewRequest("POST", "/api/ledger", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Verify insertion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM token_ledger WHERE flow_id = ?", "test-flow-123").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1 row in token_ledger, got %d", count)
	}
}

func TestHandleGetLedger(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Initialize schema
	_, err = db.Exec(data.SQLiteSchema)
	if err != nil {
		t.Fatal(err)
	}

	// Insert some test data
	insertQuery := `
		INSERT INTO token_ledger (
			flow_id, model_used, agent_role, prompt_hash, 
			input_tokens, output_tokens, total_cost_usd, 
			latency_ms, status, timestamp
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	// Insert older entry
	_, err = db.Exec(insertQuery, "flow-1", "gpt-4", "coder", "hash1", 10, 10, 0.01, 100, "SUCCESS", "2023-01-01 10:00:00")
	if err != nil {
		t.Fatal(err)
	}
	// Insert newer entry
	_, err = db.Exec(insertQuery, "flow-2", "gpt-4", "architect", "hash2", 20, 20, 0.02, 200, "SUCCESS", "2023-01-01 11:00:00")
	if err != nil {
		t.Fatal(err)
	}

	s := NewServer(db)
	handler := s.RegisterRoutes()

	// Test 1: Get all entries (default limit)
	req, _ := http.NewRequest("GET", "/api/ledger", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var entries []LedgerEntry
	if err := json.Unmarshal(rr.Body.Bytes(), &entries); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}

	// Verify ordering (newest first)
	if entries[0].FlowID != "flow-2" {
		t.Errorf("expected first entry to be flow-2, got %s", entries[0].FlowID)
	}

	// Test 2: Limit
	req, _ = http.NewRequest("GET", "/api/ledger?limit=1", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &entries); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].FlowID != "flow-2" {
		t.Errorf("expected first entry to be flow-2, got %s", entries[0].FlowID)
	}
}
