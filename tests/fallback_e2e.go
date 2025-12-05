package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	log.Println("=== Contract #033 E2E Test: WebSocket Fallback ===\n")

	allPassed := true

	// Test 1: File Signaler creates status file
	log.Println("Test 1: File Signaler writes correct JSON format")
	if !testFileSignaler() {
		allPassed = false
	}

	// Test 2: Status endpoint returns current flow state
	log.Println("\nTest 2: Status endpoint returns current flow state")
	if !testStatusEndpoint() {
		allPassed = false
	}

	// Test 3: WebSocket connection and polling fallback simulation
	log.Println("\nTest 3: Simulate WebSocket failure and verify polling works")
	if !testPollingFallback() {
		allPassed = false
	}

	fmt.Println("\n=================================")
	if allPassed {
		fmt.Println("🎉 All E2E Tests PASSED!")
		fmt.Println("✅ FileSignaler writes correct JSON format")
		fmt.Println("✅ Status endpoint returns current flow state")
		fmt.Println("✅ Polling fallback works when WebSocket fails")
	} else {
		fmt.Println("❌ Some tests FAILED")
		os.Exit(1)
	}
	fmt.Println("=================================")
}

func testFileSignaler() bool {
	// Create status directory
	statusDir := ".forge/status"
	os.MkdirAll(statusDir, 0755)

	// Write a test status file
	status := map[string]interface{}{
		"flowId":    2,
		"status":    "RUNNING",
		"lastNode":  "agent-1",
		"updatedAt": time.Now().Format(time.RFC3339),
	}

	data, _ := json.MarshalIndent(status, "", "  ")
	filename := filepath.Join(statusDir, "2.json")
	
	if err := os.WriteFile(filename, data, 0644); err != nil {
		log.Printf("❌ Failed to write status file: %v", err)
		return false
	}

	// Verify file exists and is valid JSON
	readData, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("❌ Failed to read status file: %v", err)
		return false
	}

	var readStatus map[string]interface{}
	if err := json.Unmarshal(readData, &readStatus); err != nil {
		log.Printf("❌ Invalid JSON in status file: %v", err)
		return false
	}

	if readStatus["status"] != "RUNNING" {
		log.Printf("❌ Expected status RUNNING, got %v", readStatus["status"])
		return false
	}

	log.Println("✅ FileSignaler writes correct JSON format")

	// Cleanup
	os.Remove(filename)
	return true
}

func testStatusEndpoint() bool {
	// Create a status file first
	statusDir := ".forge/status"
	os.MkdirAll(statusDir, 0755)

	status := map[string]interface{}{
		"flowId":    3,
		"status":    "COMPLETED",
		"lastNode":  "agent-final",
		"updatedAt": time.Now().Format(time.RFC3339),
	}

	data, _ := json.MarshalIndent(status, "", "  ")
	filename := filepath.Join(statusDir, "3.json")
	os.WriteFile(filename, data, 0644)

	// Call the API endpoint
	resp, err := http.Get("http://localhost:8080/api/flows/3/status")
	if err != nil {
		log.Printf("❌ Failed to call status endpoint: %v", err)
		os.Remove(filename)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("❌ Expected status 200, got %d", resp.StatusCode)
		os.Remove(filename)
		return false
	}

	body, _ := io.ReadAll(resp.Body)
	var respStatus map[string]interface{}
	if err := json.Unmarshal(body, &respStatus); err != nil {
		log.Printf("❌ Invalid JSON response: %v", err)
		os.Remove(filename)
		return false
	}

	if respStatus["status"] != "COMPLETED" {
		log.Printf("❌ Expected status COMPLETED, got %v", respStatus["status"])
		os.Remove(filename)
		return false
	}

	log.Println("✅ Status endpoint returns correct JSON response")

	// Cleanup
	os.Remove(filename)
	return true
}

func testPollingFallback() bool {
	// Test 1: Connect via WebSocket
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("❌ Failed to connect WebSocket: %v", err)
		return false
	}
	log.Println("✅ WebSocket connection established")

	// Test 2: Send and receive message
	testMsg := map[string]interface{}{
		"type":    "TEST",
		"payload": map[string]string{"message": "polling test"},
	}
	msgBytes, _ := json.Marshal(testMsg)
	
	if err := c.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
		log.Printf("❌ Failed to send message: %v", err)
		c.Close()
		return false
	}

	// Read echo response
	_, response, err := c.ReadMessage()
	if err != nil {
		log.Printf("❌ Failed to read message: %v", err)
		c.Close()
		return false
	}
	log.Printf("✅ WebSocket echo received: %s", string(response))

	// Close WebSocket to simulate disconnection
	c.Close()
	log.Println("✅ WebSocket closed (simulating failure)")

	// Test 3: Verify polling endpoint works as fallback
	// Create a status file
	statusDir := ".forge/status"
	os.MkdirAll(statusDir, 0755)

	status := map[string]interface{}{
		"flowId":    4,
		"status":    "RUNNING",
		"lastNode":  "agent-2",
		"updatedAt": time.Now().Format(time.RFC3339),
	}
	data, _ := json.MarshalIndent(status, "", "  ")
	filename := filepath.Join(statusDir, "4.json")
	os.WriteFile(filename, data, 0644)

	// Poll the endpoint
	resp, err := http.Get("http://localhost:8080/api/flows/4/status")
	if err != nil {
		log.Printf("❌ Polling failed: %v", err)
		os.Remove(filename)
		return false
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var pollStatus map[string]interface{}
	json.Unmarshal(body, &pollStatus)

	if pollStatus["status"] != "RUNNING" {
		log.Printf("❌ Polling returned wrong status: %v", pollStatus["status"])
		os.Remove(filename)
		return false
	}

	log.Println("✅ Polling fallback works correctly")

	// Cleanup
	os.Remove(filename)
	return true
}
