# Issue #034: Fix CORS Security Vulnerability

**Priority:** 🔴 CRITICAL  
**Estimated Tokens:** ~1,500 (Low complexity)  
**Agent Role:** Implementation

---

## 1. 🎫 Related Issue Context

**Gap Reference:** GAP-003 from v1-analysis.md

**SECURITY VULNERABILITY:** The WebSocket upgrader allows ALL origins:

```go
CheckOrigin: func(r *http.Request) bool {
    return true // ← DANGER: Allows any website to connect
}
```

Per Project Charter: "Go enforces strict CORS restrictions."

**Attack Vector:** Malicious websites can connect to localhost WebSocket and execute commands using victim's API keys.

---

## 2. 📋 Acceptance Criteria

### Backend (Go)
- [ ] Create `internal/server/middleware.go` with CORS middleware
- [ ] Read allowed origins from `FORGE_ALLOWED_ORIGINS` env variable
- [ ] Default to `http://localhost:8080,http://localhost:5173,http://127.0.0.1:8080,http://127.0.0.1:5173`
- [ ] Update WebSocket `CheckOrigin` to validate against whitelist
- [ ] Add CORS headers to all HTTP API responses
- [ ] Return 403 Forbidden for unauthorized origins

### Configuration
- [ ] Add `FORGE_ALLOWED_ORIGINS` to README documentation
- [ ] Log allowed origins on server startup

### Tests
- [ ] Unit test: Allowed origin passes
- [ ] Unit test: Blocked origin returns 403
- [ ] Unit test: WebSocket upgrade fails for blocked origin
- [ ] E2E test: Verify frontend still works (same-origin)

---

## 3. 📊 Token Efficiency Strategy

- Single new file (middleware.go ~60 lines)
- Minimal changes to existing files
- Standard Go HTTP middleware pattern

---

## 4. 🏗️ Technical Specification

### Middleware Implementation
```go
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
    originSet := make(map[string]bool)
    for _, o := range allowedOrigins {
        originSet[o] = true
    }
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            
            if origin != "" && !originSet[origin] {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }
            
            if origin != "" {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
                w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Forge-Api-Key")
            }
            
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

### WebSocket Origin Check
```go
CheckOrigin: func(r *http.Request) bool {
    origin := r.Header.Get("Origin")
    if origin == "" {
        return true // Same-origin requests don't send Origin header
    }
    return isAllowedOrigin(origin)
}
```

### Environment Variable
```bash
# Default (development)
FORGE_ALLOWED_ORIGINS=http://localhost:8080,http://localhost:5173

# Production example
FORGE_ALLOWED_ORIGINS=https://myapp.example.com
```

---

## 5. 📁 Files to Create/Modify

| Action | File |
|--------|------|
| CREATE | `internal/server/middleware.go` |
| MODIFY | `internal/server/websocket.go` (update CheckOrigin) |
| MODIFY | `internal/server/routes.go` (wrap with CORS middleware) |
| MODIFY | `main.go` (read env variable) |
| MODIFY | `README.md` (document FORGE_ALLOWED_ORIGINS) |

---

## 6. ✅ Definition of Done

1. Server logs allowed origins on startup
2. Cross-origin requests from unlisted origins return 403
3. Frontend continues to work (localhost is allowed)
4. WebSocket connections from unlisted origins are rejected
5. All tests pass
