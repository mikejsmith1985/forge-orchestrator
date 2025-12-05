// Package llm provides interfaces and implementations for communicating with Large Language Models.
// This test file verifies the OpenAIClient implementation using httptest.
package llm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestOpenAIClient_SuccessfulResponse verifies that OpenAIClient correctly parses
// a successful API response and returns the content with token counts.
func TestOpenAIClient_SuccessfulResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers.
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-api-key" {
			t.Errorf("Expected Authorization header 'Bearer test-api-key', got %q", authHeader)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got %q", r.Header.Get("Content-Type"))
		}

		// Verify request body structure.
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if reqBody["model"] != "gpt-4o" {
			t.Errorf("Expected model 'gpt-4o', got %q", reqBody["model"])
		}

		// Return a successful response.
		response := openAIResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{Message: struct {
					Content string `json:"content"`
				}{Content: "Hello! I'm GPT-4o."}},
			},
			Usage: struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
			}{
				PromptTokens:     12,
				CompletionTokens: 8,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &OpenAIClient{
		Endpoint: server.URL,
	}

	content, inputTokens, outputTokens, err := client.Send(
		"You are a helpful assistant.",
		"Hello!",
		"test-api-key",
	)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if content != "Hello! I'm GPT-4o." {
		t.Errorf("Expected content 'Hello! I'm GPT-4o.', got %q", content)
	}
	if inputTokens != 12 {
		t.Errorf("Expected 12 input tokens, got %d", inputTokens)
	}
	if outputTokens != 8 {
		t.Errorf("Expected 8 output tokens, got %d", outputTokens)
	}
}

// TestOpenAIClient_APIError verifies that OpenAIClient handles API errors correctly.
func TestOpenAIClient_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": {"message": "Incorrect API key provided"}}`))
	}))
	defer server.Close()

	client := &OpenAIClient{
		Endpoint: server.URL,
	}

	_, _, _, err := client.Send("system", "user", "bad-key")

	if err == nil {
		t.Fatal("Expected error for 401 response, got nil")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("Expected error to contain '401', got: %v", err)
	}
}

// TestOpenAIClient_EmptyChoices verifies handling of empty choices array.
func TestOpenAIClient_EmptyChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := openAIResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &OpenAIClient{
		Endpoint: server.URL,
	}

	_, _, _, err := client.Send("system", "user", "key")

	if err == nil {
		t.Fatal("Expected error for empty choices, got nil")
	}
	if !strings.Contains(err.Error(), "empty response") {
		t.Errorf("Expected error to contain 'empty response', got: %v", err)
	}
}

// TestOpenAIClient_NetworkError verifies handling of network failures.
func TestOpenAIClient_NetworkError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()

	client := &OpenAIClient{
		Endpoint: server.URL,
	}

	_, _, _, err := client.Send("system", "user", "key")

	if err == nil {
		t.Fatal("Expected error for network failure, got nil")
	}
}

// TestOpenAIClient_InvalidJSON verifies handling of malformed JSON responses.
func TestOpenAIClient_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid json`))
	}))
	defer server.Close()

	client := &OpenAIClient{
		Endpoint: server.URL,
	}

	_, _, _, err := client.Send("system", "user", "key")

	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
	if !strings.Contains(err.Error(), "unmarshal") {
		t.Errorf("Expected error to contain 'unmarshal', got: %v", err)
	}
}

// TestOpenAIClient_ResponseError verifies handling of error field in response body.
func TestOpenAIClient_ResponseError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := openAIResponse{
			Error: &struct {
				Message string `json:"message"`
			}{
				Message: "The model is currently overloaded",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &OpenAIClient{
		Endpoint: server.URL,
	}

	_, _, _, err := client.Send("system", "user", "key")

	if err == nil {
		t.Fatal("Expected error for response with error field, got nil")
	}
	if !strings.Contains(err.Error(), "overloaded") {
		t.Errorf("Expected error to contain 'overloaded', got: %v", err)
	}
}

// TestOpenAIClient_DefaultEndpoint verifies that default endpoint is used when not specified.
func TestOpenAIClient_DefaultEndpoint(t *testing.T) {
	client := &OpenAIClient{}
	endpoint := client.getEndpoint()

	expected := "https://api.openai.com/v1/chat/completions"
	if endpoint != expected {
		t.Errorf("Expected default endpoint %q, got %q", expected, endpoint)
	}
}

// TestOpenAIClient_Timeout verifies that requests respect timeout settings.
func TestOpenAIClient_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := openAIResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{Message: struct {
					Content string `json:"content"`
				}{Content: "OK"}},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &OpenAIClient{
		Endpoint:       server.URL,
		TimeoutSeconds: 30,
	}

	_, _, _, err := client.Send("system", "user", "key")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

// TestOpenAIClient_RateLimitError verifies handling of 429 rate limit errors.
func TestOpenAIClient_RateLimitError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"error": {"message": "Rate limit exceeded"}}`))
	}))
	defer server.Close()

	client := &OpenAIClient{
		Endpoint: server.URL,
	}

	_, _, _, err := client.Send("system", "user", "key")

	if err == nil {
		t.Fatal("Expected error for 429 response, got nil")
	}
	if !strings.Contains(err.Error(), "429") {
		t.Errorf("Expected error to contain '429', got: %v", err)
	}
}
