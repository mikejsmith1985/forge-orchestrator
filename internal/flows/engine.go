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

// Broadcaster defines the interface for broadcasting messages to WebSocket clients
type Broadcaster interface {
	Broadcast(message []byte)
}

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

// notifyStatus sends status update via WebSocket with file fallback
func notifyStatus(wsSignaler, fileSignaler Signaler, flowID int, status FlowStatus) {
	// Always write to file for fallback polling
	if fileSignaler != nil {
		if err := fileSignaler.NotifyStatus(flowID, status); err != nil {
			log.Printf("File signaler failed: %v", err)
		}
	}

	// Try WebSocket notification
	if wsSignaler != nil {
		if err := wsSignaler.NotifyStatus(flowID, status); err != nil {
			log.Printf("WebSocket notify failed, using file fallback: %v", err)
		}
	}
}

// ExecuteFlowWithSignaling runs the flow with status signaling support
func ExecuteFlowWithSignaling(flowID int, db *sql.DB, gateway *llm.Gateway, wsSignaler, fileSignaler Signaler) error {
	return ExecuteFlowWithHub(flowID, db, gateway, wsSignaler, fileSignaler, nil)
}

// ExecuteFlowWithHub runs the flow with full Hub integration for real-time broadcasts
func ExecuteFlowWithHub(flowID int, db *sql.DB, gateway *llm.Gateway, wsSignaler, fileSignaler Signaler, hub Broadcaster) error {
	startTime := time.Now()

	// Broadcast FLOW_STARTED
	if hub != nil {
		hub.Broadcast(NewFlowStartedMessage(flowID))
	}

	// Notify flow started (legacy signaler)
	notifyStatus(wsSignaler, fileSignaler, flowID, FlowStatus{
		FlowID:    flowID,
		Status:    "RUNNING",
		UpdatedAt: time.Now(),
	})

	err := executeFlowInternalWithHub(flowID, db, gateway, wsSignaler, fileSignaler, hub)

	executionTime := time.Since(startTime).Milliseconds()

	// Notify flow completed or failed
	if err != nil {
		// Broadcast FLOW_FAILED
		if hub != nil {
			hub.Broadcast(NewFlowFailedMessage(flowID, err.Error()))
		}
		notifyStatus(wsSignaler, fileSignaler, flowID, FlowStatus{
			FlowID:    flowID,
			Status:    "FAILED",
			UpdatedAt: time.Now(),
			Error:     err.Error(),
		})
		return err
	}

	// Broadcast FLOW_COMPLETED
	if hub != nil {
		hub.Broadcast(NewFlowCompletedMessage(flowID, executionTime))
	}

	notifyStatus(wsSignaler, fileSignaler, flowID, FlowStatus{
		FlowID:    flowID,
		Status:    "COMPLETED",
		UpdatedAt: time.Now(),
	})

	return nil
}

// executeFlowInternalWithHub contains the core flow execution logic with Hub broadcasting
func executeFlowInternalWithHub(flowID int, db *sql.DB, gateway *llm.Gateway, wsSignaler, fileSignaler Signaler, hub Broadcaster) error {
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
	for _, node := range graph.Nodes {
		if node.Type != "agent" {
			continue // Skip non-agent nodes if any
		}

		// Broadcast NODE_STARTED
		if hub != nil {
			hub.Broadcast(NewNodeStartedMessage(flowID, node.ID, node.Data.Label))
		}

		// Notify node starting (legacy signaler)
		notifyStatus(wsSignaler, fileSignaler, flowID, FlowStatus{
			FlowID:    flowID,
			Status:    "RUNNING",
			LastNode:  node.ID,
			UpdatedAt: time.Now(),
		})

		// Get API Key
		apiKey, err := security.GetAPIKey(node.Data.Provider)
		if err != nil {
			log.Printf("Error getting API key for provider %s: %v", node.Data.Provider, err)
			return fmt.Errorf("missing API key for provider %s: %w", node.Data.Provider, err)
		}

		// Execute Prompt
		providerType := llm.ProviderType(node.Data.Provider)

		start := time.Now()
		resp, err := gateway.ExecutePrompt(node.Data.Role, node.Data.Prompt, apiKey, providerType)
		latency := time.Since(start).Milliseconds()

		status := "SUCCESS"
		var errMsg string
		var inputTokens, outputTokens int
		var cost float64
		var promptHash string = "hash_placeholder"

		if err != nil {
			status = "FAILED"
			errMsg = err.Error()
			log.Printf("Node %s execution failed: %v", node.ID, err)
		} else {
			inputTokens = resp.InputTokens
			outputTokens = resp.OutputTokens
			cost = resp.Cost
		}

		// Broadcast NODE_COMPLETED (even if failed, we report the tokens used)
		if hub != nil {
			hub.Broadcast(NewNodeCompletedMessage(flowID, node.ID, inputTokens, outputTokens, cost))
		}

		// 4. Log to token_ledger
		insertQuery := `
			INSERT INTO token_ledger (
				flow_id, model_used, agent_role, prompt_hash, 
				input_tokens, output_tokens, total_cost_usd, 
				latency_ms, status, error_message
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		_, dbErr := db.Exec(insertQuery,
			fmt.Sprintf("%d", flowID),
			node.Data.Provider,
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
		}

		if err != nil {
			return fmt.Errorf("node %s failed: %w", node.ID, err)
		}
	}

	return nil
}

// ExecuteFlow runs the flow with the given ID (backwards compatible version without signaling)
func ExecuteFlow(flowID int, db *sql.DB, gateway *llm.Gateway) error {
	// Create file signaler for basic status tracking
	fileSignaler, _ := NewFileSignaler()
	return ExecuteFlowWithSignaling(flowID, db, gateway, nil, fileSignaler)
}
