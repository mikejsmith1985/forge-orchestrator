package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mikejsmith1985/forge-orchestrator/internal/optimizer"
)

// handleGetOptimizations triggers the analyzer and returns a list of suggestions.
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
func (s *Server) handleApplyOptimization(w http.ResponseWriter, r *http.Request) {
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

	// Apply the optimization using the applier
	result, err := optimizer.ApplyOptimization(s.db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !result.Success {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
