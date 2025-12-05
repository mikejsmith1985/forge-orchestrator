// Package llm provides interfaces and implementations for communicating with Large Language Models.
// This file defines the data models used throughout the LLM subsystem.
package llm

// PrimaryCostUnit defines how costs are measured for an LLM provider.
// Different providers charge in different ways:
// - TOKEN: Charges based on the number of tokens processed (most common)
// - PROMPT: Charges per prompt/request regardless of size
type PrimaryCostUnit string

const (
	// CostUnitToken means the provider charges per token.
	// A token is roughly 4 characters of text (like a word or part of a word).
	// Most providers like OpenAI and Anthropic use this model.
	CostUnitToken PrimaryCostUnit = "TOKEN"

	// CostUnitPrompt means the provider charges per request.
	// Some providers have flat-rate pricing per API call.
	CostUnitPrompt PrimaryCostUnit = "PROMPT"
)

// LLMConfig holds the configuration for an LLM provider.
// Think of this as the "settings" you need to talk to an AI service:
// - Which company provides it? (Provider)
// - What AI model should we use? (Model)
// - How do we calculate costs? (PrimaryCostUnit)
// - What are the rate limits? (InputRate, OutputRate, PromptRate)
type LLMConfig struct {
	// Provider is the company or service providing the LLM.
	// Examples: "Anthropic", "OpenAI", "Google"
	Provider string

	// Model is the specific AI model to use.
	// Examples: "claude-3-5-sonnet", "gpt-4o", "gemini-pro"
	Model string

	// APIEndpoint is the URL where we send requests.
	// Usually something like "https://api.anthropic.com/v1/messages"
	APIEndpoint string

	// PrimaryCostUnit indicates how this provider charges.
	// TOKEN = charges per token, PROMPT = charges per request.
	PrimaryCostUnit PrimaryCostUnit

	// InputRate is the cost per input token (in USD per 1 million tokens).
	// For example, $3.00 per million tokens.
	InputRate float64

	// OutputRate is the cost per output token (in USD per 1 million tokens).
	// Output tokens are usually more expensive than input tokens.
	OutputRate float64

	// PromptRate is the cost per prompt/request (only used if PrimaryCostUnit is PROMPT).
	// This is a flat fee per API call.
	PromptRate float64

	// MaxTokens is the maximum number of tokens the model can process in one request.
	// This includes both input and output tokens.
	MaxTokens int

	// Temperature controls how "creative" the AI responses are.
	// 0.0 = very predictable, 1.0 = more creative/random
	Temperature float64
}

// BudgetStatus tracks the current spending and limits for LLM usage.
// This helps prevent unexpected costs by setting and monitoring budgets.
// Think of it like checking your bank balance and spending limit.
type BudgetStatus struct {
	// TotalSpentUSD is how much we've spent in total (in US dollars).
	// This accumulates over time as we make API calls.
	TotalSpentUSD float64

	// DailyLimitUSD is the maximum we can spend per day.
	// If we hit this limit, we should stop making API calls.
	DailyLimitUSD float64

	// DailySpentUSD is how much we've spent today.
	// This resets at midnight (or a configured time).
	DailySpentUSD float64

	// MonthlyLimitUSD is the maximum we can spend per month.
	// Helps with longer-term budget planning.
	MonthlyLimitUSD float64

	// MonthlySpentUSD is how much we've spent this month.
	// This resets at the start of each month.
	MonthlySpentUSD float64

	// TotalTokensUsed is the total number of tokens we've processed.
	// Useful for tracking usage patterns.
	TotalTokensUsed int64

	// RemainingDailyBudget is how much more we can spend today.
	// Calculated as: DailyLimitUSD - DailySpentUSD
	RemainingDailyBudget float64

	// IsOverBudget is true if we've exceeded any of our limits.
	// When true, we should stop making API calls.
	IsOverBudget bool
}

// LLMGatewayInterface defines the contract for any LLM gateway implementation.
// This interface allows us to swap between real providers and test stubs.
type LLMGatewayInterface interface {
	// Generate sends a prompt to the LLM and returns the response.
	// It takes the agent role, user prompt, API key, and provider type.
	Generate(config LLMConfig, systemPrompt, userPrompt, apiKey string) (*LLMResponse, error)

	// GetBudgetStatus returns the current spending and budget information.
	GetBudgetStatus() BudgetStatus
}
