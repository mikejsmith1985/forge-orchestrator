# Issue #032: Implement Real WebSocket Signaling Hub

**Priority:** 🔴 CRITICAL  
**Estimated Tokens:** ~2,500 (Medium complexity)  
**Agent Role:** Implementation

---

## 1. 🎫 Related Issue Context

**Gap Reference:** GAP-001 from v1-analysis.md

The current WebSocket implementation in `internal/server/websocket.go` only echoes messages back. Per the Project Charter: "Go Flow Engine uses WebSockets for primary signaling."

**Current Problem:**
```go
// For now, we just echo messages back
for {
    mt, message, err := conn.ReadMessage()
    err = conn.WriteMessage(mt, message)
}
```

---

## 2. 📋 Acceptance Criteria

### Backend (Go)
- [ ] Create `internal/server/hub.go` with `Hub` struct managing client connections
- [ ] Create `internal/server/client.go` with `Client` struct (connection, send channel)
- [ ] Hub must support: `Register()`, `Unregister()`, `Broadcast()`, `SendToClient()`
- [ ] Rewrite `websocket.go` to use Hub pattern
- [ ] Message types to support:
  - `FLOW_STATUS` (flowId, status: STARTED|NODE_COMPLETE|FINISHED|FAILED, nodeId, timestamp)
  - `LEDGER_UPDATE` (new entry notification)
  - `OPTIMIZATION_AVAILABLE` (new suggestion notification)

### Frontend (React)
- [ ] Create `frontend/src/hooks/useWebSocket.ts` custom hook
- [ ] Hook must: connect on mount, reconnect on disconnect (exponential backoff), parse messages
- [ ] Export `sendMessage()` function and `lastMessage` state

### Tests
- [ ] Unit test: Hub registers/unregisters clients correctly
- [ ] Unit test: Broadcast sends to all connected clients
- [ ] E2E test: Flow execution triggers real-time update in Ledger view

---

## 3. 📊 Token Efficiency Strategy

- Read only: `internal/server/websocket.go`, `internal/server/routes.go`
- Create new files rather than heavy refactoring
- Use standard gorilla/websocket patterns (already imported)

---

## 4. 🏗️ Technical Specification

### Hub Structure
```go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}
```

### Message Format
```json
{
    "type": "FLOW_STATUS",
    "payload": {
        "flowId": 1,
        "status": "NODE_COMPLETE",
        "nodeId": "agent-1",
        "timestamp": "2025-12-04T10:00:00Z"
    }
}
```

### Frontend Hook Interface
```typescript
const { isConnected, lastMessage, sendMessage } = useWebSocket('/ws');
```

---

## 5. 📁 Files to Create/Modify

| Action | File |
|--------|------|
| CREATE | `internal/server/hub.go` |
| CREATE | `internal/server/client.go` |
| MODIFY | `internal/server/websocket.go` |
| MODIFY | `internal/server/server.go` (add Hub to Server struct) |
| CREATE | `frontend/src/hooks/useWebSocket.ts` |
| MODIFY | `internal/server/server_test.go` (update WebSocket test) |

---

## 6. ✅ Definition of Done

1. WebSocket test passes with real message broadcasting
2. Console shows "Client connected" / "Client disconnected" logs
3. Frontend hook connects and logs received messages
4. No TypeScript or Go compilation errors
