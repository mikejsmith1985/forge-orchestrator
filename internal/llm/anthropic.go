package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DefaultAnthropicEndpoint is the production API endpoint for Anthropic.
const DefaultAnthropicEndpoint = "https://api.anthropic.com/v1/messages"

// DefaultTimeoutSeconds is the default HTTP client timeout.
const DefaultTimeoutSeconds = 30

// AnthropicClient implements the LLMProvider interface for Anthropic.
// It supports configurable endpoints and timeouts for testing and production use.
type AnthropicClient struct {
	// Endpoint is the API URL. If empty, uses DefaultAnthropicEndpoint.
	Endpoint string

	// TimeoutSeconds is the HTTP client timeout. If 0, uses DefaultTimeoutSeconds.
	TimeoutSeconds int
}

// getEndpoint returns the configured endpoint or the default.
func (c *AnthropicClient) getEndpoint() string {
	if c.Endpoint != "" {
		return c.Endpoint
	}
	return DefaultAnthropicEndpoint
}

// getTimeout returns the configured timeout duration.
func (c *AnthropicClient) getTimeout() time.Duration {
	if c.TimeoutSeconds > 0 {
		return time.Duration(c.TimeoutSeconds) * time.Second
	}
	return DefaultTimeoutSeconds * time.Second
}

// anthropicRequest represents the payload for the Anthropic API.
type anthropicRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system"`
	Messages  []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// anthropicResponse represents the response from the Anthropic API.
type anthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Send sends a prompt to Anthropic's Claude 3.5 Sonnet model.
// It uses configurable endpoint and timeout for testability.
func (c *AnthropicClient) Send(systemPrompt, userPrompt, apiKey string) (string, int, int, error) {
	reqBody := anthropicRequest{
		Model:     "claude-3-5-sonnet-20240620",
		MaxTokens: 4096,
		System:    systemPrompt,
		Messages: []message{
			{Role: "user", Content: userPrompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.getEndpoint(), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	client := &http.Client{
		Timeout: c.getTimeout(),
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", 0, 0, fmt.Errorf("anthropic api error (status %d): %s", resp.StatusCode, string(body))
	}

	var response anthropicResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", 0, 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return "", 0, 0, fmt.Errorf("anthropic error: %s", response.Error.Message)
	}

	if len(response.Content) == 0 {
		return "", 0, 0, fmt.Errorf("empty response content")
	}

	return response.Content[0].Text, response.Usage.InputTokens, response.Usage.OutputTokens, nil
}
