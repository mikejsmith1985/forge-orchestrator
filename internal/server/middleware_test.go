package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestInitCORS_DefaultOrigins(t *testing.T) {
	// Clear environment variable
	os.Unsetenv("FORGE_ALLOWED_ORIGINS")

	origins := InitCORS()

	// Should have default origins (4 http + 4 https = 8)
	if len(origins) != 8 {
		t.Errorf("Expected 8 default origins, got %d", len(origins))
	}

	// Check that localhost:8080 is allowed
	if !IsAllowedOrigin("http://localhost:8080") {
		t.Error("Expected http://localhost:8080 to be allowed")
	}

	// Check that localhost:5173 is allowed
	if !IsAllowedOrigin("http://localhost:5173") {
		t.Error("Expected http://localhost:5173 to be allowed")
	}
}

func TestInitCORS_CustomOrigins(t *testing.T) {
	os.Setenv("FORGE_ALLOWED_ORIGINS", "https://example.com, https://app.example.com")
	defer os.Unsetenv("FORGE_ALLOWED_ORIGINS")

	origins := InitCORS()

	if len(origins) != 2 {
		t.Errorf("Expected 2 custom origins, got %d", len(origins))
	}

	if !IsAllowedOrigin("https://example.com") {
		t.Error("Expected https://example.com to be allowed")
	}

	if !IsAllowedOrigin("https://app.example.com") {
		t.Error("Expected https://app.example.com to be allowed")
	}

	// Localhost should now be blocked
	if IsAllowedOrigin("http://localhost:8080") {
		t.Error("Expected http://localhost:8080 to be blocked with custom origins")
	}
}

func TestIsAllowedOrigin_EmptyOrigin(t *testing.T) {
	InitCORS()

	// Empty origin (same-origin request) should be allowed
	if !IsAllowedOrigin("") {
		t.Error("Expected empty origin to be allowed")
	}
}

func TestCORSMiddleware_AllowedOrigin(t *testing.T) {
	os.Unsetenv("FORGE_ALLOWED_ORIGINS")
	InitCORS()

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest("GET", "/api/health", nil)
	req.Header.Set("Origin", "http://localhost:8080")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Check CORS headers
	if rr.Header().Get("Access-Control-Allow-Origin") != "http://localhost:8080" {
		t.Errorf("Expected Access-Control-Allow-Origin header to be set")
	}
}

func TestCORSMiddleware_BlockedOrigin(t *testing.T) {
	os.Unsetenv("FORGE_ALLOWED_ORIGINS")
	InitCORS()

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/health", nil)
	req.Header.Set("Origin", "https://evil-site.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected status 403 Forbidden, got %d", rr.Code)
	}
}

func TestCORSMiddleware_PreflightRequest(t *testing.T) {
	os.Unsetenv("FORGE_ALLOWED_ORIGINS")
	InitCORS()

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called for OPTIONS request")
	}))

	req := httptest.NewRequest("OPTIONS", "/api/health", nil)
	req.Header.Set("Origin", "http://localhost:8080")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 for OPTIONS, got %d", rr.Code)
	}

	// Check CORS headers are set
	if rr.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("Expected Access-Control-Allow-Methods header to be set")
	}

	// Check Max-Age is set for caching preflight
	if rr.Header().Get("Access-Control-Max-Age") != "86400" {
		t.Errorf("Expected Access-Control-Max-Age to be 86400, got %s", rr.Header().Get("Access-Control-Max-Age"))
	}
}

func TestCORSMiddleware_NoOrigin(t *testing.T) {
	os.Unsetenv("FORGE_ALLOWED_ORIGINS")
	InitCORS()

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest("GET", "/api/health", nil)
	// No Origin header (same-origin request)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// No CORS headers should be set for same-origin requests
	if rr.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Error("CORS headers should not be set for same-origin requests")
	}
}

func TestWebSocketOriginCheck_Allowed(t *testing.T) {
	os.Unsetenv("FORGE_ALLOWED_ORIGINS")
	InitCORS()

	req := httptest.NewRequest("GET", "/ws", nil)
	req.Header.Set("Origin", "http://localhost:5173")

	// Test the upgrader's CheckOrigin function
	allowed := upgrader.CheckOrigin(req)
	if !allowed {
		t.Error("Expected WebSocket connection from localhost:5173 to be allowed")
	}
}

func TestWebSocketOriginCheck_Blocked(t *testing.T) {
	os.Unsetenv("FORGE_ALLOWED_ORIGINS")
	InitCORS()

	req := httptest.NewRequest("GET", "/ws", nil)
	req.Header.Set("Origin", "https://malicious-site.com")

	// Test the upgrader's CheckOrigin function
	allowed := upgrader.CheckOrigin(req)
	if allowed {
		t.Error("Expected WebSocket connection from malicious-site.com to be blocked")
	}
}

func TestWebSocketOriginCheck_NoOrigin(t *testing.T) {
	os.Unsetenv("FORGE_ALLOWED_ORIGINS")
	InitCORS()

	req := httptest.NewRequest("GET", "/ws", nil)
	// No Origin header

	// Test the upgrader's CheckOrigin function
	allowed := upgrader.CheckOrigin(req)
	if !allowed {
		t.Error("Expected WebSocket connection without Origin header to be allowed")
	}
}
