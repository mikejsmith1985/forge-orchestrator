package flows

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3" // Use mattn/go-sqlite3
	"github.com/mikejsmith1985/forge-orchestrator/internal/data"
	"github.com/mikejsmith1985/forge-orchestrator/internal/llm"
	"github.com/mikejsmith1985/forge-orchestrator/internal/security"
	"github.com/zalando/go-keyring"
)

// MockLLMProvider implements llm.LLMProvider for testing.
type MockLLMProvider struct {
	Called      bool
	LastPrompt  string
	ReturnValue string
	Err         error
}

func (m *MockLLMProvider) Send(system, user, key string) (string, int, int, error) {
	m.Called = true
	m.LastPrompt = user
	return m.ReturnValue, 10, 20, m.Err
}

func TestExecuteFlow(t *testing.T) {
	// 1. Setup Mock Keyring
	keyring.MockInit()
	err := security.SetAPIKey("Anthropic", "dummy-key")
	if err != nil {
		t.Fatalf("Failed to set mock API key: %v", err)
	}

	// 2. Setup In-Memory DB
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory db: %v", err)
	}
	defer db.Close()

	// Initialize Schema
	_, err = db.Exec(data.SQLiteSchema)
	if err != nil {
		t.Fatalf("Failed to init schema: %v", err)
	}

	// 3. Insert Test Flow
	flowJSON := `{
		"nodes": [
			{
				"id": "1",
				"type": "agent",
				"data": {
					"label": "Coder",
					"role": "Implementation",
					"prompt": "Write a hello world function",
					"provider": "Anthropic"
				}
			}
		],
		"edges": []
	}`
	_, err = db.Exec(`INSERT INTO forge_flows (name, data, status) VALUES (?, ?, ?)`, "Test Flow", flowJSON, "active")
	if err != nil {
		t.Fatalf("Failed to insert flow: %v", err)
	}

	// 4. Setup Gateway with Mock Provider
	mockProvider := &MockLLMProvider{ReturnValue: "Here is the code"}
	gateway := &llm.Gateway{
		AnthropicClient: mockProvider,
		OpenAIClient:    &MockLLMProvider{}, // Unused in this test
	}

	// 5. Execute Flow
	// We need the ID of the inserted flow, which should be 1
	err = ExecuteFlow(1, db, gateway)
	if err != nil {
		t.Fatalf("ExecuteFlow failed: %v", err)
	}

	// 6. Verify Results
	// Check if provider was called
	if !mockProvider.Called {
		t.Error("Expected provider to be called, but it wasn't")
	}
	if mockProvider.LastPrompt != "Write a hello world function" {
		t.Errorf("Expected prompt 'Write a hello world function', got '%s'", mockProvider.LastPrompt)
	}

	// Check Ledger
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM token_ledger").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query ledger: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 ledger entry, got %d", count)
	}
}
