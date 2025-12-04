package integration

import (
	"testing"

	"github.com/mikejsmith1985/forge-orchestrator/internal/llm"
)

// Educational Comment: MockLLMProvider simulates an LLM provider (like Anthropic or OpenAI)
// for testing purposes. It allows us to verify that the Gateway sends the correct
// requests and handles responses appropriately without making real network calls.
type MockLLMProvider struct {
	ResponseContent  string
	InputTokens      int
	OutputTokens     int
	Err              error
	LastSystemPrompt string
	LastUserPrompt   string
}

func (m *MockLLMProvider) Send(systemPrompt, userPrompt, apiKey string) (string, int, int, error) {
	m.LastSystemPrompt = systemPrompt
	m.LastUserPrompt = userPrompt
	return m.ResponseContent, m.InputTokens, m.OutputTokens, m.Err
}

// Educational Comment: TestExecutePrompt verifies the Gateway's core functionality:
// 1. Routing to the correct provider.
// 2. Sending the correct prompts (System + User).
// 3. returning the expected response and cost.
func TestExecutePrompt(t *testing.T) {
	// Setup mock providers
	mockAnthropic := &MockLLMProvider{
		ResponseContent: "Anthropic response",
		InputTokens:     100,
		OutputTokens:    50,
	}
	mockOpenAI := &MockLLMProvider{
		ResponseContent: "OpenAI response",
		InputTokens:     200,
		OutputTokens:    100,
	}

	// Create a gateway and inject mock clients
	gateway := &llm.Gateway{
		AnthropicClient: mockAnthropic,
		OpenAIClient:    mockOpenAI,
	}

	// Test Case 1: Anthropic Provider
	t.Run("Anthropic Routing", func(t *testing.T) {
		resp, err := gateway.ExecutePrompt("Optimizer", "Optimize this code", "fake-key", llm.ProviderAnthropic)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resp.Content != "Anthropic response" {
			t.Errorf("expected content 'Anthropic response', got '%s'", resp.Content)
		}

		// Verify prompts were passed correctly
		// Note: We need to know what the "optimizer" system prompt is to verify it fully.
		// For now, we just check that *some* system prompt was passed.
		if mockAnthropic.LastSystemPrompt == "" {
			t.Error("expected system prompt to be set")
		}
		if mockAnthropic.LastUserPrompt != "Optimize this code" {
			t.Errorf("expected user prompt 'Optimize this code', got '%s'", mockAnthropic.LastUserPrompt)
		}

		// Verify cost calculation (approximate check)
		// Anthropic Sonnet: $3/1M input, $15/1M output
		// 100 input * 3/1M = 0.0003
		// 50 output * 15/1M = 0.00075
		// Total = 0.00105
		expectedCost := (100.0 * 3.0 / 1000000.0) + (50.0 * 15.0 / 1000000.0)
		diff := resp.Cost - expectedCost
		if diff < 0 {
			diff = -diff
		}
		if diff > 0.0000001 {
			t.Errorf("expected cost %f, got %f", expectedCost, resp.Cost)
		}
	})

	// Test Case 2: OpenAI Provider
	t.Run("OpenAI Routing", func(t *testing.T) {
		resp, err := gateway.ExecutePrompt("Optimizer", "Fix this bug", "fake-key", llm.ProviderOpenAI)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resp.Content != "OpenAI response" {
			t.Errorf("expected content 'OpenAI response', got '%s'", resp.Content)
		}

		if mockOpenAI.LastUserPrompt != "Fix this bug" {
			t.Errorf("expected user prompt 'Fix this bug', got '%s'", mockOpenAI.LastUserPrompt)
		}
	})
}
