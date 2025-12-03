package llm

import (
	"errors"
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
		anthropicClient: mockProvider,
		openAIClient:    mockProvider,
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
		anthropicClient: mockProvider,
		openAIClient:    mockProvider,
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
