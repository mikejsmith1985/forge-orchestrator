package optimizer

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// modelCosts maps model names to their approximate cost per 1M tokens (input + output averaged).
// Educational Comment: We use this lookup table to identify expensive models and suggest cheaper alternatives.
// These are approximate costs and should be updated based on current pricing.
var modelCosts = map[string]float64{
	"gpt-4":           0.06,  // High-cost model
	"gpt-4-turbo":     0.03,  // Medium-high cost
	"gpt-3.5-turbo":   0.002, // Low-cost model
	"claude-3-opus":   0.075, // High-cost model
	"claude-3-sonnet": 0.015, // Medium cost
	"claude-3-haiku":  0.001, // Low-cost model
	"claude-2":        0.024, // Medium-high cost
}

// AnalyzeLedger queries the token_ledger table and identifies optimization opportunities.
// It stores new suggestions in the database and returns both new and existing suggestions.
func AnalyzeLedger(db *sql.DB) ([]Suggestion, error) {
	// First, get existing suggestions from database
	existingSuggestions, err := GetAllSuggestions(db)
	if err != nil {
		// If table doesn't exist yet, just continue with analysis
		existingSuggestions = []Suggestion{}
	}

	// If we have pending suggestions, return those instead of generating new ones
	hasPending := false
	for _, s := range existingSuggestions {
		if s.Status == "pending" {
			hasPending = true
			break
		}
	}
	if hasPending {
		return existingSuggestions, nil
	}

	// Generate new suggestions
	newSuggestions := []Suggestion{}

	// Analysis 1: Detect repeated high-cost model usage
	highCostSuggestions, err := detectHighCostModels(db)
	if err != nil {
		return nil, fmt.Errorf("failed to detect high-cost models: %w", err)
	}
	newSuggestions = append(newSuggestions, highCostSuggestions...)

	// Analysis 2: Identify long prompts
	longPromptSuggestions, err := detectLongPrompts(db)
	if err != nil {
		return nil, fmt.Errorf("failed to detect long prompts: %w", err)
	}
	newSuggestions = append(newSuggestions, longPromptSuggestions...)

	// Analysis 3: Find failed executions
	failureSuggestions, err := detectFailures(db)
	if err != nil {
		return nil, fmt.Errorf("failed to detect failures: %w", err)
	}
	newSuggestions = append(newSuggestions, failureSuggestions...)

	// Store new suggestions in database
	for i := range newSuggestions {
		id, err := StoreSuggestion(db, newSuggestions[i])
		if err != nil {
			return nil, fmt.Errorf("failed to store suggestion: %w", err)
		}
		newSuggestions[i].ID = int(id)
		newSuggestions[i].Status = "pending"
	}

	// Return all suggestions (new + previously applied)
	return append(newSuggestions, existingSuggestions...), nil
}

// detectHighCostModels identifies flows using expensive models and suggests cheaper alternatives.
func detectHighCostModels(db *sql.DB) ([]Suggestion, error) {
	suggestions := []Suggestion{}

	query := `
		SELECT flow_id, model_used, COUNT(*) as call_count, SUM(total_cost_usd) as total_cost
		FROM token_ledger
		WHERE model_used IN ('gpt-4', 'claude-3-opus', 'gpt-4-turbo', 'claude-2')
		GROUP BY flow_id, model_used
		HAVING call_count >= 2
		ORDER BY total_cost DESC
		LIMIT 10
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var flowID, modelUsed string
		var callCount int
		var totalCost float64

		if err := rows.Scan(&flowID, &modelUsed, &callCount, &totalCost); err != nil {
			return nil, err
		}

		alternative := "gpt-3.5-turbo"
		if modelUsed == "claude-3-opus" || modelUsed == "claude-2" {
			alternative = "claude-3-haiku"
		} else if modelUsed == "claude-3-sonnet" {
			alternative = "claude-3-haiku"
		}

		currentCost := modelCosts[modelUsed]
		altCost := modelCosts[alternative]
		estimatedSavings := totalCost * ((currentCost - altCost) / currentCost)

		if estimatedSavings > 0.01 {
			applyAction := map[string]interface{}{
				"action":     "switch_model",
				"from_model": modelUsed,
				"to_model":   alternative,
				"flow_id":    flowID,
			}
			applyJSON, _ := json.Marshal(applyAction)

			suggestions = append(suggestions, Suggestion{
				Type:             "model_switch",
				Title:            fmt.Sprintf("Switch from %s to %s for flow %s", modelUsed, alternative, flowID),
				Description:      fmt.Sprintf("Flow '%s' has made %d calls using %s. Switching to %s could save approximately $%.4f.", flowID, callCount, modelUsed, alternative, estimatedSavings),
				EstimatedSavings: estimatedSavings,
				SavingsUnit:      "USD",
				TargetFlowID:     flowID,
				ApplyAction:      string(applyJSON),
			})
		}
	}

	return suggestions, nil
}

// detectLongPrompts identifies entries with unusually high input token counts.
func detectLongPrompts(db *sql.DB) ([]Suggestion, error) {
	suggestions := []Suggestion{}

	query := `
		SELECT flow_id, agent_role, AVG(input_tokens) as avg_input, COUNT(*) as call_count
		FROM token_ledger
		WHERE input_tokens > 2000
		GROUP BY flow_id, agent_role
		HAVING call_count >= 2
		ORDER BY avg_input DESC
		LIMIT 5
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var flowID, agentRole string
		var avgInput float64
		var callCount int

		if err := rows.Scan(&flowID, &agentRole, &avgInput, &callCount); err != nil {
			return nil, err
		}

		potentialReduction := 0.3
		estimatedTokenSavings := avgInput * potentialReduction * float64(callCount)

		applyAction := map[string]interface{}{
			"action":           "optimize_prompt",
			"flow_id":          flowID,
			"agent_role":       agentRole,
			"target_reduction": potentialReduction,
		}
		applyJSON, _ := json.Marshal(applyAction)

		suggestions = append(suggestions, Suggestion{
			Type:             "prompt_optimization",
			Title:            fmt.Sprintf("Optimize long prompts in flow %s", flowID),
			Description:      fmt.Sprintf("Agent '%s' in flow '%s' is using an average of %.0f input tokens (%d calls). Consider condensing prompts to reduce token usage.", agentRole, flowID, avgInput, callCount),
			EstimatedSavings: estimatedTokenSavings,
			SavingsUnit:      "tokens",
			TargetFlowID:     flowID,
			ApplyAction:      string(applyJSON),
		})
	}

	return suggestions, nil
}

// detectFailures identifies patterns of failed executions and suggests retry strategies.
func detectFailures(db *sql.DB) ([]Suggestion, error) {
	suggestions := []Suggestion{}

	query := `
		SELECT flow_id, COUNT(*) as failure_count, SUM(total_cost_usd) as wasted_cost
		FROM token_ledger
		WHERE status = 'FAILED'
		GROUP BY flow_id
		HAVING failure_count >= 2
		ORDER BY wasted_cost DESC
		LIMIT 5
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var flowID string
		var failureCount int
		var wastedCost float64

		if err := rows.Scan(&flowID, &failureCount, &wastedCost); err != nil {
			return nil, err
		}

		applyAction := map[string]interface{}{
			"action":      "implement_retry",
			"flow_id":     flowID,
			"strategy":    "exponential_backoff",
			"max_retries": 3,
		}
		applyJSON, _ := json.Marshal(applyAction)

		suggestions = append(suggestions, Suggestion{
			Type:             "retry_strategy",
			Title:            fmt.Sprintf("Implement retry logic for flow %s", flowID),
			Description:      fmt.Sprintf("Flow '%s' has experienced %d failures, wasting $%.4f. Implementing exponential backoff retry logic could help recover from transient failures.", flowID, failureCount, wastedCost),
			EstimatedSavings: wastedCost * 0.5,
			SavingsUnit:      "USD",
			TargetFlowID:     flowID,
			ApplyAction:      string(applyJSON),
		})
	}

	return suggestions, nil
}
