package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DefaultOpenAIEndpoint is the production API endpoint for OpenAI.
const DefaultOpenAIEndpoint = "https://api.openai.com/v1/chat/completions"

// OpenAIClient implements the LLMProvider interface for OpenAI.
// It supports configurable endpoints and timeouts for testing and production use.
type OpenAIClient struct {
	// Endpoint is the API URL. If empty, uses DefaultOpenAIEndpoint.
	Endpoint string

	// TimeoutSeconds is the HTTP client timeout. If 0, uses DefaultTimeoutSeconds.
	TimeoutSeconds int
}

// getEndpoint returns the configured endpoint or the default.
func (c *OpenAIClient) getEndpoint() string {
	if c.Endpoint != "" {
		return c.Endpoint
	}
	return DefaultOpenAIEndpoint
}

// getTimeout returns the configured timeout duration.
func (c *OpenAIClient) getTimeout() time.Duration {
	if c.TimeoutSeconds > 0 {
		return time.Duration(c.TimeoutSeconds) * time.Second
	}
	return DefaultTimeoutSeconds * time.Second
}

// openAIRequest represents the payload for the OpenAI API.
type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openAIResponse represents the response from the OpenAI API.
type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Send sends a prompt to OpenAI's GPT-4o model.
// It uses configurable endpoint and timeout for testability.
func (c *OpenAIClient) Send(systemPrompt, userPrompt, apiKey string) (string, int, int, error) {
	reqBody := openAIRequest{
		Model: "gpt-4o",
		Messages: []openAIMessage{
			{Role: "system", Content: systemPrompt},
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

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

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
		return "", 0, 0, fmt.Errorf("openai api error (status %d): %s", resp.StatusCode, string(body))
	}

	var response openAIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", 0, 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return "", 0, 0, fmt.Errorf("openai error: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return "", 0, 0, fmt.Errorf("empty response choices")
	}

	return response.Choices[0].Message.Content, response.Usage.PromptTokens, response.Usage.CompletionTokens, nil
}
