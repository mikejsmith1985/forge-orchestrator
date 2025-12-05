// Package data provides database initialization and management for the Forge Orchestrator.
// This file implements the ledger service for logging LLM API usage.
package data

import (
	"database/sql"
	"time"
)

// LedgerService handles all operations related to the token ledger.
// It provides methods to log and retrieve API usage records.
// Think of it as a librarian that manages the "receipt book" for all AI calls.
type LedgerService struct {
	// db is the database connection used by this service.
	// It's a pointer to a sql.DB, which handles the actual database communication.
	db *sql.DB
}

// NewLedgerService creates a new LedgerService with the given database connection.
// This is a constructor that wires up the service to a specific database.
func NewLedgerService(db *sql.DB) *LedgerService {
	return &LedgerService{db: db}
}

// LogUsage inserts a complete record into the token_ledger table.
// This function is called after every LLM API call to track usage and costs.
//
// It takes a TokenLedgerEntry containing all the details about the API call:
// - Which flow triggered it
// - Which model was used
// - How many tokens were consumed
// - How much it cost
// - Whether it succeeded or failed
//
// Returns an error if the insert fails, nil on success.
func (s *LedgerService) LogUsage(entry TokenLedgerEntry) error {
	// SQL query to insert a new record into the token_ledger table.
	// We use a parameterized query (with ?) to prevent SQL injection attacks.
	// SQL injection is when bad actors try to sneak malicious commands into queries.
	query := `
		INSERT INTO token_ledger (
			timestamp,
			flow_id,
			model_used,
			agent_role,
			prompt_hash,
			input_tokens,
			output_tokens,
			total_cost_usd,
			latency_ms,
			status,
			error_message
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// If no timestamp is provided, use the current time.
	// This ensures every record has a timestamp.
	timestamp := entry.Timestamp
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	// Execute the insert query with all the values from the entry.
	// The order of values must match the order of columns in the query.
	result, err := s.db.Exec(
		query,
		timestamp,
		entry.FlowID,
		entry.ModelUsed,
		entry.AgentRole,
		entry.PromptHash,
		entry.InputTokens,
		entry.OutputTokens,
		entry.TotalCostUSD,
		entry.LatencyMs,
		entry.Status,
		entry.ErrorMessage,
	)

	if err != nil {
		return err
	}

	// Optionally, we can get the ID of the newly inserted row.
	// This is useful for confirmation but we don't need it for logging.
	_, err = result.LastInsertId()
	return err
}

// GetEntry retrieves a specific token ledger entry by its ID.
// This is useful for testing and for displaying individual records.
func (s *LedgerService) GetEntry(id int64) (*TokenLedgerEntry, error) {
	// SQL query to select a single record by ID.
	query := `
		SELECT 
			id,
			timestamp,
			flow_id,
			model_used,
			agent_role,
			prompt_hash,
			input_tokens,
			output_tokens,
			total_cost_usd,
			latency_ms,
			status,
			error_message
		FROM token_ledger
		WHERE id = ?
	`

	// Create an entry to hold the result.
	var entry TokenLedgerEntry

	// QueryRow returns at most one row. We then Scan to populate our struct.
	err := s.db.QueryRow(query, id).Scan(
		&entry.ID,
		&entry.Timestamp,
		&entry.FlowID,
		&entry.ModelUsed,
		&entry.AgentRole,
		&entry.PromptHash,
		&entry.InputTokens,
		&entry.OutputTokens,
		&entry.TotalCostUSD,
		&entry.LatencyMs,
		&entry.Status,
		&entry.ErrorMessage,
	)

	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// GetEntriesByFlowID retrieves all ledger entries for a specific flow.
// This is useful for analyzing the costs of a particular workflow.
func (s *LedgerService) GetEntriesByFlowID(flowID string) ([]TokenLedgerEntry, error) {
	query := `
		SELECT 
			id,
			timestamp,
			flow_id,
			model_used,
			agent_role,
			prompt_hash,
			input_tokens,
			output_tokens,
			total_cost_usd,
			latency_ms,
			status,
			error_message
		FROM token_ledger
		WHERE flow_id = ?
		ORDER BY timestamp DESC
	`

	rows, err := s.db.Query(query, flowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []TokenLedgerEntry
	for rows.Next() {
		var entry TokenLedgerEntry
		err := rows.Scan(
			&entry.ID,
			&entry.Timestamp,
			&entry.FlowID,
			&entry.ModelUsed,
			&entry.AgentRole,
			&entry.PromptHash,
			&entry.InputTokens,
			&entry.OutputTokens,
			&entry.TotalCostUSD,
			&entry.LatencyMs,
			&entry.Status,
			&entry.ErrorMessage,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

// GetLastInsertedEntry retrieves the most recently inserted token ledger entry.
// This is useful for testing to verify that an entry was inserted correctly.
func (s *LedgerService) GetLastInsertedEntry() (*TokenLedgerEntry, error) {
	query := `
		SELECT 
			id,
			timestamp,
			flow_id,
			model_used,
			agent_role,
			prompt_hash,
			input_tokens,
			output_tokens,
			total_cost_usd,
			latency_ms,
			status,
			error_message
		FROM token_ledger
		ORDER BY id DESC
		LIMIT 1
	`

	var entry TokenLedgerEntry
	err := s.db.QueryRow(query).Scan(
		&entry.ID,
		&entry.Timestamp,
		&entry.FlowID,
		&entry.ModelUsed,
		&entry.AgentRole,
		&entry.PromptHash,
		&entry.InputTokens,
		&entry.OutputTokens,
		&entry.TotalCostUSD,
		&entry.LatencyMs,
		&entry.Status,
		&entry.ErrorMessage,
	)

	if err != nil {
		return nil, err
	}

	return &entry, nil
}
