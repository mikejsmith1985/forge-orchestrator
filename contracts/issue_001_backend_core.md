# Issue: Implement Core WebSocket & API Routing (Backend)

**Role**: Implementation Agent (Go)
**Priority**: Critical
**Context**: Project scaffolding is complete. We need the core `internal/server` package to handle HTTP requests and WebSocket connections for agent signaling.

## 1. Requirements
1.  **Create `internal/server/server.go`**:
    - Define a `Server` struct.
    - Method `NewServer(db *sql.DB) *Server`.
2.  **Create `internal/server/routes.go`**:
    - Define `RegisterRoutes()` method.
    - Endpoint: `GET /api/health` (Returns 200 OK).
    - Endpoint: `GET /ws` (WebSocket upgrade).
3.  **Create `internal/server/websocket.go`**:
    - Implement basic WebSocket upgrader.
    - Maintain a list of connected clients (pool).
    - Handle basic message reading/writing (echo for now).
4.  **Update `main.go`**:
    - Replace inline HTTP handler with `server.NewServer()`.

## 2. TDD & Verification Protocol
> [!IMPORTANT]
> You must follow this TDD workflow.
1.  **Create Test**: Create `internal/server/server_test.go`.
    -   Test `GET /api/health` returns 200 and JSON.
    -   Test `GET /ws` upgrades connection (mocking WS if possible, or just checking headers).
2.  **Implement Feature**: Write the code in `internal/server/` to pass the test.
3.  **Verify**: Run `go test ./internal/server/...`. If it fails, fix and retry.

## 3. Handoff & Deliverables
Upon completion, you must provide:
1.  **Committed Code**: `internal/server/*.go`, `main.go`.
2.  **Committed Test**: `internal/server/server_test.go`.
3.  **Token Efficiency Report**:
    -   Estimated Input Tokens: [Value]
    -   Actual Output Tokens: [Value]
    -   Optimization Strategy: [e.g., "Only read necessary files"]
4.  **WebSocket Signal**: Send signal `BACKEND_READY` (Simulated).

## 4. Acceptance Criteria
- [ ] `go test ./internal/server/...` passes.
- [ ] `curl http://localhost:8080/api/health` returns `{"status": "ok"}`.
- [ ] Code is self-documenting (comments).
