package optimizer

import (
	"database/sql"
	"encoding/json"
	"testing"

	_ "modernc.org/sqlite"
)

func setupApplierTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Create required tables
	schema := `
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

	return db
}

func TestStoreSuggestion(t *testing.T) {
	db := setupApplierTestDB(t)
	defer db.Close()

	suggestion := Suggestion{
		Type:             "model_switch",
		Title:            "Switch to cheaper model",
		Description:      "Test description",
		EstimatedSavings: 0.05,
		SavingsUnit:      "USD",
		TargetFlowID:     "1",
		ApplyAction:      `{"action":"switch_model","from_model":"gpt-4","to_model":"gpt-3.5-turbo","flow_id":"1"}`,
	}

	id, err := StoreSuggestion(db, suggestion)
	if err != nil {
		t.Fatalf("Failed to store suggestion: %v", err)
	}

	if id == 0 {
		t.Error("Expected non-zero ID")
	}

	// Verify it was stored
	stored, err := GetSuggestionByID(db, int(id))
	if err != nil {
		t.Fatalf("Failed to get suggestion: %v", err)
	}

	if stored.Title != suggestion.Title {
		t.Errorf("Expected title '%s', got '%s'", suggestion.Title, stored.Title)
	}

	if stored.Status != "pending" {
		t.Errorf("Expected status 'pending', got '%s'", stored.Status)
	}
}

func TestApplyModelSwitch(t *testing.T) {
	db := setupApplierTestDB(t)
	defer db.Close()

	// Create a test flow with a node using gpt-4
	flowData := `{
		"nodes": [
			{"id": "1", "type": "input", "data": {"label": "Start", "provider": "gpt-4"}, "position": {"x": 100, "y": 100}},
			{"id": "2", "type": "default", "data": {"label": "Process", "provider": "gpt-4"}, "position": {"x": 100, "y": 200}}
		],
		"edges": [{"id": "e1-2", "source": "1", "target": "2"}]
	}`

	_, err := db.Exec("INSERT INTO forge_flows (name, data, status) VALUES (?, ?, ?)",
		"Test Flow", flowData, "active")
	if err != nil {
		t.Fatalf("Failed to create flow: %v", err)
	}

	// Create and store the suggestion
	applyAction := map[string]interface{}{
		"action":     "switch_model",
		"from_model": "gpt-4",
		"to_model":   "gpt-3.5-turbo",
		"flow_id":    "1",
	}
	applyJSON, _ := json.Marshal(applyAction)

	suggestion := Suggestion{
		Type:             "model_switch",
		Title:            "Switch from gpt-4 to gpt-3.5-turbo",
		Description:      "Test switch",
		EstimatedSavings: 0.05,
		SavingsUnit:      "USD",
		TargetFlowID:     "1",
		ApplyAction:      string(applyJSON),
	}

	id, err := StoreSuggestion(db, suggestion)
	if err != nil {
		t.Fatalf("Failed to store suggestion: %v", err)
	}

	// Apply the optimization
	result, err := ApplyOptimization(db, int(id))
	if err != nil {
		t.Fatalf("Failed to apply optimization: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success, got failure: %s", result.Message)
	}

	// Verify the flow was updated
	var updatedFlowData string
	err = db.QueryRow("SELECT data FROM forge_flows WHERE id = 1").Scan(&updatedFlowData)
	if err != nil {
		t.Fatalf("Failed to get updated flow: %v", err)
	}

	var graph FlowGraph
	json.Unmarshal([]byte(updatedFlowData), &graph)

	for _, node := range graph.Nodes {
		if node.Data.Provider == "gpt-4" {
			t.Errorf("Node %s still using gpt-4, should be gpt-3.5-turbo", node.ID)
		}
	}

	// Verify suggestion is marked as applied
	stored, _ := GetSuggestionByID(db, int(id))
	if stored.Status != "applied" {
		t.Errorf("Expected status 'applied', got '%s'", stored.Status)
	}
}

func TestApplySuggestionAlreadyApplied(t *testing.T) {
	db := setupApplierTestDB(t)
	defer db.Close()

	// Create a flow
	_, err := db.Exec("INSERT INTO forge_flows (name, data, status) VALUES (?, ?, ?)",
		"Test Flow", `{"nodes":[],"edges":[]}`, "active")
	if err != nil {
		t.Fatalf("Failed to create flow: %v", err)
	}

	// Store and apply a suggestion
	suggestion := Suggestion{
		Type:             "model_switch",
		Title:            "Test",
		Description:      "Test",
		EstimatedSavings: 0.01,
		SavingsUnit:      "USD",
		TargetFlowID:     "1",
		ApplyAction:      `{"action":"switch_model","from_model":"gpt-4","to_model":"gpt-3.5-turbo","flow_id":"1"}`,
	}

	id, _ := StoreSuggestion(db, suggestion)

	// Mark as already applied
	db.Exec("UPDATE optimization_suggestions SET status = 'applied' WHERE id = ?", id)

	// Try to apply again
	result, err := ApplyOptimization(db, int(id))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Success {
		t.Error("Expected failure for already applied suggestion")
	}

	if result.Message != "Suggestion has already been applied" {
		t.Errorf("Unexpected message: %s", result.Message)
	}
}

func TestApplyPromptOptimization(t *testing.T) {
	db := setupApplierTestDB(t)
	defer db.Close()

	// Create a flow
	_, err := db.Exec("INSERT INTO forge_flows (name, data, status) VALUES (?, ?, ?)",
		"Test Flow", `{"nodes":[],"edges":[]}`, "active")
	if err != nil {
		t.Fatalf("Failed to create flow: %v", err)
	}

	// Create suggestion for prompt optimization
	applyAction := map[string]interface{}{
		"action":     "optimize_prompt",
		"flow_id":    "1",
		"agent_role": "coder",
	}
	applyJSON, _ := json.Marshal(applyAction)

	suggestion := Suggestion{
		Type:             "prompt_optimization",
		Title:            "Optimize prompts",
		Description:      "Test",
		EstimatedSavings: 500,
		SavingsUnit:      "tokens",
		TargetFlowID:     "1",
		ApplyAction:      string(applyJSON),
	}

	id, _ := StoreSuggestion(db, suggestion)

	// Apply
	result, err := ApplyOptimization(db, int(id))
	if err != nil {
		t.Fatalf("Failed to apply: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success: %s", result.Message)
	}

	// Verify flow status changed
	var status string
	db.QueryRow("SELECT status FROM forge_flows WHERE id = 1").Scan(&status)
	if status != "needs_optimization" {
		t.Errorf("Expected status 'needs_optimization', got '%s'", status)
	}
}

func TestApplyRetryStrategy(t *testing.T) {
	db := setupApplierTestDB(t)
	defer db.Close()

	// Create a flow
	_, err := db.Exec("INSERT INTO forge_flows (name, data, status) VALUES (?, ?, ?)",
		"Test Flow", `{"nodes":[{"id":"1","type":"input","data":{"label":"Start"}}],"edges":[]}`, "active")
	if err != nil {
		t.Fatalf("Failed to create flow: %v", err)
	}

	// Create suggestion for retry strategy
	applyAction := map[string]interface{}{
		"action":      "implement_retry",
		"flow_id":     "1",
		"strategy":    "exponential_backoff",
		"max_retries": 3,
	}
	applyJSON, _ := json.Marshal(applyAction)

	suggestion := Suggestion{
		Type:             "retry_strategy",
		Title:            "Add retry logic",
		Description:      "Test",
		EstimatedSavings: 0.02,
		SavingsUnit:      "USD",
		TargetFlowID:     "1",
		ApplyAction:      string(applyJSON),
	}

	id, _ := StoreSuggestion(db, suggestion)

	// Apply
	result, err := ApplyOptimization(db, int(id))
	if err != nil {
		t.Fatalf("Failed to apply: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success: %s", result.Message)
	}

	// Verify retry config was added
	var flowData string
	db.QueryRow("SELECT data FROM forge_flows WHERE id = 1").Scan(&flowData)

	var flowDataMap map[string]interface{}
	json.Unmarshal([]byte(flowData), &flowDataMap)

	retryConfig, ok := flowDataMap["retryConfig"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected retryConfig to be added")
	}

	if retryConfig["enabled"] != true {
		t.Error("Expected retry to be enabled")
	}
	if retryConfig["strategy"] != "exponential_backoff" {
		t.Errorf("Expected strategy 'exponential_backoff', got '%v'", retryConfig["strategy"])
	}
}
