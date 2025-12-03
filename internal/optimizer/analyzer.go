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
// Educational Comment: This is the core analysis function that applies multiple heuristics to detect
// inefficiencies in token usage. It returns a slice of actionable suggestions.
func AnalyzeLedger(db *sql.DB) ([]Suggestion, error) {
	suggestions := []Suggestion{}
	nextID := 1 // Simple ID counter for suggestions

	// Analysis 1: Detect repeated high-cost model usage and suggest cheaper alternatives
	// Educational Comment: We look for flows that consistently use expensive models.
	// If a flow has multiple calls using high-cost models, we suggest switching to cheaper alternatives.
	highCostSuggestions, err := detectHighCostModels(db, &nextID)
	if err != nil {
		return nil, fmt.Errorf("failed to detect high-cost models: %w", err)
	}
	suggestions = append(suggestions, highCostSuggestions...)

	// Analysis 2: Identify long prompts that might benefit from condensing
	// Educational Comment: Large input token counts suggest verbose prompts that could be optimized.
	// We look for entries with unusually high input token counts.
	longPromptSuggestions, err := detectLongPrompts(db, &nextID)
	if err != nil {
		return nil, fmt.Errorf("failed to detect long prompts: %w", err)
	}
	suggestions = append(suggestions, longPromptSuggestions...)

	// Analysis 3: Find failed executions and suggest retry strategies
	// Educational Comment: Failed API calls waste tokens and money. We identify patterns of failures
	// and suggest implementing retry logic or switching to more reliable configurations.
	failureSuggestions, err := detectFailures(db, &nextID)
	if err != nil {
		return nil, fmt.Errorf("failed to detect failures: %w", err)
	}
	suggestions = append(suggestions, failureSuggestions...)

	return suggestions, nil
}

// detectHighCostModels identifies flows using expensive models and suggests cheaper alternatives.
func detectHighCostModels(db *sql.DB, nextID *int) ([]Suggestion, error) {
	suggestions := []Suggestion{}

	// Educational Comment: We group by flow_id and model_used to find high-cost usage patterns.
	// A threshold of $0.03 per million tokens is considered "high cost" for this analysis.
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

		// Educational Comment: Suggest a cheaper alternative based on the current model.
		// For GPT models, suggest gpt-3.5-turbo. For Claude models, suggest claude-3-haiku.
		alternative := "gpt-3.5-turbo"
		if modelUsed == "claude-3-opus" || modelUsed == "claude-2" {
			alternative = "claude-3-haiku"
		} else if modelUsed == "claude-3-sonnet" {
			alternative = "claude-3-haiku"
		}

		// Estimate savings: difference in cost per million tokens times average usage
		currentCost := modelCosts[modelUsed]
		altCost := modelCosts[alternative]
		estimatedSavings := totalCost * ((currentCost - altCost) / currentCost)

		// Only suggest if savings are meaningful (>$0.01)
		if estimatedSavings > 0.01 {
			applyAction := map[string]interface{}{
				"action":     "switch_model",
				"from_model": modelUsed,
				"to_model":   alternative,
				"flow_id":    flowID,
			}
			applyJSON, _ := json.Marshal(applyAction)

			suggestions = append(suggestions, Suggestion{
				ID:               *nextID,
				Type:             "model_switch",
				Title:            fmt.Sprintf("Switch from %s to %s for flow %s", modelUsed, alternative, flowID),
				Description:      fmt.Sprintf("Flow '%s' has made %d calls using %s. Switching to %s could save approximately $%.4f.", flowID, callCount, modelUsed, alternative, estimatedSavings),
				EstimatedSavings: estimatedSavings,
				SavingsUnit:      "USD",
				TargetFlowID:     flowID,
				ApplyAction:      string(applyJSON),
			})
			*nextID++
		}
	}

	return suggestions, nil
}

// detectLongPrompts identifies entries with unusually high input token counts.
func detectLongPrompts(db *sql.DB, nextID *int) ([]Suggestion, error) {
	suggestions := []Suggestion{}

	// Educational Comment: We look for entries where input_tokens > 2000, which suggests
	// a very long prompt that might benefit from condensing or restructuring.
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

		// Educational Comment: Estimate potential savings by assuming a 30% reduction in input tokens
		// through prompt optimization techniques like removing redundancy and using more concise language.
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
			ID:               *nextID,
			Type:             "prompt_optimization",
			Title:            fmt.Sprintf("Optimize long prompts in flow %s", flowID),
			Description:      fmt.Sprintf("Agent '%s' in flow '%s' is using an average of %.0f input tokens (%d calls). Consider condensing prompts to reduce token usage.", agentRole, flowID, avgInput, callCount),
			EstimatedSavings: estimatedTokenSavings,
			SavingsUnit:      "tokens",
			TargetFlowID:     flowID,
			ApplyAction:      string(applyJSON),
		})
		*nextID++
	}

	return suggestions, nil
}

// detectFailures identifies patterns of failed executions and suggests retry strategies.
func detectFailures(db *sql.DB, nextID *int) ([]Suggestion, error) {
	suggestions := []Suggestion{}

	// Educational Comment: Failed calls are wasteful because they consume tokens without producing results.
	// We detect flows with multiple failures and suggest implementing retry logic.
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

		// Educational Comment: Implementing exponential backoff retry logic can help recover from
		// transient failures without wasting additional tokens on failed attempts.
		applyAction := map[string]interface{}{
			"action":      "implement_retry",
			"flow_id":     flowID,
			"strategy":    "exponential_backoff",
			"max_retries": 3,
		}
		applyJSON, _ := json.Marshal(applyAction)

		suggestions = append(suggestions, Suggestion{
			ID:               *nextID,
			Type:             "retry_strategy",
			Title:            fmt.Sprintf("Implement retry logic for flow %s", flowID),
			Description:      fmt.Sprintf("Flow '%s' has experienced %d failures, wasting $%.4f. Implementing exponential backoff retry logic could help recover from transient failures.", flowID, failureCount, wastedCost),
			EstimatedSavings: wastedCost * 0.5, // Assume retry logic could save 50% of wasted costs
			SavingsUnit:      "USD",
			TargetFlowID:     flowID,
			ApplyAction:      string(applyJSON),
		})
		*nextID++
	}

	return suggestions, nil
}
