package server

import (
	"encoding/json"
	"net/http"
)

// LedgerEntry represents a row in the token_ledger table.
type LedgerEntry struct {
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
// It uses a simple approximation: len(text) / 4.
func (s *Server) handleEstimateTokens(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Simple approximation: len(text) / 4
	count := len(req.Text) / 4

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"count": count})
}
