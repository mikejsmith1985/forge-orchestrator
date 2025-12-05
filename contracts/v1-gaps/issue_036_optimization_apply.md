# Issue #036: Implement Real Optimization Apply Logic

**Priority:** 🟡 HIGH  
**Estimated Tokens:** ~2,200 (Medium complexity)  
**Agent Role:** Implementation

---

## 1. 🎫 Related Issue Context

**Gap Reference:** GAP-011 from v1-analysis.md

The "Apply" button for optimizations only logs to ledger but doesn't actually modify configurations:

```go
// In a real implementation, we would fetch the suggestion by ID and apply it.
// Here, we just log that it was applied.
```

Per Charter: "Token Ledger Audit Flow Playback: Displays suggestions with a 'Click to Apply' button that instantly updates the associated Flow Node/Command Card."

---

## 2. 📋 Acceptance Criteria

### Backend (Go)
- [ ] Create `internal/optimizer/applier.go` with apply logic
- [ ] Parse `ApplyAction` JSON to determine action type
- [ ] Implement `model_switch`: Update provider in flow's node data
- [ ] Implement `prompt_optimization`: Flag flow for review (add `needs_optimization: true`)
- [ ] Implement `retry_strategy`: Add retry config to flow metadata
- [ ] Return updated configuration in response

### Suggestion Persistence
- [ ] Create `optimization_suggestions` table to store generated suggestions
- [ ] Store suggestions with IDs so they can be fetched by ID when applying
- [ ] Mark suggestions as "applied" after successful application

### Frontend (React)
- [ ] Update OptimizationCard to show what will change before applying
- [ ] Add confirmation modal: "This will change [flow X] from [gpt-4] to [gpt-3.5-turbo]. Continue?"
- [ ] Show success message with summary of changes made

### Tests
- [ ] Unit test: model_switch updates flow data correctly
- [ ] Unit test: Applied suggestion marked in database
- [ ] E2E test: Apply optimization → Flow configuration changes

---

## 3. 📊 Token Efficiency Strategy

- Add new table via migration in sqlite_schema.go
- Single new file for apply logic
- Minimal frontend changes (add modal)

---

## 4. 🏗️ Technical Specification

### New Database Table
```sql
CREATE TABLE IF NOT EXISTS optimization_suggestions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    estimated_savings REAL NOT NULL,
    savings_unit TEXT NOT NULL,
    target_flow_id TEXT,
    target_command_id INTEGER,
    apply_action TEXT NOT NULL,
    status TEXT DEFAULT 'pending',
    applied_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Apply Action Parsing
```go
type ApplyAction struct {
    Action    string `json:"action"` // model_switch, prompt_optimization, retry_strategy
    FlowID    string `json:"flow_id,omitempty"`
    FromModel string `json:"from_model,omitempty"`
    ToModel   string `json:"to_model,omitempty"`
}

func ApplyOptimization(db *sql.DB, suggestionID int) (*ApplyResult, error) {
    // 1. Fetch suggestion by ID
    // 2. Parse ApplyAction JSON
    // 3. Execute appropriate action
    // 4. Mark suggestion as applied
    // 5. Return result
}
```

### Model Switch Implementation
```go
func applyModelSwitch(db *sql.DB, action ApplyAction) error {
    // 1. Fetch flow data
    var flowData string
    db.QueryRow("SELECT data FROM forge_flows WHERE id = ?", action.FlowID).Scan(&flowData)
    
    // 2. Parse and update nodes
    var graph FlowGraph
    json.Unmarshal([]byte(flowData), &graph)
    
    for i, node := range graph.Nodes {
        if node.Data.Provider == action.FromModel {
            graph.Nodes[i].Data.Provider = action.ToModel
        }
    }
    
    // 3. Save updated flow
    updatedData, _ := json.Marshal(graph)
    db.Exec("UPDATE forge_flows SET data = ? WHERE id = ?", string(updatedData), action.FlowID)
    
    return nil
}
```

---

## 5. 📁 Files to Create/Modify

| Action | File |
|--------|------|
| CREATE | `internal/optimizer/applier.go` |
| MODIFY | `internal/data/sqlite_schema.go` (add table) |
| MODIFY | `internal/optimizer/analyzer.go` (store suggestions) |
| MODIFY | `internal/server/optimizer.go` (implement apply) |
| MODIFY | `frontend/src/components/Ledger/OptimizationCard.tsx` (add modal) |
| CREATE | `internal/optimizer/applier_test.go` |

---

## 6. ✅ Definition of Done

1. Suggestions are stored in database with unique IDs
2. `POST /api/ledger/optimizations/{id}/apply` updates flow configuration
3. Frontend shows confirmation before applying
4. Applied suggestions show as "Applied" and are disabled
5. Flow data in database reflects the applied change
6. All tests pass
