// Package llm provides interfaces and implementations for communicating with Large Language Models.
// This test file verifies that the StubAdapter works correctly.
package llm

import (
	"testing"
)

// TestStubAdapterGenerateReturnsZeroCostResponse verifies that the stub_adapter.Generate()
// method returns the expected zero-cost, hardcoded response.
// This is the main validation for Contract 3.
func TestStubAdapterGenerateReturnsZeroCostResponse(t *testing.T) {
	// Create a new stub adapter.
	stub := NewStubAdapter()

	// Create a test configuration.
	config := LLMConfig{
		Provider:        "TestProvider",
		Model:           "test-model",
		PrimaryCostUnit: CostUnitToken,
	}

	// Call Generate with test parameters.
	response, err := stub.Generate(config, "system prompt", "user prompt", "fake-api-key")

	// Verify no error occurred.
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify response is not nil.
	if response == nil {
		t.Fatal("Expected response to not be nil")
	}

	// Verify token counts are zero (stub doesn't process real tokens).
	if response.InputTokens != 0 {
		t.Errorf("Expected InputTokens to be 0, got %d", response.InputTokens)
	}
	if response.OutputTokens != 0 {
		t.Errorf("Expected OutputTokens to be 0, got %d", response.OutputTokens)
	}

	// Verify cost is zero (stub is free to use).
	if response.Cost != 0 {
		t.Errorf("Expected Cost to be 0, got %f", response.Cost)
	}

	// Verify content is not empty (should have hardcoded response).
	if response.Content == "" {
		t.Error("Expected Content to not be empty")
	}
}

// TestStubAdapterPrimaryCostUnitIsReadable verifies that the PrimaryCostUnit field
// can be read successfully from the configuration struct.
// This is part of Contract 3 validation.
func TestStubAdapterPrimaryCostUnitIsReadable(t *testing.T) {
	// Create a new stub adapter.
	stub := NewStubAdapter()

	// Get the configuration.
	config := stub.GetConfig()

	// Verify PrimaryCostUnit is set and readable.
	if config.PrimaryCostUnit == "" {
		t.Error("Expected PrimaryCostUnit to be set, got empty string")
	}

	// Verify it's set to TOKEN (the default for the stub).
	if config.PrimaryCostUnit != CostUnitToken {
		t.Errorf("Expected PrimaryCostUnit to be %q, got %q", CostUnitToken, config.PrimaryCostUnit)
	}
}

// TestStubAdapterGetBudgetStatus verifies that budget status returns sensible defaults.
func TestStubAdapterGetBudgetStatus(t *testing.T) {
	stub := NewStubAdapter()

	status := stub.GetBudgetStatus()

	// Verify no spending has occurred.
	if status.TotalSpentUSD != 0 {
		t.Errorf("Expected TotalSpentUSD to be 0, got %f", status.TotalSpentUSD)
	}

	// Verify not over budget.
	if status.IsOverBudget {
		t.Error("Expected IsOverBudget to be false for stub")
	}

	// Verify daily limit is set (stub has very high limits).
	if status.DailyLimitUSD <= 0 {
		t.Error("Expected DailyLimitUSD to be greater than 0")
	}
}

// TestLLMConfigPrimaryCostUnitValues verifies that the PrimaryCostUnit constants are defined.
func TestLLMConfigPrimaryCostUnitValues(t *testing.T) {
	// Verify TOKEN constant is defined correctly.
	if CostUnitToken != "TOKEN" {
		t.Errorf("Expected CostUnitToken to be 'TOKEN', got %q", CostUnitToken)
	}

	// Verify PROMPT constant is defined correctly.
	if CostUnitPrompt != "PROMPT" {
		t.Errorf("Expected CostUnitPrompt to be 'PROMPT', got %q", CostUnitPrompt)
	}
}

// TestLLMConfigCanBeCreated verifies that LLMConfig struct can be created with all fields.
func TestLLMConfigCanBeCreated(t *testing.T) {
	config := LLMConfig{
		Provider:        "Anthropic",
		Model:           "claude-3-5-sonnet",
		APIEndpoint:     "https://api.anthropic.com/v1/messages",
		PrimaryCostUnit: CostUnitToken,
		InputRate:       3.00,
		OutputRate:      15.00,
		PromptRate:      0.0,
		MaxTokens:       200000,
		Temperature:     0.7,
	}

	// Verify all fields are accessible.
	if config.Provider != "Anthropic" {
		t.Errorf("Expected Provider to be 'Anthropic', got %q", config.Provider)
	}
	if config.PrimaryCostUnit != CostUnitToken {
		t.Errorf("Expected PrimaryCostUnit to be TOKEN, got %q", config.PrimaryCostUnit)
	}
	if config.InputRate != 3.00 {
		t.Errorf("Expected InputRate to be 3.00, got %f", config.InputRate)
	}
}

// TestBudgetStatusCanBeCreated verifies that BudgetStatus struct can be created with all fields.
func TestBudgetStatusCanBeCreated(t *testing.T) {
	status := BudgetStatus{
		TotalSpentUSD:        10.50,
		DailyLimitUSD:        50.00,
		DailySpentUSD:        5.25,
		MonthlyLimitUSD:      500.00,
		MonthlySpentUSD:      150.00,
		TotalTokensUsed:      1000000,
		RemainingDailyBudget: 44.75,
		IsOverBudget:         false,
	}

	// Verify all fields are accessible.
	if status.TotalSpentUSD != 10.50 {
		t.Errorf("Expected TotalSpentUSD to be 10.50, got %f", status.TotalSpentUSD)
	}
	if status.TotalTokensUsed != 1000000 {
		t.Errorf("Expected TotalTokensUsed to be 1000000, got %d", status.TotalTokensUsed)
	}
	if status.IsOverBudget != false {
		t.Error("Expected IsOverBudget to be false")
	}
}

// TestStubAdapterImplementsLLMGatewayInterface verifies that StubAdapter
// properly implements the LLMGatewayInterface interface.
func TestStubAdapterImplementsLLMGatewayInterface(t *testing.T) {
	// This test verifies at compile time that StubAdapter implements LLMGatewayInterface.
	// If it doesn't, this won't compile.
	var _ LLMGatewayInterface = (*StubAdapter)(nil)
	var _ LLMGatewayInterface = NewStubAdapter()

	// If we got here, the interface is properly implemented.
	t.Log("StubAdapter correctly implements the LLMGatewayInterface interface")
}
