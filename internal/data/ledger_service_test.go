// Package data provides database initialization and management for the Forge Orchestrator.
// This test file verifies that the ledger service correctly logs and retrieves entries.
package data

import (
	"os"
	"testing"
	"time"
)

// TestLogUsageInsertsAndRetrievesEntry verifies that LogUsage inserts a complete
// TokenLedgerEntry into the database and that it can be retrieved with matching content.
// This is the main validation for Contract 4.
func TestLogUsageInsertsAndRetrievesEntry(t *testing.T) {
	// Create a temporary database for testing.
	tempDB := "test_ledger_service.db"
	defer os.Remove(tempDB)

	// Initialize the database with all tables.
	db, err := InitializeDatabase(tempDB)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create the ledger service.
	service := NewLedgerService(db)

	// Create a test entry with all fields populated.
	testTimestamp := time.Now().Truncate(time.Second)
	testEntry := TokenLedgerEntry{
		Timestamp:    testTimestamp,
		FlowID:       "test-flow-123",
		ModelUsed:    "claude-3-5-sonnet",
		AgentRole:    "Architect",
		PromptHash:   "abc123hash",
		InputTokens:  1000,
		OutputTokens: 500,
		TotalCostUSD: 0.018,
		LatencyMs:    1500,
		Status:       "SUCCESS",
		ErrorMessage: "",
	}

	// Insert the entry using LogUsage.
	err = service.LogUsage(testEntry)
	if err != nil {
		t.Fatalf("Failed to log usage: %v", err)
	}

	// Retrieve the last inserted entry.
	retrieved, err := service.GetLastInsertedEntry()
	if err != nil {
		t.Fatalf("Failed to retrieve entry: %v", err)
	}

	// Verify each field matches the inserted data.
	if retrieved.FlowID != testEntry.FlowID {
		t.Errorf("FlowID mismatch: expected %q, got %q", testEntry.FlowID, retrieved.FlowID)
	}
	if retrieved.ModelUsed != testEntry.ModelUsed {
		t.Errorf("ModelUsed mismatch: expected %q, got %q", testEntry.ModelUsed, retrieved.ModelUsed)
	}
	if retrieved.AgentRole != testEntry.AgentRole {
		t.Errorf("AgentRole mismatch: expected %q, got %q", testEntry.AgentRole, retrieved.AgentRole)
	}
	if retrieved.PromptHash != testEntry.PromptHash {
		t.Errorf("PromptHash mismatch: expected %q, got %q", testEntry.PromptHash, retrieved.PromptHash)
	}
	if retrieved.InputTokens != testEntry.InputTokens {
		t.Errorf("InputTokens mismatch: expected %d, got %d", testEntry.InputTokens, retrieved.InputTokens)
	}
	if retrieved.OutputTokens != testEntry.OutputTokens {
		t.Errorf("OutputTokens mismatch: expected %d, got %d", testEntry.OutputTokens, retrieved.OutputTokens)
	}
	if retrieved.TotalCostUSD != testEntry.TotalCostUSD {
		t.Errorf("TotalCostUSD mismatch: expected %f, got %f", testEntry.TotalCostUSD, retrieved.TotalCostUSD)
	}
	if retrieved.LatencyMs != testEntry.LatencyMs {
		t.Errorf("LatencyMs mismatch: expected %d, got %d", testEntry.LatencyMs, retrieved.LatencyMs)
	}
	if retrieved.Status != testEntry.Status {
		t.Errorf("Status mismatch: expected %q, got %q", testEntry.Status, retrieved.Status)
	}
	if retrieved.ErrorMessage != testEntry.ErrorMessage {
		t.Errorf("ErrorMessage mismatch: expected %q, got %q", testEntry.ErrorMessage, retrieved.ErrorMessage)
	}

	// Verify the ID was assigned (should be 1 for first insert).
	if retrieved.ID == 0 {
		t.Error("Expected ID to be assigned, got 0")
	}
}

// TestLogUsageWithFailedStatus verifies that failed entries are logged correctly.
func TestLogUsageWithFailedStatus(t *testing.T) {
	tempDB := "test_ledger_failed.db"
	defer os.Remove(tempDB)

	db, err := InitializeDatabase(tempDB)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	service := NewLedgerService(db)

	// Create a failed entry.
	testEntry := TokenLedgerEntry{
		FlowID:       "failed-flow-456",
		ModelUsed:    "gpt-4o",
		AgentRole:    "Coder",
		PromptHash:   "def456hash",
		InputTokens:  500,
		OutputTokens: 0, // No output because it failed
		TotalCostUSD: 0.0025,
		LatencyMs:    30000, // 30 seconds (timeout)
		Status:       "TIMEOUT",
		ErrorMessage: "Request timed out after 30 seconds",
	}

	err = service.LogUsage(testEntry)
	if err != nil {
		t.Fatalf("Failed to log failed usage: %v", err)
	}

	retrieved, err := service.GetLastInsertedEntry()
	if err != nil {
		t.Fatalf("Failed to retrieve entry: %v", err)
	}

	// Verify the error message was stored.
	if retrieved.ErrorMessage != testEntry.ErrorMessage {
		t.Errorf("ErrorMessage mismatch: expected %q, got %q", testEntry.ErrorMessage, retrieved.ErrorMessage)
	}
	if retrieved.Status != "TIMEOUT" {
		t.Errorf("Status mismatch: expected TIMEOUT, got %q", retrieved.Status)
	}
}

// TestGetEntriesByFlowID verifies that entries can be retrieved by flow ID.
func TestGetEntriesByFlowID(t *testing.T) {
	tempDB := "test_ledger_by_flow.db"
	defer os.Remove(tempDB)

	db, err := InitializeDatabase(tempDB)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	service := NewLedgerService(db)

	// Insert multiple entries for the same flow.
	flowID := "multi-entry-flow"
	for i := 0; i < 3; i++ {
		entry := TokenLedgerEntry{
			FlowID:       flowID,
			ModelUsed:    "claude-3-5-sonnet",
			AgentRole:    "Architect",
			PromptHash:   "hash" + string(rune('A'+i)),
			InputTokens:  100 * (i + 1),
			OutputTokens: 50 * (i + 1),
			TotalCostUSD: 0.01 * float64(i+1),
			LatencyMs:    1000 * (i + 1),
			Status:       "SUCCESS",
		}
		err = service.LogUsage(entry)
		if err != nil {
			t.Fatalf("Failed to log entry %d: %v", i, err)
		}
	}

	// Insert one entry for a different flow.
	differentEntry := TokenLedgerEntry{
		FlowID:       "different-flow",
		ModelUsed:    "gpt-4o",
		AgentRole:    "Coder",
		PromptHash:   "differentHash",
		InputTokens:  200,
		OutputTokens: 100,
		TotalCostUSD: 0.02,
		LatencyMs:    2000,
		Status:       "SUCCESS",
	}
	err = service.LogUsage(differentEntry)
	if err != nil {
		t.Fatalf("Failed to log different entry: %v", err)
	}

	// Retrieve entries for the target flow.
	entries, err := service.GetEntriesByFlowID(flowID)
	if err != nil {
		t.Fatalf("Failed to get entries by flow ID: %v", err)
	}

	// Should have exactly 3 entries for our test flow.
	if len(entries) != 3 {
		t.Errorf("Expected 3 entries for flow %q, got %d", flowID, len(entries))
	}

	// All entries should have the correct flow ID.
	for _, entry := range entries {
		if entry.FlowID != flowID {
			t.Errorf("Entry has wrong flow ID: expected %q, got %q", flowID, entry.FlowID)
		}
	}
}

// TestGetEntryByID verifies that a specific entry can be retrieved by its ID.
func TestGetEntryByID(t *testing.T) {
	tempDB := "test_ledger_by_id.db"
	defer os.Remove(tempDB)

	db, err := InitializeDatabase(tempDB)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	service := NewLedgerService(db)

	// Insert an entry.
	testEntry := TokenLedgerEntry{
		FlowID:       "id-test-flow",
		ModelUsed:    "claude-3-5-sonnet",
		AgentRole:    "Architect",
		PromptHash:   "idTestHash",
		InputTokens:  750,
		OutputTokens: 250,
		TotalCostUSD: 0.015,
		LatencyMs:    1200,
		Status:       "SUCCESS",
	}

	err = service.LogUsage(testEntry)
	if err != nil {
		t.Fatalf("Failed to log entry: %v", err)
	}

	// Get the entry by ID (should be 1 for first insert).
	retrieved, err := service.GetEntry(1)
	if err != nil {
		t.Fatalf("Failed to get entry by ID: %v", err)
	}

	if retrieved.FlowID != testEntry.FlowID {
		t.Errorf("FlowID mismatch: expected %q, got %q", testEntry.FlowID, retrieved.FlowID)
	}
}

// TestLogUsageAutoTimestamp verifies that a timestamp is auto-generated if not provided.
func TestLogUsageAutoTimestamp(t *testing.T) {
	tempDB := "test_ledger_auto_timestamp.db"
	defer os.Remove(tempDB)

	db, err := InitializeDatabase(tempDB)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	service := NewLedgerService(db)

	// Create an entry without a timestamp (zero value).
	beforeInsert := time.Now()
	testEntry := TokenLedgerEntry{
		// Timestamp is zero (not set)
		FlowID:       "auto-timestamp-flow",
		ModelUsed:    "claude-3-5-sonnet",
		AgentRole:    "Architect",
		PromptHash:   "autoTimestampHash",
		InputTokens:  100,
		OutputTokens: 50,
		TotalCostUSD: 0.005,
		LatencyMs:    500,
		Status:       "SUCCESS",
	}

	err = service.LogUsage(testEntry)
	if err != nil {
		t.Fatalf("Failed to log entry: %v", err)
	}
	afterInsert := time.Now()

	retrieved, err := service.GetLastInsertedEntry()
	if err != nil {
		t.Fatalf("Failed to retrieve entry: %v", err)
	}

	// Verify timestamp is within the expected range.
	if retrieved.Timestamp.Before(beforeInsert) || retrieved.Timestamp.After(afterInsert) {
		t.Errorf("Timestamp %v not in expected range [%v, %v]", retrieved.Timestamp, beforeInsert, afterInsert)
	}
}
