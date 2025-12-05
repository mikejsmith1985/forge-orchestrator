package optimizer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// FlowGraph represents the structure of a flow's data field
type FlowGraph struct {
	Nodes []FlowNode `json:"nodes"`
	Edges []FlowEdge `json:"edges"`
}

// FlowNode represents a single node in the flow
type FlowNode struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Position NodePosition `json:"position"`
	Data     NodeData     `json:"data"`
}

// NodePosition represents the x,y position of a node
type NodePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// NodeData contains the configuration data for a node
type NodeData struct {
	Label    string `json:"label"`
	Role     string `json:"role,omitempty"`
	Prompt   string `json:"prompt,omitempty"`
	Provider string `json:"provider,omitempty"`
}

// FlowEdge represents a connection between nodes
type FlowEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

// StoreSuggestion saves a suggestion to the database
func StoreSuggestion(db *sql.DB, s Suggestion) (int64, error) {
	query := `
		INSERT INTO optimization_suggestions 
		(type, title, description, estimated_savings, savings_unit, target_flow_id, target_command_id, apply_action, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'pending')
	`
	result, err := db.Exec(query, s.Type, s.Title, s.Description, s.EstimatedSavings, s.SavingsUnit, s.TargetFlowID, s.TargetCommandID, s.ApplyAction)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetSuggestionByID retrieves a suggestion from the database by ID
func GetSuggestionByID(db *sql.DB, id int) (*Suggestion, error) {
	query := `
		SELECT id, type, title, description, estimated_savings, savings_unit, 
		       target_flow_id, target_command_id, apply_action, status, applied_at, created_at
		FROM optimization_suggestions WHERE id = ?
	`
	var s Suggestion
	var targetFlowID, applyAction sql.NullString
	var targetCommandID sql.NullInt64
	var appliedAt, createdAt sql.NullTime

	err := db.QueryRow(query, id).Scan(
		&s.ID, &s.Type, &s.Title, &s.Description, &s.EstimatedSavings, &s.SavingsUnit,
		&targetFlowID, &targetCommandID, &applyAction, &s.Status, &appliedAt, &createdAt,
	)
	if err != nil {
		return nil, err
	}

	if targetFlowID.Valid {
		s.TargetFlowID = targetFlowID.String
	}
	if targetCommandID.Valid {
		s.TargetCommandID = int(targetCommandID.Int64)
	}
	if applyAction.Valid {
		s.ApplyAction = applyAction.String
	}
	if appliedAt.Valid {
		s.AppliedAt = &appliedAt.Time
	}
	if createdAt.Valid {
		s.CreatedAt = &createdAt.Time
	}

	return &s, nil
}

// GetAllSuggestions retrieves all suggestions from the database
func GetAllSuggestions(db *sql.DB) ([]Suggestion, error) {
	query := `
		SELECT id, type, title, description, estimated_savings, savings_unit, 
		       target_flow_id, target_command_id, apply_action, status, applied_at, created_at
		FROM optimization_suggestions ORDER BY created_at DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suggestions []Suggestion
	for rows.Next() {
		var s Suggestion
		var targetFlowID, applyAction sql.NullString
		var targetCommandID sql.NullInt64
		var appliedAt, createdAt sql.NullTime

		err := rows.Scan(
			&s.ID, &s.Type, &s.Title, &s.Description, &s.EstimatedSavings, &s.SavingsUnit,
			&targetFlowID, &targetCommandID, &applyAction, &s.Status, &appliedAt, &createdAt,
		)
		if err != nil {
			return nil, err
		}

		if targetFlowID.Valid {
			s.TargetFlowID = targetFlowID.String
		}
		if targetCommandID.Valid {
			s.TargetCommandID = int(targetCommandID.Int64)
		}
		if applyAction.Valid {
			s.ApplyAction = applyAction.String
		}
		if appliedAt.Valid {
			s.AppliedAt = &appliedAt.Time
		}
		if createdAt.Valid {
			s.CreatedAt = &createdAt.Time
		}

		suggestions = append(suggestions, s)
	}

	return suggestions, nil
}

// ApplyOptimization applies a stored suggestion by ID
func ApplyOptimization(db *sql.DB, suggestionID int) (*ApplyResult, error) {
	// 1. Fetch suggestion by ID
	suggestion, err := GetSuggestionByID(db, suggestionID)
	if err != nil {
		return nil, fmt.Errorf("suggestion not found: %w", err)
	}

	// Check if already applied
	if suggestion.Status == "applied" {
		return &ApplyResult{
			Success: false,
			Message: "Suggestion has already been applied",
		}, nil
	}

	// 2. Parse ApplyAction JSON
	var action ApplyAction
	if err := json.Unmarshal([]byte(suggestion.ApplyAction), &action); err != nil {
		return nil, fmt.Errorf("failed to parse apply action: %w", err)
	}

	// 3. Execute appropriate action based on type
	var result *ApplyResult
	switch action.Action {
	case "switch_model":
		result, err = applyModelSwitch(db, action)
	case "optimize_prompt":
		result, err = applyPromptOptimization(db, action)
	case "implement_retry":
		result, err = applyRetryStrategy(db, action)
	default:
		return nil, fmt.Errorf("unknown action type: %s", action.Action)
	}

	if err != nil {
		return nil, err
	}

	// 4. Mark suggestion as applied
	if result.Success {
		now := time.Now()
		_, err = db.Exec(
			"UPDATE optimization_suggestions SET status = 'applied', applied_at = ? WHERE id = ?",
			now, suggestionID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to mark suggestion as applied: %w", err)
		}
	}

	return result, nil
}

// applyModelSwitch updates a flow's node data to use a different model
func applyModelSwitch(db *sql.DB, action ApplyAction) (*ApplyResult, error) {
	if action.FlowID == "" {
		return &ApplyResult{Success: false, Message: "Flow ID is required"}, nil
	}

	// 1. Fetch flow data
	var flowData string
	err := db.QueryRow("SELECT data FROM forge_flows WHERE id = ?", action.FlowID).Scan(&flowData)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch flow: %w", err)
	}

	// 2. Parse and update nodes
	var graph FlowGraph
	if err := json.Unmarshal([]byte(flowData), &graph); err != nil {
		return nil, fmt.Errorf("failed to parse flow data: %w", err)
	}

	nodesUpdated := 0
	for i, node := range graph.Nodes {
		if node.Data.Provider == action.FromModel {
			graph.Nodes[i].Data.Provider = action.ToModel
			nodesUpdated++
		}
	}

	if nodesUpdated == 0 {
		return &ApplyResult{
			Success: false,
			Message: fmt.Sprintf("No nodes found using model %s", action.FromModel),
		}, nil
	}

	// 3. Save updated flow
	updatedData, err := json.Marshal(graph)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize flow data: %w", err)
	}

	_, err = db.Exec("UPDATE forge_flows SET data = ?, updated_at = ? WHERE id = ?",
		string(updatedData), time.Now(), action.FlowID)
	if err != nil {
		return nil, fmt.Errorf("failed to update flow: %w", err)
	}

	return &ApplyResult{
		Success:        true,
		Message:        fmt.Sprintf("Successfully switched %d node(s) from %s to %s", nodesUpdated, action.FromModel, action.ToModel),
		ChangesApplied: fmt.Sprintf("Updated %d nodes in flow %s", nodesUpdated, action.FlowID),
	}, nil
}

// applyPromptOptimization flags a flow for prompt review
func applyPromptOptimization(db *sql.DB, action ApplyAction) (*ApplyResult, error) {
	if action.FlowID == "" {
		return &ApplyResult{Success: false, Message: "Flow ID is required"}, nil
	}

	// For prompt optimization, we flag the flow by updating its status
	// In a real implementation, this could trigger a review workflow
	_, err := db.Exec(
		"UPDATE forge_flows SET status = 'needs_optimization', updated_at = ? WHERE id = ?",
		time.Now(), action.FlowID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to flag flow for optimization: %w", err)
	}

	return &ApplyResult{
		Success:        true,
		Message:        fmt.Sprintf("Flow %s has been flagged for prompt optimization review", action.FlowID),
		ChangesApplied: fmt.Sprintf("Flow %s status changed to 'needs_optimization'", action.FlowID),
	}, nil
}

// applyRetryStrategy adds retry configuration to flow metadata
func applyRetryStrategy(db *sql.DB, action ApplyAction) (*ApplyResult, error) {
	if action.FlowID == "" {
		return &ApplyResult{Success: false, Message: "Flow ID is required"}, nil
	}

	// Fetch current flow data
	var flowData string
	err := db.QueryRow("SELECT data FROM forge_flows WHERE id = ?", action.FlowID).Scan(&flowData)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch flow: %w", err)
	}

	// Parse flow data
	var flowDataMap map[string]interface{}
	if err := json.Unmarshal([]byte(flowData), &flowDataMap); err != nil {
		return nil, fmt.Errorf("failed to parse flow data: %w", err)
	}

	// Add retry configuration
	flowDataMap["retryConfig"] = map[string]interface{}{
		"enabled":     true,
		"strategy":    "exponential_backoff",
		"maxRetries":  3,
		"baseDelayMs": 1000,
	}

	// Save updated flow
	updatedData, err := json.Marshal(flowDataMap)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize flow data: %w", err)
	}

	_, err = db.Exec("UPDATE forge_flows SET data = ?, updated_at = ? WHERE id = ?",
		string(updatedData), time.Now(), action.FlowID)
	if err != nil {
		return nil, fmt.Errorf("failed to update flow: %w", err)
	}

	return &ApplyResult{
		Success:        true,
		Message:        fmt.Sprintf("Retry strategy added to flow %s", action.FlowID),
		ChangesApplied: "Added exponential backoff retry configuration with max 3 retries",
	}, nil
}
