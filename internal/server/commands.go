package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mikejsmith1985/forge-orchestrator/internal/llm"
	"github.com/mikejsmith1985/forge-orchestrator/internal/security"
)

// CommandCard represents a reusable terminal command.
type CommandCard struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	Description string `json:"description"`
}

// handleGetCommands returns all command cards.
// Educational Comment: We use a simple SELECT query to retrieve all rows.
// In a production app with many users, we'd likely need pagination or filtering here.
func (s *Server) handleGetCommands(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query("SELECT id, name, command, description FROM command_cards ORDER BY id DESC")
	if err != nil {
		http.Error(w, "Failed to query commands: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var commands []CommandCard
	for rows.Next() {
		var c CommandCard
		if err := rows.Scan(&c.ID, &c.Name, &c.Command, &c.Description); err != nil {
			http.Error(w, "Failed to scan command: "+err.Error(), http.StatusInternalServerError)
			return
		}
		commands = append(commands, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commands)
}

// handleCreateCommand adds a new command card.
// Educational Comment: We decode the JSON body into a struct, then execute an INSERT statement.
// We return the ID of the newly created row so the frontend can update its state immediately.
func (s *Server) handleCreateCommand(w http.ResponseWriter, r *http.Request) {
	var c CommandCard
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if c.Name == "" || c.Command == "" {
		http.Error(w, "Name and Command are required", http.StatusBadRequest)
		return
	}

	res, err := s.db.Exec("INSERT INTO command_cards (name, command, description) VALUES (?, ?, ?)", c.Name, c.Command, c.Description)
	if err != nil {
		http.Error(w, "Failed to insert command: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	c.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

// handleDeleteCommand removes a command card by ID.
// Educational Comment: We parse the ID from the URL path (e.g., /api/commands/123).
// Then we execute a DELETE statement.
func (s *Server) handleDeleteCommand(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = s.db.Exec("DELETE FROM command_cards WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete command: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RunCommandRequest represents the payload for running a command via LLM.
type RunCommandRequest struct {
	AgentRole  string `json:"agent_role"`
	UserPrompt string `json:"user_prompt"`
	Provider   string `json:"provider"`
}

// handleRunCommand executes a prompt using the LLM Gateway.
// Educational Comment: This handler acts as the bridge between the frontend and the LLM Gateway.
// It extracts the API key from the header for security and delegates the complex routing logic to the Gateway.
func (s *Server) handleRunCommand(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-Forge-Api-Key")

	// If API key is not in header, try to get it from the keyring
	// Educational Comment: We check the header first to allow per-request overrides.
	// If not found, we fallback to the secure keyring.
	// Note: We need the provider to look up the key, so we'll do this check
	// after decoding the body if it's still missing.
	// However, to keep the flow simple, we can decode first.
	// But wait, the original code checked the header *before* decoding.
	// Let's decode first now so we have the provider.

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req RunCommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.AgentRole == "" || req.Provider == "" {
		http.Error(w, "agent_role and provider are required", http.StatusBadRequest)
		return
	}

	if apiKey == "" {
		// Try to get from keyring
		key, err := security.GetAPIKey(req.Provider)
		if err == nil && key != "" {
			apiKey = key
		}
	}

	if apiKey == "" {
		http.Error(w, "Missing X-Forge-Api-Key header and no key found in keyring", http.StatusUnauthorized)
		return
	}

	// Fetch command from database
	var commandPrompt string
	err = s.db.QueryRow("SELECT command FROM command_cards WHERE id = ?", id).Scan(&commandPrompt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Command not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Convert string provider to ProviderType
	provider := llm.ProviderType(req.Provider)

	// Execute via Gateway
	response, err := s.gateway.ExecutePrompt(req.AgentRole, commandPrompt, apiKey, provider)

	// Prepare ledger entry
	ledgerEntry := LedgerEntry{
		FlowID:     "cmd-" + strconv.Itoa(id), // Simple flow ID for now
		ModelUsed:  string(provider),
		AgentRole:  req.AgentRole,
		PromptHash: "hash-" + strconv.Itoa(len(commandPrompt)), // Placeholder hash
		Status:     "SUCCESS",
	}

	if err != nil {
		ledgerEntry.Status = "FAILED"
		ledgerEntry.ErrorMessage = err.Error()
		// Log failure to ledger
		s.logToLedger(ledgerEntry)
		http.Error(w, "LLM execution failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update ledger entry with success details
	ledgerEntry.InputTokens = response.InputTokens
	ledgerEntry.OutputTokens = response.OutputTokens
	ledgerEntry.TotalCostUSD = response.Cost
	// Note: Latency is not captured here yet, could add timing around ExecutePrompt

	// Log success to ledger
	s.logToLedger(ledgerEntry)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// logToLedger helper to insert into token_ledger
func (s *Server) logToLedger(entry LedgerEntry) {
	query := `
		INSERT INTO token_ledger (
			flow_id, model_used, agent_role, prompt_hash, 
			input_tokens, output_tokens, total_cost_usd, 
			latency_ms, status, error_message
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	s.db.Exec(query,
		entry.FlowID, entry.ModelUsed, entry.AgentRole, entry.PromptHash,
		entry.InputTokens, entry.OutputTokens, entry.TotalCostUSD,
		entry.LatencyMS, entry.Status, entry.ErrorMessage,
	)
}
