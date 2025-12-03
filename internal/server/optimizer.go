package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mikejsmith1985/forge-orchestrator/internal/optimizer"
)

// handleGetOptimizations triggers the analyzer and returns a list of suggestions.
// Educational Comment: This endpoint is the entry point for the Token Optimizer.
// It queries the ledger, runs the analysis heuristics, and returns actionable suggestions.
func (s *Server) handleGetOptimizations(w http.ResponseWriter, r *http.Request) {
	suggestions, err := optimizer.AnalyzeLedger(s.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}

// handleApplyOptimization applies a selected optimization suggestion.
// Educational Comment: Currently, this is a placeholder that logs the action to the ledger.
// In a full implementation, this would trigger the actual configuration change (e.g., updating a flow definition).
func (s *Server) handleApplyOptimization(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path (e.g., /api/ledger/optimizations/{id}/apply)
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Optimization ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid optimization ID", http.StatusBadRequest)
		return
	}

	// In a real implementation, we would fetch the suggestion by ID and apply it.
	// Here, we just log that it was applied.
	// We'll insert a special entry into the ledger to track this action.

	query := `
		INSERT INTO token_ledger (
			flow_id, model_used, agent_role, prompt_hash, 
			input_tokens, output_tokens, total_cost_usd, 
			latency_ms, status, error_message
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Log the application of the optimization
	_, err = s.db.Exec(query,
		"system", "optimizer", "optimizer_agent", "applied_suggestion_"+strconv.Itoa(id),
		0, 0, 0.0,
		0, "OPTIMIZED", "Applied suggestion ID: "+strconv.Itoa(id),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "applied", "message": "Optimization applied successfully"})
}
