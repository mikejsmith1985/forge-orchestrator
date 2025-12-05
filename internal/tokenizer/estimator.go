package tokenizer

import (
	"regexp"
	"strings"

	"github.com/pkoukk/tiktoken-go"
)

// EstimationResult contains the token estimation details
type EstimationResult struct {
	Count    int    `json:"count"`
	Method   string `json:"method"`
	Provider string `json:"provider"`
	Model    string `json:"model,omitempty"`
}

// Estimator handles token estimation with multiple methods
type Estimator struct{}

// NewEstimator creates a new token estimator
func NewEstimator() *Estimator {
	return &Estimator{}
}

// Estimate returns an accurate token count using tiktoken for OpenAI models
// or falls back to heuristic estimation for other providers
func (e *Estimator) Estimate(text string, provider string, model string) EstimationResult {
	// Normalize provider name
	provider = strings.ToLower(provider)
	if provider == "" {
		provider = "openai"
	}

	// Try tiktoken for OpenAI models
	if provider == "openai" {
		if model == "" {
			model = "gpt-4"
		}
		count, err := e.estimateWithTiktoken(text, model)
		if err == nil {
			return EstimationResult{
				Count:    count,
				Method:   "tiktoken",
				Provider: provider,
				Model:    model,
			}
		}
	}

	// Fall back to heuristic for Anthropic or when tiktoken fails
	count := e.estimateHeuristic(text, provider)
	return EstimationResult{
		Count:    count,
		Method:   "heuristic",
		Provider: provider,
	}
}

// estimateWithTiktoken uses the tiktoken library for accurate OpenAI token counts
func (e *Estimator) estimateWithTiktoken(text string, model string) (int, error) {
	encoding, err := tiktoken.EncodingForModel(model)
	if err != nil {
		// Try cl100k_base as fallback (GPT-4 and GPT-3.5-turbo use this)
		encoding, err = tiktoken.GetEncoding("cl100k_base")
		if err != nil {
			return 0, err
		}
	}
	tokens := encoding.Encode(text, nil, nil)
	return len(tokens), nil
}

// estimateHeuristic provides an improved word-based estimation
func (e *Estimator) estimateHeuristic(text string, provider string) int {
	// Base: word count
	words := strings.Fields(text)
	wordCount := len(words)

	// Code pattern adjustments (brackets, operators are usually single tokens)
	codePatterns := regexp.MustCompile(`[{}\[\]();:,<>=+\-*/&|!@#$%^~]`)
	codeTokens := len(codePatterns.FindAllString(text, -1))

	// Non-ASCII penalty (typically 2-4 tokens each)
	nonAscii := 0
	for _, r := range text {
		if r > 127 {
			nonAscii++
		}
	}

	// Base estimate: words * 1.3 (average word is ~1.3 tokens)
	// Plus code tokens * 0.5 (already counted partially in words)
	// Plus non-ASCII * 2 (multi-byte characters use more tokens)
	estimate := float64(wordCount)*1.3 + float64(codeTokens)*0.5 + float64(nonAscii)*2.0

	// Provider adjustments
	if provider == "anthropic" {
		estimate *= 1.1 // Claude tends to tokenize slightly higher
	}

	// Minimum of 1 token for non-empty text
	if estimate < 1 && len(text) > 0 {
		return 1
	}

	return int(estimate)
}
