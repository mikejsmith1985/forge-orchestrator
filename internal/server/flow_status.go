package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mikejsmith1985/forge-orchestrator/internal/flows"
)

// handleGetFlowStatus returns the current status of a flow
func (s *Server) handleGetFlowStatus(w http.ResponseWriter, r *http.Request) {
	// Extract flow ID from URL path
	idStr := r.PathValue("id")
	flowID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid flow ID", http.StatusBadRequest)
		return
	}

	// Try to get status from file signaler (fallback storage)
	fileSignaler, err := flows.NewFileSignaler()
	if err != nil {
		http.Error(w, "Failed to initialize status reader", http.StatusInternalServerError)
		return
	}

	status, err := fileSignaler.GetStatus(flowID)
	if err != nil {
		// Return a default pending status if not found
		status = &flows.FlowStatus{
			FlowID: flowID,
			Status: "UNKNOWN",
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
