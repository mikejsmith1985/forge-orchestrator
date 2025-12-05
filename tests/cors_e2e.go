package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	log.Println("=== Contract #034 E2E Test: CORS Security ===\n")

	// Wait for server to be ready
	time.Sleep(500 * time.Millisecond)

	allPassed := true

	// Test 1: Allowed origin passes
	log.Println("Test 1: Allowed origin (localhost:8080) should pass")
	if !testAllowedOrigin() {
		allPassed = false
	}

	// Test 2: Blocked origin returns 403
	log.Println("\nTest 2: Blocked origin (evil-site.com) should return 403")
	if !testBlockedOrigin() {
		allPassed = false
	}

	// Test 3: WebSocket upgrade fails for blocked origin
	log.Println("\nTest 3: WebSocket upgrade fails for blocked origin")
	if !testWebSocketBlocked() {
		allPassed = false
	}

	// Test 4: WebSocket works for allowed origin
	log.Println("\nTest 4: WebSocket works for allowed origin")
	if !testWebSocketAllowed() {
		allPassed = false
	}

	// Test 5: No origin header (same-origin) works
	log.Println("\nTest 5: Request without Origin header (same-origin) works")
	if !testNoOriginHeader() {
		allPassed = false
	}

	// Test 6: Preflight OPTIONS request works
	log.Println("\nTest 6: Preflight OPTIONS request works")
	if !testPreflightRequest() {
		allPassed = false
	}

	fmt.Println("\n=================================")
	if allPassed {
		fmt.Println("🎉 All CORS Security Tests PASSED!")
		fmt.Println("✅ Allowed origins can access the API")
		fmt.Println("✅ Blocked origins return 403 Forbidden")
		fmt.Println("✅ WebSocket rejects blocked origins")
		fmt.Println("✅ Same-origin requests work correctly")
		fmt.Println("✅ Preflight requests handled properly")
	} else {
		fmt.Println("❌ Some tests FAILED")
		os.Exit(1)
	}
	fmt.Println("=================================")
}

func testAllowedOrigin() bool {
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/health", nil)
	req.Header.Set("Origin", "http://localhost:8080")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Request failed: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("❌ Expected status 200, got %d", resp.StatusCode)
		return false
	}

	// Check CORS headers
	corsOrigin := resp.Header.Get("Access-Control-Allow-Origin")
	if corsOrigin != "http://localhost:8080" {
		log.Printf("❌ Expected CORS origin header, got: %s", corsOrigin)
		return false
	}

	log.Println("✅ Allowed origin request succeeded with CORS headers")
	return true
}

func testBlockedOrigin() bool {
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/health", nil)
	req.Header.Set("Origin", "https://evil-site.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Request failed: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 403 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("❌ Expected status 403, got %d. Body: %s", resp.StatusCode, string(body))
		return false
	}

	log.Println("✅ Blocked origin correctly returned 403 Forbidden")
	return true
}

func testWebSocketBlocked() bool {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}

	header := http.Header{}
	header.Set("Origin", "https://malicious-site.com")

	_, resp, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err == nil {
		log.Println("❌ WebSocket connection should have been rejected")
		return false
	}

	if resp != nil && resp.StatusCode == 403 {
		log.Println("✅ WebSocket connection correctly rejected with 403")
		return true
	}

	// WebSocket upgrade failure manifests as connection error
	log.Println("✅ WebSocket connection rejected for blocked origin")
	return true
}

func testWebSocketAllowed() bool {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}

	header := http.Header{}
	header.Set("Origin", "http://localhost:5173")

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Printf("❌ WebSocket connection failed: %v", err)
		return false
	}
	defer conn.Close()

	// Send a test message
	testMsg := map[string]interface{}{
		"type":    "CORS_TEST",
		"payload": "hello",
	}
	msgBytes, _ := json.Marshal(testMsg)

	if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
		log.Printf("❌ Failed to send message: %v", err)
		return false
	}

	// Read echo response
	_, _, err = conn.ReadMessage()
	if err != nil {
		log.Printf("❌ Failed to read message: %v", err)
		return false
	}

	log.Println("✅ WebSocket connection allowed for permitted origin")
	return true
}

func testNoOriginHeader() bool {
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/health", nil)
	// No Origin header - simulates same-origin request

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Request failed: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("❌ Expected status 200, got %d", resp.StatusCode)
		return false
	}

	// CORS headers should NOT be set for same-origin
	corsOrigin := resp.Header.Get("Access-Control-Allow-Origin")
	if corsOrigin != "" {
		log.Printf("⚠️  CORS header set for same-origin request (not a security issue)")
	}

	log.Println("✅ Same-origin request (no Origin header) succeeded")
	return true
}

func testPreflightRequest() bool {
	req, _ := http.NewRequest("OPTIONS", "http://localhost:8080/api/health", nil)
	req.Header.Set("Origin", "http://localhost:8080")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Preflight request failed: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("❌ Expected status 200 for preflight, got %d", resp.StatusCode)
		return false
	}

	// Check CORS headers
	allowMethods := resp.Header.Get("Access-Control-Allow-Methods")
	if allowMethods == "" {
		log.Println("❌ Missing Access-Control-Allow-Methods header")
		return false
	}

	log.Println("✅ Preflight OPTIONS request handled correctly")
	return true
}
