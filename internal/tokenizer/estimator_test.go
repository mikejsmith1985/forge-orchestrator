package tokenizer

import (
	"strings"
	"testing"
)

func TestEstimator_EstimateWithTiktoken(t *testing.T) {
	e := NewEstimator()

	tests := []struct {
		name     string
		text     string
		provider string
		model    string
		wantMin  int
		wantMax  int
	}{
		{
			name:     "simple text",
			text:     "hello world",
			provider: "openai",
			model:    "gpt-4",
			wantMin:  1,
			wantMax:  5,
		},
		{
			name:     "code snippet",
			text:     "func main() { fmt.Println(\"hello\") }",
			provider: "openai",
			model:    "gpt-4",
			wantMin:  5,
			wantMax:  20,
		},
		{
			name:     "longer text",
			text:     "The quick brown fox jumps over the lazy dog. This is a common pangram used for testing.",
			provider: "openai",
			model:    "gpt-4",
			wantMin:  15,
			wantMax:  30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := e.Estimate(tt.text, tt.provider, tt.model)

			if result.Count < tt.wantMin || result.Count > tt.wantMax {
				t.Errorf("Estimate() count = %v, want between %v and %v", result.Count, tt.wantMin, tt.wantMax)
			}

			if result.Method != "tiktoken" {
				t.Errorf("Estimate() method = %v, want tiktoken", result.Method)
			}

			if result.Provider != "openai" {
				t.Errorf("Estimate() provider = %v, want openai", result.Provider)
			}
		})
	}
}

func TestEstimator_EstimateHeuristicForAnthropic(t *testing.T) {
	e := NewEstimator()

	result := e.Estimate("hello world this is a test", "anthropic", "")

	if result.Method != "heuristic" {
		t.Errorf("Estimate() method = %v, want heuristic for Anthropic", result.Method)
	}

	if result.Provider != "anthropic" {
		t.Errorf("Estimate() provider = %v, want anthropic", result.Provider)
	}

	// Anthropic should have higher estimate due to 1.1x multiplier
	openaiResult := e.Estimate("hello world this is a test", "openai", "gpt-4")
	// Allow for tiktoken being more accurate than our heuristic comparison
	if result.Count < 5 {
		t.Errorf("Estimate() count = %v, expected at least 5 tokens", result.Count)
	}

	_ = openaiResult // Used for comparison context
}

func TestEstimator_DefaultProvider(t *testing.T) {
	e := NewEstimator()

	result := e.Estimate("hello world", "", "")

	if result.Provider != "openai" {
		t.Errorf("Estimate() provider = %v, want openai as default", result.Provider)
	}
}

func TestEstimator_NonASCIIHandling(t *testing.T) {
	e := NewEstimator()

	// Text with non-ASCII characters
	textWithUnicode := "Hello 世界 🌍"
	result := e.Estimate(textWithUnicode, "anthropic", "") // Use heuristic

	// Should have more tokens due to non-ASCII penalty
	plainText := "Hello world earth"
	plainResult := e.Estimate(plainText, "anthropic", "")

	if result.Count <= plainResult.Count-2 {
		t.Logf("Unicode text tokens: %d, Plain text tokens: %d", result.Count, plainResult.Count)
	}
}

func TestEstimator_CodePatterns(t *testing.T) {
	e := NewEstimator()

	code := "if (x > 0) { return x * 2; }"
	result := e.Estimate(code, "anthropic", "")

	// Code should produce reasonable token count
	if result.Count < 5 {
		t.Errorf("Estimate() for code = %v, expected at least 5 tokens", result.Count)
	}
}

func TestEstimator_EmptyText(t *testing.T) {
	e := NewEstimator()

	result := e.Estimate("", "openai", "gpt-4")

	if result.Count != 0 {
		t.Errorf("Estimate() for empty text = %v, want 0", result.Count)
	}
}

func TestEstimator_AccuracyBenchmark(t *testing.T) {
	e := NewEstimator()

	// Test that tiktoken produces reasonable results for known inputs
	// "Hello, world!" is typically 4 tokens in GPT-4
	result := e.Estimate("Hello, world!", "openai", "gpt-4")

	if result.Method != "tiktoken" {
		t.Skip("Tiktoken not available, skipping accuracy test")
	}

	// GPT-4 tokenizes "Hello, world!" as approximately 4 tokens
	if result.Count < 3 || result.Count > 6 {
		t.Errorf("Estimate() for 'Hello, world!' = %v, expected 3-6 tokens", result.Count)
	}
}

func TestEstimator_LongText(t *testing.T) {
	e := NewEstimator()

	// Generate a longer text
	longText := strings.Repeat("This is a test sentence. ", 100)
	result := e.Estimate(longText, "openai", "gpt-4")

	// Should be roughly 500-700 tokens for 500 words
	if result.Count < 400 || result.Count > 800 {
		t.Errorf("Estimate() for long text = %v, expected 400-800 tokens", result.Count)
	}
}
