package llm

import (
	"errors"
	"strings"
	"testing"

	"github.com/mikejsmith1985/forge-orchestrator/internal/agents"
)

// MockProvider implements LLMProvider for testing.
type MockProvider struct {
	SendFunc func(systemPrompt, userPrompt, apiKey string) (string, int, int, error)
}

func (m *MockProvider) Send(systemPrompt, userPrompt, apiKey string) (string, int, int, error) {
	return m.SendFunc(systemPrompt, userPrompt, apiKey)
}

func TestExecutePrompt_SystemPromptSelection(t *testing.T) {
	// Setup mock
	mockProvider := &MockProvider{
		SendFunc: func(systemPrompt, userPrompt, apiKey string) (string, int, int, error) {
			// Verify system prompt matches the agent role
			expectedPrompt, _ := agents.GetAgentPrompt("Architect")
			if systemPrompt != expectedPrompt {
				return "", 0, 0, errors.New("incorrect system prompt")
			}
			return "response", 10, 20, nil
		},
	}

	gateway := &Gateway{
		AnthropicClient: mockProvider,
		OpenAIClient:    mockProvider,
	}

	// Test Architect role
	_, err := gateway.ExecutePrompt("Architect", "hello", "key", ProviderAnthropic)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestExecutePrompt_TokenCountingAndCost(t *testing.T) {
	mockProvider := &MockProvider{
		SendFunc: func(systemPrompt, userPrompt, apiKey string) (string, int, int, error) {
			return "response", 1000, 1000, nil
		},
	}

	gateway := &Gateway{
		AnthropicClient: mockProvider,
		OpenAIClient:    mockProvider,
	}

	// Test Anthropic Cost
	// Input: 1000 tokens * $3/1M = $0.003
	// Output: 1000 tokens * $15/1M = $0.015
	// Total: $0.018
	resp, err := gateway.ExecutePrompt("Architect", "hello", "key", ProviderAnthropic)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedCost := 0.018
	if !floatEquals(resp.Cost, expectedCost) {
		t.Errorf("expected cost %f, got %f", expectedCost, resp.Cost)
	}

	// Test OpenAI Cost
	// Input: 1000 tokens * $5/1M = $0.005
	// Output: 1000 tokens * $15/1M = $0.015
	// Total: $0.020
	resp, err = gateway.ExecutePrompt("Architect", "hello", "key", ProviderOpenAI)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedCost = 0.020
	if !floatEquals(resp.Cost, expectedCost) {
		t.Errorf("expected cost %f, got %f", expectedCost, resp.Cost)
	}
}

func floatEquals(a, b float64) bool {
	const epsilon = 1e-9
	return (a-b) < epsilon && (b-a) < epsilon
}

func TestExecutePrompt_InvalidRole(t *testing.T) {
	gateway := NewGateway()
	_, err := gateway.ExecutePrompt("InvalidRole", "hello", "key", ProviderAnthropic)
	if err == nil {
		t.Error("expected error for invalid role, got nil")
	}
}

// ========== ERROR HANDLING TESTS ==========

func TestExecutePrompt_UnsupportedProvider(t *testing.T) {
	gateway := NewGateway()
	_, err := gateway.ExecutePrompt("Architect", "hello", "key", ProviderType("UnsupportedProvider"))
	if err == nil {
		t.Error("expected error for unsupported provider, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "unsupported provider") {
		t.Errorf("error should mention unsupported provider: %v", err)
	}
}

func TestExecutePrompt_ProviderError(t *testing.T) {
	mockProvider := &MockProvider{
		SendFunc: func(systemPrompt, userPrompt, apiKey string) (string, int, int, error) {
			return "", 0, 0, errors.New("network timeout")
		},
	}

	gateway := &Gateway{
		AnthropicClient: mockProvider,
		OpenAIClient:    mockProvider,
	}

	_, err := gateway.ExecutePrompt("Architect", "hello", "key", ProviderAnthropic)
	if err == nil {
		t.Error("expected error from provider, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "network timeout") {
		t.Errorf("error should contain provider error: %v", err)
	}
}

func TestExecutePrompt_InvalidAPIKey(t *testing.T) {
	mockProvider := &MockProvider{
		SendFunc: func(systemPrompt, userPrompt, apiKey string) (string, int, int, error) {
			if apiKey == "" {
				return "", 0, 0, errors.New("invalid API key")
			}
			return "response", 10, 20, nil
		},
	}

	gateway := &Gateway{
		AnthropicClient: mockProvider,
		OpenAIClient:    mockProvider,
	}

	_, err := gateway.ExecutePrompt("Architect", "hello", "", ProviderAnthropic)
	if err == nil {
		t.Error("expected error for empty API key, got nil")
	}
}

func TestExecutePrompt_EmptyPrompt(t *testing.T) {
	mockProvider := &MockProvider{
		SendFunc: func(systemPrompt, userPrompt, apiKey string) (string, int, int, error) {
			// Should handle empty prompt gracefully
			return "response", 0, 5, nil
		},
	}

	gateway := &Gateway{
		AnthropicClient: mockProvider,
		OpenAIClient:    mockProvider,
	}

	// Empty prompt should still succeed (edge case but valid)
	resp, err := gateway.ExecutePrompt("Architect", "", "key", ProviderAnthropic)
	if err != nil {
		t.Errorf("unexpected error for empty prompt: %v", err)
	}
	if resp == nil {
		t.Error("response should not be nil")
	}
}

func TestCalculateCost_UnknownProvider(t *testing.T) {
	// Test that unknown provider returns zero cost (doesn't panic)
	cost := calculateCost(ProviderType("Unknown"), 1000, 1000)
	if cost != 0 {
		t.Errorf("expected 0 cost for unknown provider, got %f", cost)
	}
}
