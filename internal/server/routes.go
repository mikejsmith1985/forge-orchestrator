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
