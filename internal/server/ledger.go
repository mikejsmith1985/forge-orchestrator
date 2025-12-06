package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/mikejsmith1985/forge-orchestrator/internal/data"
	"github.com/mikejsmith1985/forge-orchestrator/internal/tokenizer"
)

// LedgerEntryResponse is the JSON API representation of a ledger entry.
// It uses string for timestamp to ensure consistent JSON serialization.
type LedgerEntryResponse struct {
	ID           int64   `json:"id"`
	Timestamp    string  `json:"timestamp"`
	FlowID       string  `json:"flow_id"`
	ModelUsed    string  `json:"model_used"`
	AgentRole    string  `json:"agent_role"`
	PromptHash   string  `json:"prompt_hash"`
	InputTokens  int     `json:"input_tokens"`
	OutputTokens int     `json:"output_tokens"`
	TotalCostUSD float64 `json:"total_cost_usd"`
	LatencyMs    int     `json:"latency_ms"`
	Status       string  `json:"status"`
	ErrorMessage string  `json:"error_message,omitempty"`
}

// ToResponse converts a TokenLedgerEntry to its API response format.
func ToLedgerResponse(entry data.TokenLedgerEntry) LedgerEntryResponse {
	return LedgerEntryResponse{
		ID:           entry.ID,
		Timestamp:    entry.Timestamp.Format(time.RFC3339),
		FlowID:       entry.FlowID,
		ModelUsed:    entry.ModelUsed,
		AgentRole:    entry.AgentRole,
		PromptHash:   entry.PromptHash,
		InputTokens:  entry.InputTokens,
		OutputTokens: entry.OutputTokens,
		TotalCostUSD: entry.TotalCostUSD,
		LatencyMs:    entry.LatencyMs,
		Status:       entry.Status,
		ErrorMessage: entry.ErrorMessage,
	}
}

// LedgerEntryRequest is the JSON API representation for creating a ledger entry.
type LedgerEntryRequest struct {
	FlowID       string  `json:"flow_id"`
	ModelUsed    string  `json:"model_used"`
	AgentRole    string  `json:"agent_role"`
	PromptHash   string  `json:"prompt_hash"`
	InputTokens  int     `json:"input_tokens"`
	OutputTokens int     `json:"output_tokens"`
	TotalCostUSD float64 `json:"total_cost_usd"`
	LatencyMs    int     `json:"latency_ms"`
	Status       string  `json:"status"`
	ErrorMessage string  `json:"error_message,omitempty"`
}

// ToEntry converts a LedgerEntryRequest to a TokenLedgerEntry.
func (r LedgerEntryRequest) ToEntry() data.TokenLedgerEntry {
	return data.TokenLedgerEntry{
		Timestamp:    time.Now(),
		FlowID:       r.FlowID,
		ModelUsed:    r.ModelUsed,
		AgentRole:    r.AgentRole,
		PromptHash:   r.PromptHash,
		InputTokens:  r.InputTokens,
		OutputTokens: r.OutputTokens,
		TotalCostUSD: r.TotalCostUSD,
		LatencyMs:    r.LatencyMs,
		Status:       r.Status,
		ErrorMessage: r.ErrorMessage,
	}
}

// handleCreateLedgerEntry inserts a new entry into the token_ledger table.
func (s *Server) handleCreateLedgerEntry(w http.ResponseWriter, r *http.Request) {
	var req LedgerEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entry := req.ToEntry()
	ledgerService := data.NewLedgerService(s.db)
	if err := ledgerService.LogUsage(entry); err != nil {
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

	entries := []LedgerEntryResponse{}
	for rows.Next() {
		var e data.TokenLedgerEntry
		var errMsg sql.NullString
		if err := rows.Scan(
			&e.ID, &e.Timestamp, &e.FlowID, &e.ModelUsed, &e.AgentRole, &e.PromptHash,
			&e.InputTokens, &e.OutputTokens, &e.TotalCostUSD, &e.LatencyMs, &e.Status, &errMsg,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if errMsg.Valid {
			e.ErrorMessage = errMsg.String
		}
		entries = append(entries, ToLedgerResponse(e))
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// BudgetResponse represents the current budget status for the UI.
// Task 4.2: This provides the Dynamic Budget Meter data.
type BudgetResponse struct {
	TotalBudget      float64 `json:"totalBudget"`
	SpentToday       float64 `json:"spentToday"`
	RemainingBudget  float64 `json:"remainingBudget"`
	RemainingPrompts int     `json:"remainingPrompts"`
	CostUnit         string  `json:"costUnit"`
	Model            string  `json:"model"`
}

// handleGetBudget returns the current budget status for the selected model.
// Task 4.2: Provides data for the Dynamic Budget Meter UI.
func (s *Server) handleGetBudget(w http.ResponseWriter, r *http.Request) {
	// Get optional model parameter (defaults to current active model)
	model := r.URL.Query().Get("model")
	if model == "" {
		model = "gpt-4o"
	}

	// Calculate spent today from ledger
	var spentToday float64
	query := `SELECT COALESCE(SUM(total_cost_usd), 0) FROM token_ledger WHERE date(timestamp) = date('now')`
	if err := s.db.QueryRow(query).Scan(&spentToday); err != nil {
		spentToday = 0
	}

	// Default budget configuration (could be made configurable)
	totalBudget := 10.00 // $10 daily budget
	remainingBudget := totalBudget - spentToday
	if remainingBudget < 0 {
		remainingBudget = 0
	}

	// Calculate remaining prompts based on average cost per prompt
	// Using $0.01 per prompt as a reasonable estimate for GPT-4o
	avgCostPerPrompt := 0.01
	if model == "gpt-3.5-turbo" {
		avgCostPerPrompt = 0.002
	}
	remainingPrompts := int(remainingBudget / avgCostPerPrompt)

	response := BudgetResponse{
		TotalBudget:      totalBudget,
		SpentToday:       spentToday,
		RemainingBudget:  remainingBudget,
		RemainingPrompts: remainingPrompts,
		CostUnit:         "TOKEN", // Default to TOKEN for most models
		Model:            model,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
