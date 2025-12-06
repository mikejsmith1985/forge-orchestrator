// Package server provides HTTP handlers for the Forge Orchestrator API.
// This file contains API handlers for executing commands via the Executor interface.
package server

import (
	"encoding/json"
	"net/http"

	"github.com/mikejsmith1985/forge-orchestrator/internal/execution"
)

// ExecuteRequest represents the JSON payload for the /api/execute endpoint.
// This is what the frontend sends when it wants to run a command.
type ExecuteRequest struct {
	// Command is the shell command to execute.
	Command string `json:"command"`

	// WorkingDir is optional - if not provided, uses the current directory.
	WorkingDir string `json:"workingDir,omitempty"`

	// TimeoutSeconds is optional - if not provided, uses default (no timeout).
	TimeoutSeconds int `json:"timeoutSeconds,omitempty"`
}

// ExecuteResponse represents the JSON response from the /api/execute endpoint.
// This tells the frontend what happened when we ran the command.
type ExecuteResponse struct {
	// Message is a human-readable status message.
	Message string `json:"message"`

	// Stdout is the standard output from the command (if executed).
	Stdout string `json:"stdout,omitempty"`

	// Stderr is the standard error output from the command (if executed).
	Stderr string `json:"stderr,omitempty"`

	// ExitCode is the numeric exit code from the command.
	ExitCode int `json:"exitCode"`

	// Success indicates whether the execution was successful.
	Success bool `json:"success"`
}

// PTYCommandRequest represents the JSON payload for injecting commands into PTY.
// Task 2.2: This is used by Flow Nodes and Command Cards to execute in the terminal.
type PTYCommandRequest struct {
	// SessionID identifies which PTY session to inject the command into.
	SessionID string `json:"sessionId"`

	// Command is the command string to inject (simulates typing).
	Command string `json:"command"`
}

// PTYCommandResponse represents the response from PTY command injection.
type PTYCommandResponse struct {
	// Success indicates whether the command was injected successfully.
	Success bool `json:"success"`

	// Message is a human-readable status message.
	Message string `json:"message"`
}

// executor is the Executor interface instance used for running commands.
// We use the interface type so we can swap implementations (e.g., for testing).
var executor execution.Executor

// init initializes the executor with a LocalRunner by default.
// This runs when the package is loaded.
func init() {
	executor = execution.NewLocalRunner()
}

// SetExecutor allows replacing the executor (useful for testing).
func SetExecutor(e execution.Executor) {
	executor = e
}

// handleExecute processes requests to the /api/execute endpoint.
// It receives a command from the frontend, runs it using the Executor interface,
// and returns the result.
//
// Flow:
// 1. Parse the JSON request body
// 2. Create an ExecutionContext from the request
// 3. Call the Executor interface to run the command
// 4. Return the result as JSON
func (s *Server) handleExecute(w http.ResponseWriter, r *http.Request) {
	// Set response content type to JSON.
	w.Header().Set("Content-Type", "application/json")

	// Only accept POST requests.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ExecuteResponse{
			Message: "Method not allowed. Use POST.",
			Success: false,
		})
		return
	}

	// Parse the request body.
	var req ExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ExecuteResponse{
			Message: "Invalid request body: " + err.Error(),
			Success: false,
		})
		return
	}

	// Validate that a command was provided.
	if req.Command == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ExecuteResponse{
			Message: "Command is required",
			Success: false,
		})
		return
	}

	// Create an ExecutionContext from the request.
	ctx := execution.ExecutionContext{
		Command:        req.Command,
		WorkingDir:     req.WorkingDir,
		TimeoutSeconds: req.TimeoutSeconds,
	}

	// Execute the command using the Executor interface.
	result := executor.Execute(ctx)

	// Build the response.
	response := ExecuteResponse{
		Message:  "Execution Request Received",
		Stdout:   result.Stdout,
		Stderr:   result.Stderr,
		ExitCode: result.ExitCode,
		Success:  result.ExitCode == 0 && result.Error == nil,
	}

	// If there was an error running the command, include it in the message.
	if result.Error != nil {
		response.Message = "Execution failed: " + result.Error.Error()
		response.Success = false
	}

	// Return the response.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// handlePTYCommandExecute injects a command into an active PTY session.
// Task 2.2: This is the API used by Flow Nodes and Command Cards to execute
// commands in the integrated terminal, simulating a human typing.
func (s *Server) handlePTYCommandExecute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req PTYCommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PTYCommandResponse{
			Success: false,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	if req.Command == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PTYCommandResponse{
			Success: false,
			Message: "Command is required",
		})
		return
	}

	// Get the PTY session
	session := s.ptyManager.GetSession(req.SessionID)
	if session == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(PTYCommandResponse{
			Success: false,
			Message: "PTY session not found",
		})
		return
	}

	// Write the command to the PTY (simulates typing + Enter)
	if err := session.WriteCommand(req.Command); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PTYCommandResponse{
			Success: false,
			Message: "Failed to inject command: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PTYCommandResponse{
		Success: true,
		Message: "Command injected successfully",
	})
}
