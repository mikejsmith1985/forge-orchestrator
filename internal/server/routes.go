package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", s.healthHandler)
	mux.HandleFunc("/ws", s.websocketHandler)
	mux.HandleFunc("POST /api/ledger", s.handleCreateLedgerEntry)
	mux.HandleFunc("GET /api/ledger", s.handleGetLedger)
	mux.HandleFunc("POST /api/tokens/estimate", s.handleEstimateTokens)

	// Command Cards Routes
	mux.HandleFunc("GET /api/commands", s.handleGetCommands)
	mux.HandleFunc("POST /api/commands", s.handleCreateCommand)
	mux.HandleFunc("DELETE /api/commands/{id}", s.handleDeleteCommand)
	mux.HandleFunc("POST /api/commands/{id}/run", s.handleRunCommand)

	// Optimizer Routes
	mux.HandleFunc("GET /api/ledger/optimizations", s.handleGetOptimizations)
	mux.HandleFunc("POST /api/ledger/optimizations/{id}/apply", s.handleApplyOptimization)

	// Keyring Routes
	mux.HandleFunc("POST /api/keys", s.handleSetAPIKey)
	mux.HandleFunc("GET /api/keys/status", s.handleGetAPIKeyStatus)
	mux.HandleFunc("DELETE /api/keys/{provider}", s.handleDeleteAPIKey)

	// Flows Routes
	mux.HandleFunc("GET /api/flows", s.handleGetFlows)
	mux.HandleFunc("POST /api/flows", s.handleCreateFlow)
	mux.HandleFunc("PUT /api/flows/{id}", s.handleUpdateFlow)
	mux.HandleFunc("DELETE /api/flows/{id}", s.handleDeleteFlow)
	mux.HandleFunc("POST /api/flows/{id}/execute", s.handleExecuteFlow)

	return mux
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	s.handleWebSocket(w, r)
}
