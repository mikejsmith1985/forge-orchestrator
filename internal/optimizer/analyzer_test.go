package optimizer

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates an in-memory SQLite database and initializes the schema for testing.
func setupTestDB(t *testing.T) *sql.DB {
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

	return db
}

func TestAnalyzeLedger(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	// 1. High-cost model usage (GPT-4)
	// 2. Long prompts
	// 3. Failed executions
	queries := []string{
		// High cost flow
		`INSERT INTO token_ledger (flow_id, model_used, agent_role, prompt_hash, input_tokens, output_tokens, total_cost_usd, latency_ms, status) 
		 VALUES ('flow_expensive', 'gpt-4', 'coder', 'hash1', 1000, 500, 0.09, 1000, 'SUCCESS')`,
		`INSERT INTO token_ledger (flow_id, model_used, agent_role, prompt_hash, input_tokens, output_tokens, total_cost_usd, latency_ms, status) 
		 VALUES ('flow_expensive', 'gpt-4', 'coder', 'hash2', 1000, 500, 0.09, 1000, 'SUCCESS')`,

		// Long prompt flow
		`INSERT INTO token_ledger (flow_id, model_used, agent_role, prompt_hash, input_tokens, output_tokens, total_cost_usd, latency_ms, status) 
		 VALUES ('flow_long', 'gpt-3.5-turbo', 'writer', 'hash3', 3000, 100, 0.01, 500, 'SUCCESS')`,
		`INSERT INTO token_ledger (flow_id, model_used, agent_role, prompt_hash, input_tokens, output_tokens, total_cost_usd, latency_ms, status) 
		 VALUES ('flow_long', 'gpt-3.5-turbo', 'writer', 'hash4', 3500, 100, 0.012, 600, 'SUCCESS')`,

		// Failed flow
		`INSERT INTO token_ledger (flow_id, model_used, agent_role, prompt_hash, input_tokens, output_tokens, total_cost_usd, latency_ms, status) 
		 VALUES ('flow_fail', 'gpt-3.5-turbo', 'tester', 'hash5', 500, 50, 0.002, 200, 'FAILED')`,
		`INSERT INTO token_ledger (flow_id, model_used, agent_role, prompt_hash, input_tokens, output_tokens, total_cost_usd, latency_ms, status) 
		 VALUES ('flow_fail', 'gpt-3.5-turbo', 'tester', 'hash6', 500, 50, 0.002, 200, 'FAILED')`,
	}

	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}

	// Run analysis
	suggestions, err := AnalyzeLedger(db)
	if err != nil {
		t.Fatalf("AnalyzeLedger failed: %v", err)
	}

	// Verify results
	if len(suggestions) != 3 {
		t.Errorf("Expected 3 suggestions, got %d", len(suggestions))
	}

	foundHighCost := false
	foundLongPrompt := false
	foundRetry := false

	for _, s := range suggestions {
		switch s.Type {
		case "model_switch":
			foundHighCost = true
			if s.TargetFlowID != "flow_expensive" {
				t.Errorf("Expected model_switch for flow_expensive, got %s", s.TargetFlowID)
			}
		case "prompt_optimization":
			foundLongPrompt = true
			if s.TargetFlowID != "flow_long" {
				t.Errorf("Expected prompt_optimization for flow_long, got %s", s.TargetFlowID)
			}
		case "retry_strategy":
			foundRetry = true
			if s.TargetFlowID != "flow_fail" {
				t.Errorf("Expected retry_strategy for flow_fail, got %s", s.TargetFlowID)
			}
		}
	}

	if !foundHighCost {
		t.Error("Missing model_switch suggestion")
	}
	if !foundLongPrompt {
		t.Error("Missing prompt_optimization suggestion")
	}
	if !foundRetry {
		t.Error("Missing retry_strategy suggestion")
	}
}
