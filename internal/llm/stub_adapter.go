// Package llm provides interfaces and implementations for communicating with Large Language Models.
// This file implements a stub adapter for testing and development purposes.
// The stub returns hardcoded responses instead of making real API calls.
package llm

// StubAdapter is a non-functional LLM gateway that returns hardcoded responses.
// It's used for:
// 1. Testing without making real (expensive) API calls
// 2. Development when API keys aren't available
// 3. Preventing front-end API errors during initial development
//
// Think of it like a practice dummy for martial arts - it looks real but doesn't hit back!
type StubAdapter struct {
	// Config holds the configuration for this stub adapter.
	Config LLMConfig
}

// NewStubAdapter creates a new StubAdapter with default configuration.
// This is a simple constructor that sets up a basic stub.
func NewStubAdapter() *StubAdapter {
	return &StubAdapter{
		Config: LLMConfig{
			Provider:        "Stub",
			Model:           "stub-model-v1",
			PrimaryCostUnit: CostUnitToken,
			InputRate:       0.0,  // Free! (because it's fake)
			OutputRate:      0.0,  // Free! (because it's fake)
			MaxTokens:       4096,
			Temperature:     0.0,
		},
	}
}

// Generate returns a hardcoded, successful JSON response.
// It sets token/cost to zero because this is a stub - no real work is done.
//
// The response is designed to be valid JSON that matches what a real LLM might return.
// This prevents front-end errors during development.
func (s *StubAdapter) Generate(config LLMConfig, systemPrompt, userPrompt, apiKey string) (*LLMResponse, error) {
	// Return a hardcoded successful response.
	// The content is valid JSON that a front-end can parse without errors.
	// We set all costs to zero because this is a stub and doesn't use real resources.
	return &LLMResponse{
		Content: `{
  "status": "success",
  "message": "This is a stub response for testing purposes.",
  "data": {
    "generated": true,
    "model": "stub-model-v1",
    "note": "Replace with real LLM provider for production use."
  }
}`,
		InputTokens:  0, // Zero because we didn't process any real tokens
		OutputTokens: 0, // Zero because we didn't generate any real output
		Cost:         0, // Zero because this is free (it's fake!)
	}, nil
}

// GetBudgetStatus returns a budget status indicating unlimited usage.
// Since the stub doesn't use real resources, we report zero spending and no limits.
func (s *StubAdapter) GetBudgetStatus() BudgetStatus {
	return BudgetStatus{
		TotalSpentUSD:        0.0,
		DailyLimitUSD:        999999.0, // Effectively unlimited
		DailySpentUSD:        0.0,
		MonthlyLimitUSD:      999999.0, // Effectively unlimited
		MonthlySpentUSD:      0.0,
		TotalTokensUsed:      0,
		RemainingDailyBudget: 999999.0,
		IsOverBudget:         false,
	}
}

// GetConfig returns the configuration for this stub adapter.
// This allows code to check the PrimaryCostUnit and other settings.
func (s *StubAdapter) GetConfig() LLMConfig {
	return s.Config
}
