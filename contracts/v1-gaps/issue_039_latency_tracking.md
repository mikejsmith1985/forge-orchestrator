# Issue #039: Add Command Execution Latency Tracking

**Priority:** 🟢 MEDIUM  
**Estimated Tokens:** ~1,000 (Low complexity)  
**Agent Role:** Implementation

---

## 1. 🎫 Related Issue Context

**Gap Reference:** GAP-013 from v1-analysis.md

Command execution doesn't record latency:
```go
// Note: Latency is not captured here yet, could add timing around ExecutePrompt
```

The ledger shows `0` for latency on command executions, making performance analysis incomplete.

---

## 2. 📋 Acceptance Criteria

### Backend (Go)
- [ ] Add `time.Now()` before `ExecutePrompt()` call
- [ ] Calculate `time.Since()` after call completes
- [ ] Store latency in milliseconds in ledger entry
- [ ] Also add latency tracking to flow execution

### Frontend (React)
- [ ] Add "Latency" column to Ledger table
- [ ] Format as "XXX ms" or "X.XX s" for longer durations
- [ ] Color code: Green < 1s, Yellow 1-5s, Red > 5s

### Tests
- [ ] Unit test: Ledger entry contains non-zero latency after command execution
- [ ] E2E test: Latency column visible and shows reasonable values

---

## 3. 📊 Token Efficiency Strategy

- Very small change to commands.go
- Add one column to frontend table
- Reuse existing patterns

---

## 4. 🏗️ Technical Specification

### Backend Timing
```go
func (s *Server) handleRunCommand(w http.ResponseWriter, r *http.Request) {
    // ... existing setup code ...
    
    // Time the execution
    startTime := time.Now()
    response, err := s.gateway.ExecutePrompt(req.AgentRole, commandPrompt, apiKey, provider)
    latencyMs := time.Since(startTime).Milliseconds()
    
    ledgerEntry := LedgerEntry{
        // ... existing fields ...
        LatencyMS: int(latencyMs), // Add this
    }
    
    // ... rest of handler ...
}
```

### Frontend Column
```tsx
<th className="px-6 py-4 text-right">Latency</th>

// In table body:
<td className="px-6 py-4 text-right font-mono">
    <span className={
        entry.latency_ms < 1000 ? 'text-green-400' :
        entry.latency_ms < 5000 ? 'text-yellow-400' : 'text-red-400'
    }>
        {entry.latency_ms < 1000 
            ? `${entry.latency_ms}ms` 
            : `${(entry.latency_ms / 1000).toFixed(2)}s`}
    </span>
</td>
```

---

## 5. 📁 Files to Create/Modify

| Action | File |
|--------|------|
| MODIFY | `internal/server/commands.go` (add timing) |
| MODIFY | `internal/flows/engine.go` (add timing to flow nodes) |
| MODIFY | `frontend/src/components/Ledger/LedgerView.tsx` (add column) |

---

## 6. ✅ Definition of Done

1. Command execution records accurate latency in ledger
2. Flow node execution records accurate latency in ledger
3. Ledger UI shows latency column with color coding
4. Previously recorded entries with 0 latency still display correctly
5. All tests pass
