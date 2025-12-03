package llm

import (
	"fmt"

	"github.com/mikejsmith1985/forge-orchestrator/internal/agents"
)

// ProviderType defines the supported LLM providers.
type ProviderType string

const (
	ProviderAnthropic ProviderType = "Anthropic"
	ProviderOpenAI    ProviderType = "OpenAI"
)

// LLMResponse holds the result of an LLM generation.
type LLMResponse struct {
	Content      string
	InputTokens  int
	OutputTokens int
	Cost         float64
}

// LLMProvider is the interface that specific provider clients must implement.
type LLMProvider interface {
	Send(systemPrompt, userPrompt, apiKey string) (string, int, int, error)
}

// Gateway handles routing prompts to the appropriate provider.
type Gateway struct {
	anthropicClient LLMProvider
	openAIClient    LLMProvider
}

// NewGateway creates a new Gateway with initialized clients.
func NewGateway() *Gateway {
	return &Gateway{
		anthropicClient: &AnthropicClient{},
		openAIClient:    &OpenAIClient{},
	}
}

// ExecutePrompt routes the prompt to the specified provider and calculates cost.
// It selects the system prompt based on the agentRole.
func (g *Gateway) ExecutePrompt(agentRole, userPrompt, apiKey string, provider ProviderType) (*LLMResponse, error) {
	systemPrompt, err := agents.GetAgentPrompt(agentRole)
	if err != nil {
		return nil, err
	}

	var content string
	var inputTokens, outputTokens int
	var sendErr error

	switch provider {
	case ProviderAnthropic:
		content, inputTokens, outputTokens, sendErr = g.anthropicClient.Send(systemPrompt, userPrompt, apiKey)
	case ProviderOpenAI:
		content, inputTokens, outputTokens, sendErr = g.openAIClient.Send(systemPrompt, userPrompt, apiKey)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	if sendErr != nil {
		return nil, sendErr
	}

	cost := calculateCost(provider, inputTokens, outputTokens)

	return &LLMResponse{
		Content:      content,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		Cost:         cost,
	}, nil
}

// calculateCost estimates the cost based on provider pricing (as of late 2024/2025).
// Educational Comment: Token counting and cost estimation are crucial for budget management in LLM apps.
// We use hardcoded rates here, but in production, these should be configurable.
func calculateCost(provider ProviderType, input, output int) float64 {
	var inputRate, outputRate float64

	switch provider {
	case ProviderAnthropic:
		// Claude 3.5 Sonnet: ~$3.00/1M input, ~$15.00/1M output
		inputRate = 3.00 / 1_000_000
		outputRate = 15.00 / 1_000_000
	case ProviderOpenAI:
		// GPT-4o: ~$5.00/1M input, ~$15.00/1M output
		inputRate = 5.00 / 1_000_000
		outputRate = 15.00 / 1_000_000
	}

	return (float64(input) * inputRate) + (float64(output) * outputRate)
}
