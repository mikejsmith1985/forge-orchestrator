# Issue #045: WebSocket Hub Integration with Flow Engine

**Priority:** 🟡 HIGH  
**Estimated Tokens:** ~1,500 (Medium complexity)  
**Agent Role:** Implementation  
**Depends On:** Issue #032 (WebSocket Hub)

---

## 1. 🎫 Related Issue Context

Issue #032 creates the WebSocket Hub infrastructure. This issue integrates it with the Flow Engine to broadcast real-time status updates during flow execution.

---

## 2. 📋 Acceptance Criteria

### Backend (Go)
- [ ] Inject Hub reference into Flow Engine
- [ ] Broadcast `FLOW_STARTED` when flow execution begins
- [ ] Broadcast `NODE_STARTED` before each node executes
- [ ] Broadcast `NODE_COMPLETED` after each node executes (with result preview)
- [ ] Broadcast `FLOW_COMPLETED` when flow finishes successfully
- [ ] Broadcast `FLOW_FAILED` when flow fails (with error message)

### Message Payloads
- [ ] Include flowId, nodeId, timestamp in all messages
- [ ] Include token counts in NODE_COMPLETED
- [ ] Include error message in FLOW_FAILED
- [ ] Include total execution time in FLOW_COMPLETED

### Frontend (React)
- [ ] LedgerView subscribes to WebSocket updates
- [ ] Auto-refresh ledger table when FLOW_COMPLETED received
- [ ] Show toast notification for FLOW_STARTED
- [ ] Show toast notification for FLOW_FAILED

### Tests
- [ ] Unit test: Flow execution triggers correct sequence of broadcasts
- [ ] E2E test: Run flow → Verify ledger updates in real-time

---

## 3. 📊 Token Efficiency Strategy

- Small modifications to existing engine.go
- Reuse Hub from #032
- Add notification helper function

---

## 4. 🏗️ Technical Specification

### Engine with Hub
```go
// internal/flows/engine.go

type Engine struct {
    db      *sql.DB
    gateway *llm.Gateway
    hub     *server.Hub // Injected
}

func NewEngine(db *sql.DB, gateway *llm.Gateway, hub *server.Hub) *Engine {
    return &Engine{db: db, gateway: gateway, hub: hub}
}

func (e *Engine) ExecuteFlow(flowID int) error {
    // Broadcast start
    e.broadcast(FlowMessage{
        Type: "FLOW_STARTED",
        Payload: map[string]interface{}{
            "flowId":    flowID,
            "timestamp": time.Now(),
        },
    })
    
    // ... existing flow logic ...
    
    for _, node := range graph.Nodes {
        // Broadcast node start
        e.broadcast(FlowMessage{
            Type: "NODE_STARTED",
            Payload: map[string]interface{}{
                "flowId":    flowID,
                "nodeId":    node.ID,
                "label":     node.Data.Label,
                "timestamp": time.Now(),
            },
        })
        
        // Execute node...
        resp, err := e.gateway.ExecutePrompt(...)
        
        // Broadcast node complete
        e.broadcast(FlowMessage{
            Type: "NODE_COMPLETED",
            Payload: map[string]interface{}{
                "flowId":       flowID,
                "nodeId":       node.ID,
                "inputTokens":  resp.InputTokens,
                "outputTokens": resp.OutputTokens,
                "cost":         resp.Cost,
                "timestamp":    time.Now(),
            },
        })
    }
    
    // Broadcast complete
    e.broadcast(FlowMessage{
        Type: "FLOW_COMPLETED",
        Payload: map[string]interface{}{
            "flowId":    flowID,
            "timestamp": time.Now(),
        },
    })
    
    return nil
}

func (e *Engine) broadcast(msg FlowMessage) {
    if e.hub == nil {
        return // No hub configured (testing)
    }
    data, _ := json.Marshal(msg)
    e.hub.Broadcast(data)
}
```

### Frontend Integration
```typescript
// In LedgerView.tsx
const { lastMessage } = useWebSocket('/ws');

useEffect(() => {
    if (!lastMessage) return;
    
    const msg = JSON.parse(lastMessage);
    
    switch (msg.type) {
        case 'FLOW_STARTED':
            toast.info(`Flow ${msg.payload.flowId} started`);
            break;
        case 'FLOW_COMPLETED':
            toast.success(`Flow ${msg.payload.flowId} completed`);
            fetchData(); // Refresh ledger
            break;
        case 'FLOW_FAILED':
            toast.error(`Flow ${msg.payload.flowId} failed: ${msg.payload.error}`);
            fetchData();
            break;
    }
}, [lastMessage]);
```

---

## 5. 📁 Files to Create/Modify

| Action | File |
|--------|------|
| MODIFY | `internal/flows/engine.go` (add Hub, broadcasts) |
| MODIFY | `internal/server/flows.go` (pass Hub to engine) |
| MODIFY | `internal/server/server.go` (store Hub reference) |
| MODIFY | `frontend/src/components/Ledger/LedgerView.tsx` |
| CREATE | `internal/flows/messages.go` (message types) |

---

## 6. ✅ Definition of Done

1. Flow execution broadcasts all lifecycle events
2. Frontend receives and processes WebSocket messages
3. Ledger auto-refreshes when flow completes
4. Toast notifications appear for flow start/complete/fail
5. Tests verify broadcast sequence
