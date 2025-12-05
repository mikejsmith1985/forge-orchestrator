package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mikejsmith1985/forge-orchestrator/internal/tokenizer"
)

// LedgerEntry represents a row in the token_ledger table.
type LedgerEntry struct {
	ID           int     `json:"id"`
	Timestamp    string  `json:"timestamp"`
	FlowID       string  `json:"flow_id"`
	ModelUsed    string  `json:"model_used"`
	AgentRole    string  `json:"agent_role"`
	PromptHash   string  `json:"prompt_hash"`
	InputTokens  int     `json:"input_tokens"`
	OutputTokens int     `json:"output_tokens"`
	TotalCostUSD float64 `json:"total_cost_usd"`
	LatencyMS    int     `json:"latency_ms"`
	Status       string  `json:"status"`
	ErrorMessage string  `json:"error_message,omitempty"`
}

// handleCreateLedgerEntry inserts a new entry into the token_ledger table.
func (s *Server) handleCreateLedgerEntry(w http.ResponseWriter, r *http.Request) {
	var entry LedgerEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO token_ledger (
			flow_id, model_used, agent_role, prompt_hash, 
			input_tokens, output_tokens, total_cost_usd, 
			latency_ms, status, error_message
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(query,
		entry.FlowID, entry.ModelUsed, entry.AgentRole, entry.PromptHash,
		entry.InputTokens, entry.OutputTokens, entry.TotalCostUSD,
		entry.LatencyMS, entry.Status, entry.ErrorMessage,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// handleEstimateTokens estimates the number of tokens in a given text string.
// It uses tiktoken for accurate OpenAI tokenization or falls back to heuristic
// for other providers.
func (s *Server) handleEstimateTokens(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Text     string `json:"text"`
		Provider string `json:"provider"`
		Model    string `json:"model"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	estimator := tokenizer.NewEstimator()
	result := estimator.Estimate(req.Text, req.Provider, req.Model)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleGetLedger retrieves the history of agent executions.
// Educational Comment: We use a limit to prevent fetching too many rows at once,
// which could impact performance. The default is 50, but clients can override it.
func (s *Server) handleGetLedger(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	query := `
		SELECT id, timestamp, flow_id, model_used, agent_role, prompt_hash, 
		       input_tokens, output_tokens, total_cost_usd, latency_ms, status, error_message
		FROM token_ledger
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	entries := []LedgerEntry{}
	for rows.Next() {
		var e LedgerEntry
		var errMsg sql.NullString
		if err := rows.Scan(
			&e.ID, &e.Timestamp, &e.FlowID, &e.ModelUsed, &e.AgentRole, &e.PromptHash,
			&e.InputTokens, &e.OutputTokens, &e.TotalCostUSD, &e.LatencyMS, &e.Status, &errMsg,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if errMsg.Valid {
			e.ErrorMessage = errMsg.String
		}
		entries = append(entries, e)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}
