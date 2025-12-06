package server

import (
	"log"
	"net/http"
	"os"
	"strings"
)

// Default allowed origins for development
var defaultAllowedOrigins = []string{
	"http://localhost:8080",
	"http://localhost:8081", // Playwright dev server
	"http://localhost:5173",
	"http://127.0.0.1:8080",
	"http://127.0.0.1:8081", // Playwright dev server
	"http://127.0.0.1:5173",
	"https://localhost:8080",
	"https://localhost:8081", // Playwright dev server
	"https://localhost:5173",
	"https://127.0.0.1:8080",
	"https://127.0.0.1:8081", // Playwright dev server
	"https://127.0.0.1:5173",
}

// allowedOrigins holds the set of allowed origins for CORS
var allowedOrigins map[string]bool

// InitCORS initializes the CORS configuration from environment variables
func InitCORS() []string {
	envOrigins := os.Getenv("FORGE_ALLOWED_ORIGINS")
	var origins []string

	if envOrigins != "" {
		origins = strings.Split(envOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}
	} else {
		origins = defaultAllowedOrigins
	}

	allowedOrigins = make(map[string]bool)
	for _, o := range origins {
		allowedOrigins[o] = true
	}

	log.Printf("CORS: Allowed origins: %v", origins)
	return origins
}

// IsAllowedOrigin checks if the given origin is in the whitelist
func IsAllowedOrigin(origin string) bool {
	if origin == "" {
		return true // Same-origin requests don't send Origin header
	}
	return allowedOrigins[origin]
}

// CORSMiddleware creates a middleware that handles CORS
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// If origin is provided and not allowed, return 403
		if origin != "" && !IsAllowedOrigin(origin) {
			log.Printf("CORS: Blocked request from origin: %s", origin)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Set CORS headers for allowed origins
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Forge-Api-Key")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
		}

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
