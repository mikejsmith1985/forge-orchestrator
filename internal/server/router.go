// Package server provides HTTP routing for the Forge Orchestrator.
// This file defines the router setup as specified in Contract 5.
// It works alongside routes.go to organize API endpoints.
package server

import (
	"net/http"
)

// Router wraps an http.ServeMux and provides methods for registering routes.
// This is a simple abstraction that makes route registration more organized.
type Router struct {
	mux *http.ServeMux
}

// NewRouter creates a new Router instance.
// The router uses Go's standard library ServeMux under the hood.
func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

// Handle registers a handler for a specific pattern.
// This is a convenience method that delegates to the underlying ServeMux.
func (r *Router) Handle(pattern string, handler http.Handler) {
	r.mux.Handle(pattern, handler)
}

// HandleFunc registers a handler function for a specific pattern.
// This is used for registering API endpoints.
func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc) {
	r.mux.HandleFunc(pattern, handler)
}

// ServeHTTP implements the http.Handler interface.
// This allows the Router to be used directly with http.ListenAndServe.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// GetMux returns the underlying http.ServeMux.
// This is useful when you need direct access to the mux (e.g., for static file serving).
func (r *Router) GetMux() *http.ServeMux {
	return r.mux
}

// SetupAPIRoutes registers all API routes on the given Server.
// This function is called during server initialization to wire up all endpoints.
// It includes the /api/execute endpoint required by Contract 5.
func SetupAPIRoutes(s *Server, r *Router) {
	// Health check endpoint - always first for monitoring/load balancers.
	r.HandleFunc("/api/health", s.healthHandler)

	// Execute endpoint - Contract 5 requirement.
	// This calls the Executor interface to run shell commands.
	r.HandleFunc("POST /api/execute", s.handleExecute)

	// WebSocket endpoint for real-time terminal communication.
	r.HandleFunc("/ws", s.websocketHandler)

	// Ledger endpoints - for tracking API usage.
	r.HandleFunc("POST /api/ledger", s.handleCreateLedgerEntry)
	r.HandleFunc("GET /api/ledger", s.handleGetLedger)
	r.HandleFunc("POST /api/tokens/estimate", s.handleEstimateTokens)

	// Command Cards endpoints - for saved commands.
	r.HandleFunc("GET /api/commands", s.handleGetCommands)
	r.HandleFunc("POST /api/commands", s.handleCreateCommand)
	r.HandleFunc("DELETE /api/commands/{id}", s.handleDeleteCommand)
	r.HandleFunc("POST /api/commands/{id}/run", s.handleRunCommand)

	// Optimizer endpoints - for cost optimization suggestions.
	r.HandleFunc("GET /api/ledger/optimizations", s.handleGetOptimizations)
	r.HandleFunc("POST /api/ledger/optimizations/{id}/apply", s.handleApplyOptimization)

	// Keyring endpoints - for API key management.
	r.HandleFunc("POST /api/keys", s.handleSetAPIKey)
	r.HandleFunc("GET /api/keys/status", s.handleGetAPIKeyStatus)
	r.HandleFunc("DELETE /api/keys/{provider}", s.handleDeleteAPIKey)

	// Flows endpoints - for workflow management.
	r.HandleFunc("GET /api/flows", s.handleGetFlows)
	r.HandleFunc("GET /api/flows/{id}", s.handleGetFlow)
	r.HandleFunc("POST /api/flows", s.handleCreateFlow)
	r.HandleFunc("PUT /api/flows/{id}", s.handleUpdateFlow)
	r.HandleFunc("DELETE /api/flows/{id}", s.handleDeleteFlow)
	r.HandleFunc("POST /api/flows/{id}/execute", s.handleExecuteFlow)
	r.HandleFunc("GET /api/flows/{id}/status", s.handleGetFlowStatus)
}
