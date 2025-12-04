package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	BackendURL  = "http://localhost:8080"
	FrontendURL = "http://localhost:8081"
)

func main() {
	fmt.Println("🔥 Starting Smoke Tests...")

	// 1. Check Health (Direct Backend)
	checkEndpoint(BackendURL+"/api/health", "GET", nil, 200)

	// 2. Check Health (Via Frontend Proxy)
	checkEndpoint(FrontendURL+"/api/health", "GET", nil, 200)

	// 3. Create Ledger Entry
	ledgerPayload := map[string]interface{}{
		"model_used":     "smoke-test-model",
		"input_tokens":   10,
		"output_tokens":  20,
		"total_cost_usd": 0.001,
		"status":         "success",
	}
	checkEndpoint(BackendURL+"/api/ledger", "POST", ledgerPayload, 201)

	// 4. Get Ledger Entries
	checkEndpoint(BackendURL+"/api/ledger", "GET", nil, 200)

	// 5. Create Command
	cmdPayload := map[string]interface{}{
		"name":        "Smoke Test Command",
		"command":     "Echo hello",
		"description": "Created by smoke test",
	}
	resp := checkEndpoint(BackendURL+"/api/commands", "POST", cmdPayload, 201)

	// Extract ID for execution
	var cmd map[string]interface{}
	json.Unmarshal(resp, &cmd)
	cmdID := int(cmd["id"].(float64))
	fmt.Printf("   Created Command ID: %d\n", cmdID)

	// 6. Execute Command (Expect 500 or 401 if no key, but NOT HTML)
	// We expect this to fail gracefully (JSON error) not HTML error
	// If we have no keys set, it returns 401 or 500.
	// Let's just check it returns JSON.
	runResp := checkEndpoint(BackendURL+fmt.Sprintf("/api/commands/%d/run", cmdID), "POST", nil, 0) // 0 means don't check status code strict
	if bytes.Contains(runResp, []byte("<!doctype")) || bytes.Contains(runResp, []byte("<html")) {
		fmt.Println("❌ FAILED: Received HTML response for Command Execution!")
		os.Exit(1)
	} else {
		fmt.Println("✅ Command Execution returned non-HTML response (Pass)")
	}

	fmt.Println("✨ All Smoke Tests Passed!")
}

func checkEndpoint(url, method string, payload interface{}, expectedStatus int) []byte {
	fmt.Printf("👉 Checking %s %s... ", method, url)

	var body io.Reader
	if payload != nil {
		jsonBytes, _ := json.Marshal(payload)
		body = bytes.NewBuffer(jsonBytes)
	}

	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Check for HTML (The "Unexpected token <" killer)
	if bytes.Contains(respBody, []byte("<!doctype")) || bytes.Contains(respBody, []byte("<html")) {
		fmt.Printf("❌ FAILED: Received HTML instead of JSON!\n")
		fmt.Println(string(respBody))
		os.Exit(1)
	}

	if expectedStatus != 0 && resp.StatusCode != expectedStatus {
		fmt.Printf("❌ FAILED: Expected status %d, got %d\n", expectedStatus, resp.StatusCode)
		fmt.Println(string(respBody))
		os.Exit(1)
	}

	fmt.Println("✅ OK")
	return respBody
}
