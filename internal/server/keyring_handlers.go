package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mikejsmith1985/forge-orchestrator/internal/security"
	"github.com/zalando/go-keyring"
)

// SetAPIKeyRequest represents the payload for setting an API key.
type SetAPIKeyRequest struct {
	Provider string `json:"provider"`
	Key      string `json:"key"`
}

// SetAPIKeyResponse represents the response after setting an API key.
type SetAPIKeyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// KeyStatus represents the status of a single API key provider.
type KeyStatus struct {
	Provider string `json:"provider"`
	IsSet    bool   `json:"isSet"`
}

// KeyStatusResponse represents the response for the key status endpoint.
type KeyStatusResponse struct {
	Keys []KeyStatus `json:"keys"`
}

// handleSetAPIKey stores an API key in the secure keyring.
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SetAPIKeyResponse{
		Status:  "ok",
		Message: "API key saved successfully",
	})
}

// handleGetAPIKeyStatus checks which API keys are present in the keyring.
// Returns a structured response with lowercase provider names for frontend consistency.
func (s *Server) handleGetAPIKeyStatus(w http.ResponseWriter, r *http.Request) {
	// List of providers to check
	providers := []string{"Anthropic", "OpenAI", "Google"}
	keys := make([]KeyStatus, 0, len(providers))

	for _, p := range providers {
		_, err := security.GetAPIKey(p)
		isSet := err == nil

		keys = append(keys, KeyStatus{
			Provider: strings.ToLower(p),
			IsSet:    isSet,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(KeyStatusResponse{Keys: keys})
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
