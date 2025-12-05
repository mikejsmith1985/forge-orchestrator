# Issue #038: Fix Agent Role Name Mapping

**Priority:** 🟡 MEDIUM  
**Estimated Tokens:** ~800 (Low complexity)  
**Agent Role:** Implementation

---

## 1. 🎫 Related Issue Context

**Gap Reference:** GAP-014 from v1-analysis.md

Agent prompts use formal names but flow data uses informal names:

**Prompts define:**
- `"Architect"`, `"Implementation"`, `"Test"`, `"Optimizer"`

**Flow JSON uses:**
- `"planner"`, `"coder"`, `"tester"`, `"auditor"`

This causes `GetAgentPrompt()` to return errors for valid roles.

---

## 2. 📋 Acceptance Criteria

### Backend (Go)
- [ ] Create role alias map in `agent_prompts.go`
- [ ] Update `GetAgentPrompt()` to resolve aliases before lookup
- [ ] Support both formal and informal names
- [ ] Aliases are case-insensitive

### Mapping
| Alias | Canonical Name |
|-------|----------------|
| planner | Architect |
| coder | Implementation |
| developer | Implementation |
| dev | Implementation |
| tester | Test |
| qa | Test |
| auditor | Optimizer |
| optimizer | Optimizer |

### Frontend (React)
- [ ] Update FlowEditor node configuration to use consistent role names
- [ ] Add role dropdown with user-friendly labels

### Tests
- [ ] Unit test: `GetAgentPrompt("coder")` returns Implementation prompt
- [ ] Unit test: `GetAgentPrompt("PLANNER")` returns Architect prompt (case insensitive)
- [ ] Unit test: Unknown role still returns error

---

## 3. 📊 Token Efficiency Strategy

- Changes only in agent_prompts.go (~20 new lines)
- Simple map-based lookup
- Minimal frontend change (dropdown options)

---

## 4. 🏗️ Technical Specification

### Role Alias Map
```go
var roleAliases = map[string]string{
    "planner":   "Architect",
    "coder":     "Implementation",
    "developer": "Implementation",
    "dev":       "Implementation",
    "tester":    "Test",
    "qa":        "Test",
    "auditor":   "Optimizer",
    "optimizer": "Optimizer",
}

func resolveRole(role string) string {
    normalized := strings.ToLower(strings.TrimSpace(role))
    
    // Check if it's an alias
    if canonical, ok := roleAliases[normalized]; ok {
        return canonical
    }
    
    // Check if it's already a canonical name (case-insensitive)
    canonicalNames := []string{"Architect", "Implementation", "Test", "Optimizer"}
    for _, name := range canonicalNames {
        if strings.EqualFold(role, name) {
            return name
        }
    }
    
    return role // Return as-is, GetAgentPrompt will error
}

func GetAgentPrompt(role string) (string, error) {
    resolvedRole := resolveRole(role)
    
    switch resolvedRole {
    case "Architect":
        return SystemPromptArchitect, nil
    // ... rest of cases
    }
}
```

### Frontend Role Dropdown
```typescript
const agentRoles = [
    { value: 'Architect', label: 'Planner / Architect' },
    { value: 'Implementation', label: 'Developer / Coder' },
    { value: 'Test', label: 'QA / Tester' },
    { value: 'Optimizer', label: 'Auditor / Optimizer' },
];
```

---

## 5. 📁 Files to Create/Modify

| Action | File |
|--------|------|
| MODIFY | `internal/agents/agent_prompts.go` |
| MODIFY | `frontend/src/components/Flows/FlowEditor.tsx` (role dropdown) |
| MODIFY | `internal/agents/agent_prompts_test.go` (new file or extend) |

---

## 6. ✅ Definition of Done

1. `GetAgentPrompt("coder")` returns the Implementation prompt
2. All aliases work case-insensitively
3. FlowEditor shows role dropdown with friendly labels
4. Existing flows with informal role names continue to work
5. All tests pass
