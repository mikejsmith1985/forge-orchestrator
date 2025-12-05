# Issue #033: Add Fallback Signaling for WebSocket Failures

**Priority:** 🔴 CRITICAL  
**Estimated Tokens:** ~1,800 (Low-Medium complexity)  
**Agent Role:** Implementation  
**Depends On:** Issue #032 (WebSocket Hub)

---

## 1. 🎫 Related Issue Context

**Gap Reference:** GAP-002 from v1-analysis.md

Per Project Charter: "implement a GitHub Commit Check as a fallback if the signal fails."

When WebSocket is unavailable, the system needs an alternative way to check flow execution status.

---

## 2. 📋 Acceptance Criteria

### Backend (Go)
- [ ] Create `internal/flows/signaler.go` with `Signaler` interface
- [ ] Implement `WebSocketSignaler` (uses Hub from #032)
- [ ] Implement `FileSignaler` (writes to `.forge/status/{flowId}.json`)
- [ ] Flow engine: Try WebSocket first, fall back to file after 5 second timeout
- [ ] Add `GET /api/flows/{id}/status` endpoint for polling

### File-Based Fallback
- [ ] Create `.forge/status/` directory on startup if not exists
- [ ] Write status file on each state change: `{"flowId": 1, "status": "RUNNING", "lastNode": "agent-1", "updatedAt": "..."}`
- [ ] Frontend polls this endpoint every 3 seconds when WebSocket disconnected

### Frontend (React)
- [ ] Update `useWebSocket` hook to detect disconnection
- [ ] When disconnected > 5 seconds, start polling `/api/flows/{id}/status`
- [ ] Stop polling when WebSocket reconnects

### Tests
- [ ] Unit test: FileSignaler writes correct JSON format
- [ ] Unit test: Status endpoint returns current flow state
- [ ] Integration test: Simulate WebSocket failure, verify polling works

---

## 3. 📊 Token Efficiency Strategy

- Small, focused files (signaler.go ~80 lines)
- Reuse existing patterns from flows/engine.go
- Minimal frontend changes (extend existing hook)

---

## 4. 🏗️ Technical Specification

### Signaler Interface
```go
type Signaler interface {
    NotifyStatus(flowID int, status FlowStatus) error
    GetStatus(flowID int) (*FlowStatus, error)
}

type FlowStatus struct {
    FlowID    int       `json:"flowId"`
    Status    string    `json:"status"` // PENDING, RUNNING, COMPLETED, FAILED
    LastNode  string    `json:"lastNode,omitempty"`
    UpdatedAt time.Time `json:"updatedAt"`
    Error     string    `json:"error,omitempty"`
}
```

### Fallback Logic in Engine
```go
func (e *Engine) notifyStatus(flowID int, status FlowStatus) {
    err := e.wsSignaler.NotifyStatus(flowID, status)
    if err != nil {
        log.Printf("WebSocket notify failed, using file fallback: %v", err)
        e.fileSignaler.NotifyStatus(flowID, status)
    }
}
```

---

## 5. 📁 Files to Create/Modify

| Action | File |
|--------|------|
| CREATE | `internal/flows/signaler.go` |
| CREATE | `internal/flows/file_signaler.go` |
| CREATE | `internal/flows/ws_signaler.go` |
| MODIFY | `internal/flows/engine.go` (add signaler calls) |
| MODIFY | `internal/server/routes.go` (add status endpoint) |
| CREATE | `internal/server/flow_status.go` (handler) |
| MODIFY | `frontend/src/hooks/useWebSocket.ts` (add polling fallback) |

---

## 6. ✅ Definition of Done

1. Flow execution creates status file in `.forge/status/`
2. `GET /api/flows/1/status` returns valid JSON
3. Frontend switches to polling when WebSocket disconnects
4. All tests pass
