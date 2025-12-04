package server

import (
	"encoding/json"
	"net/http"

	"github.com/mikejsmith1985/forge-orchestrator/internal/security"
	"github.com/zalando/go-keyring"
)

// SetAPIKeyRequest represents the payload for setting an API key.
type SetAPIKeyRequest struct {
	Provider string `json:"provider"`
	Key      string `json:"key"`
}

// handleSetAPIKey stores an API key in the secure keyring.
// Educational Comment: We decode the JSON body, validate the input, and then
// delegate the storage to the security package. This keeps the handler logic
// clean and focused on HTTP concerns.
func (s *Server) handleSetAPIKey(w http.ResponseWriter, r *http.Request) {
	var req SetAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Provider == "" || req.Key == "" {
		http.Error(w, "Provider and Key are required", http.StatusBadRequest)
		return
	}

	if err := security.SetAPIKey(req.Provider, req.Key); err != nil {
		http.Error(w, "Failed to set API key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

// handleGetAPIKeyStatus checks which API keys are present in the keyring.
// Educational Comment: For security reasons, we never return the actual API keys.
// Instead, we return a boolean status indicating whether a key exists for each
// known provider. This allows the frontend to show a "configured" state without
// exposing secrets.
func (s *Server) handleGetAPIKeyStatus(w http.ResponseWriter, r *http.Request) {
	// List of providers to check. In a real app, this might be dynamic or from a config.
	providers := []string{"Anthropic", "OpenAI", "Google"}
	status := make(map[string]bool)

	for _, p := range providers {
		_, err := security.GetAPIKey(p)
		if err == nil {
			status[p] = true
		} else if err == keyring.ErrNotFound {
			status[p] = false
		} else {
			// If there's an error other than NotFound, we might log it but still return false
			// or handle it differently. For now, we assume false implies "not available".
			status[p] = false
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// handleDeleteAPIKey removes an API key from the keyring.
func (s *Server) handleDeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	if provider == "" {
		http.Error(w, "Provider is required", http.StatusBadRequest)
		return
	}

	if err := security.DeleteAPIKey(provider); err != nil {
		if err == keyring.ErrNotFound {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete API key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
