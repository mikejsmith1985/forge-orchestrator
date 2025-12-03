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

	// Optimizer Routes
	mux.HandleFunc("GET /api/ledger/optimizations", s.handleGetOptimizations)
	mux.HandleFunc("POST /api/ledger/optimizations/{id}/apply", s.handleApplyOptimization)

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
