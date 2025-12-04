package flows

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mikejsmith1985/forge-orchestrator/internal/llm"
	"github.com/mikejsmith1985/forge-orchestrator/internal/security"
)

// Flow represents a flow definition from the database.
type Flow struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Data      string    `json:"data"` // JSON string of the graph
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// FlowGraph represents the parsed JSON structure of the flow.
type FlowGraph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// Node represents a single step in the flow.
type Node struct {
	ID   string   `json:"id"`
	Type string   `json:"type"` // e.g., "agent"
	Data NodeData `json:"data"`
}

// NodeData contains the configuration for a node.
type NodeData struct {
	Label    string `json:"label"`
	Role     string `json:"role"`     // e.g., "coder", "planner"
	Prompt   string `json:"prompt"`   // The user input/task for this agent
	Provider string `json:"provider"` // e.g., "Anthropic", "OpenAI"
}

// Edge represents a connection between nodes.
type Edge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

// ExecuteFlow runs the flow with the given ID.
// Educational Comment: This function orchestrates the execution of a flow.
// It fetches the flow definition, parses the graph, and executes nodes sequentially.
// In a real-world scenario, this would handle topological sorting for parallel execution
// and complex dependency management.
func ExecuteFlow(flowID int, db *sql.DB, gateway *llm.Gateway) error {
	// 1. Fetch flow data
	var flowData string
	query := `SELECT data FROM forge_flows WHERE id = ?`
	err := db.QueryRow(query, flowID).Scan(&flowData)
	if err != nil {
		return fmt.Errorf("failed to fetch flow: %w", err)
	}

	// 2. Parse JSON graph
	var graph FlowGraph
	if err := json.Unmarshal([]byte(flowData), &graph); err != nil {
		return fmt.Errorf("failed to parse flow data: %w", err)
	}

	// 3. Execute nodes (Sequential for now)
	// Educational Comment: We iterate through the nodes in the order they appear in the JSON.
	// A more advanced implementation would use the Edges to determine execution order (DAG).
	for _, node := range graph.Nodes {
		if node.Type != "agent" {
			continue // Skip non-agent nodes if any
		}

		// Get API Key
		apiKey, err := security.GetAPIKey(node.Data.Provider)
		if err != nil {
			log.Printf("Error getting API key for provider %s: %v", node.Data.Provider, err)
			// Log failure to ledger? Or just continue/fail?
			// Let's fail the flow execution for now.
			return fmt.Errorf("missing API key for provider %s: %w", node.Data.Provider, err)
		}

		// Execute Prompt
		// We map the string provider from JSON to the ProviderType enum
		providerType := llm.ProviderType(node.Data.Provider)

		start := time.Now()
		resp, err := gateway.ExecutePrompt(node.Data.Role, node.Data.Prompt, apiKey, providerType)
		latency := time.Since(start).Milliseconds()

		status := "SUCCESS"
		var errMsg string
		var inputTokens, outputTokens int
		var cost float64
		var promptHash string = "hash_placeholder" // In a real app, hash the prompt

		if err != nil {
			status = "FAILED"
			errMsg = err.Error()
			log.Printf("Node %s execution failed: %v", node.ID, err)
		} else {
			inputTokens = resp.InputTokens
			outputTokens = resp.OutputTokens
			cost = resp.Cost
		}

		// 4. Log to token_ledger
		// Educational Comment: We log every execution step to the ledger for auditing and cost tracking.
		// This is critical for understanding system behavior and managing expenses.
		insertQuery := `
			INSERT INTO token_ledger (
				flow_id, model_used, agent_role, prompt_hash, 
				input_tokens, output_tokens, total_cost_usd, 
				latency_ms, status, error_message
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		_, dbErr := db.Exec(insertQuery,
			fmt.Sprintf("%d", flowID), // Use Flow ID as the tracking ID
			node.Data.Provider,        // Using Provider as Model for now, or could be specific model name if available
			node.Data.Role,
			promptHash,
			inputTokens,
			outputTokens,
			cost,
			latency,
			status,
			errMsg,
		)
		if dbErr != nil {
			log.Printf("Failed to log to ledger: %v", dbErr)
			// Don't fail the flow just because logging failed, but it's bad practice.
		}

		if err != nil {
			return fmt.Errorf("node %s failed: %w", node.ID, err)
		}
	}

	return nil
}
