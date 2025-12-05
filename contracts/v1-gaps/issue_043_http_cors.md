# Issue #043: Add HTTP CORS Middleware for API Routes

**Priority:** 🟡 HIGH  
**Estimated Tokens:** ~800 (Low complexity)  
**Agent Role:** Implementation  
**Note:** This is a sub-task split from #034 to ensure single-chat scope

---

## 1. 🎫 Related Issue Context

**Gap Reference:** GAP-005 from v1-analysis.md

While Issue #034 addresses WebSocket CORS, the HTTP API endpoints have NO CORS headers at all. This is a separate concern that should be handled by middleware.

---

## 2. 📋 Acceptance Criteria

### Backend (Go)
- [ ] Create CORS middleware in `internal/server/middleware.go`
- [ ] Set `Access-Control-Allow-Origin` from allowed origins list
- [ ] Set `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- [ ] Set `Access-Control-Allow-Headers: Content-Type, X-Forge-Api-Key`
- [ ] Handle OPTIONS preflight requests with 200 OK
- [ ] Apply middleware to all `/api/*` routes

### Configuration
- [ ] Share allowed origins list with WebSocket (from #034)
- [ ] Default to localhost variants for development

### Tests
- [ ] Unit test: CORS headers present on API response
- [ ] Unit test: OPTIONS request returns 200 with correct headers
- [ ] Unit test: Cross-origin request from blocked origin returns 403

---

## 3. 📊 Token Efficiency Strategy

- Single middleware function (~40 lines)
- Reuse origin validation from #034
- Standard Go HTTP middleware pattern

---

## 4. 🏗️ Technical Specification

### CORS Middleware
```go
// internal/server/middleware.go

func CORSMiddleware(allowedOrigins map[string]bool) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            
            // If no origin, it's same-origin request
            if origin == "" {
                next.ServeHTTP(w, r)
                return
            }
            
            // Check if origin is allowed
            if !allowedOrigins[origin] {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }
            
            // Set CORS headers
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Forge-Api-Key")
            w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
            
            // Handle preflight
            if r.Method == http.MethodOptions {
                w.WriteHeader(http.StatusOK)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

### Applying to Routes
```go
// internal/server/routes.go

func (s *Server) RegisterRoutes() http.Handler {
    mux := http.NewServeMux()
    
    // ... register all routes ...
    
    // Wrap with CORS middleware
    allowedOrigins := getallowedOrigins() // From env or default
    return CORSMiddleware(allowedOrigins)(mux)
}
```

---

## 5. 📁 Files to Create/Modify

| Action | File |
|--------|------|
| CREATE | `internal/server/middleware.go` |
| MODIFY | `internal/server/routes.go` (wrap with middleware) |
| CREATE | `internal/server/middleware_test.go` |

---

## 6. ✅ Definition of Done

1. All API responses include CORS headers when Origin is set
2. OPTIONS requests return 200 with CORS headers
3. Blocked origins receive 403 Forbidden
4. Frontend continues to work (localhost allowed)
5. All tests pass
