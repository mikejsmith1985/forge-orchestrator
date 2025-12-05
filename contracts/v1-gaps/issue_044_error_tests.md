# Issue #044: Add Error Handling Tests

**Priority:** 🟢 MEDIUM  
**Estimated Tokens:** ~1,200 (Low complexity)  
**Agent Role:** Test

---

## 1. 🎫 Related Issue Context

**Gap Reference:** GAP-006 from v1-analysis.md

Current tests only cover "happy path" scenarios. Missing error handling tests:

| File | Missing Tests |
|------|---------------|
| `commands_test.go` | Malformed JSON body |
| `commands_test.go` | SQL injection attempts |
| `ledger_test.go` | Invalid timestamp format |
| `flows/engine_test.go` | Flow with no nodes |
| `flows/engine_test.go` | Flow with invalid provider |

---

## 2. 📋 Acceptance Criteria

### Commands Tests
- [ ] Test: Malformed JSON returns 400 Bad Request
- [ ] Test: Empty name field returns 400 Bad Request
- [ ] Test: Delete non-existent ID returns appropriate status
- [ ] Test: SQL injection in name field doesn't break database

### Ledger Tests
- [ ] Test: Invalid JSON body returns 400
- [ ] Test: Missing required fields returns 400
- [ ] Test: Non-numeric limit parameter is handled gracefully

### Flows Engine Tests
- [ ] Test: Flow with empty nodes array completes without error
- [ ] Test: Flow with invalid provider returns meaningful error
- [ ] Test: Flow with missing API key returns clear error message
- [ ] Test: Database connection failure is handled

### Gateway Tests
- [ ] Test: Invalid API key returns error (not panic)
- [ ] Test: Network timeout is handled
- [ ] Test: Malformed response from LLM is handled

---

## 3. 📊 Token Efficiency Strategy

- Add tests to existing test files
- Use table-driven tests for efficiency
- ~20-30 new test cases total

---

## 4. 🏗️ Technical Specification

### Commands Error Tests
```go
// internal/server/commands_test.go

func TestHandleCreateCommand_Errors(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    server := NewServer(db)
    handler := server.RegisterRoutes()
    
    tests := []struct {
        name       string
        body       string
        wantStatus int
    }{
        {
            name:       "malformed JSON",
            body:       `{"name": "test"`, // missing closing brace
            wantStatus: http.StatusBadRequest,
        },
        {
            name:       "empty name",
            body:       `{"name": "", "command": "echo hello"}`,
            wantStatus: http.StatusBadRequest,
        },
        {
            name:       "empty command",
            body:       `{"name": "Test", "command": ""}`,
            wantStatus: http.StatusBadRequest,
        },
        {
            name:       "SQL injection attempt",
            body:       `{"name": "'; DROP TABLE command_cards; --", "command": "test"}`,
            wantStatus: http.StatusCreated, // Should succeed but sanitized
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req, _ := http.NewRequest("POST", "/api/commands", 
                bytes.NewBufferString(tt.body))
            rr := httptest.NewRecorder()
            handler.ServeHTTP(rr, req)
            
            if rr.Code != tt.wantStatus {
                t.Errorf("got status %d, want %d", rr.Code, tt.wantStatus)
            }
        })
    }
}

func TestHandleDeleteCommand_NotFound(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    server := NewServer(db)
    handler := server.RegisterRoutes()
    
    req, _ := http.NewRequest("DELETE", "/api/commands/99999", nil)
    rr := httptest.NewRecorder()
    handler.ServeHTTP(rr, req)
    
    // Should succeed even if not found (idempotent delete)
    // OR return 404 - document expected behavior
    if rr.Code != http.StatusNoContent && rr.Code != http.StatusNotFound {
        t.Errorf("unexpected status: %d", rr.Code)
    }
}
```

### Flows Engine Error Tests
```go
// internal/flows/engine_test.go

func TestExecuteFlow_EmptyNodes(t *testing.T) {
    // Setup...
    flowJSON := `{"nodes": [], "edges": []}`
    _, err := db.Exec(`INSERT INTO forge_flows (name, data, status) VALUES (?, ?, ?)`, 
        "Empty Flow", flowJSON, "active")
    
    err = ExecuteFlow(1, db, gateway)
    
    // Should complete without error (nothing to do)
    if err != nil {
        t.Errorf("unexpected error for empty flow: %v", err)
    }
}

func TestExecuteFlow_InvalidProvider(t *testing.T) {
    // Setup...
    flowJSON := `{
        "nodes": [{
            "id": "1",
            "type": "agent",
            "data": {"role": "coder", "provider": "InvalidProvider", "prompt": "test"}
        }],
        "edges": []
    }`
    
    err := ExecuteFlow(1, db, gateway)
    
    if err == nil {
        t.Error("expected error for invalid provider")
    }
    if !strings.Contains(err.Error(), "unsupported provider") {
        t.Errorf("error should mention unsupported provider: %v", err)
    }
}
```

---

## 5. 📁 Files to Create/Modify

| Action | File |
|--------|------|
| MODIFY | `internal/server/commands_test.go` |
| MODIFY | `internal/server/ledger_test.go` |
| MODIFY | `internal/flows/engine_test.go` |
| MODIFY | `internal/llm/gateway_test.go` |

---

## 6. ✅ Definition of Done

1. All error handling tests pass
2. Each error scenario returns appropriate HTTP status code
3. Error messages are user-friendly (not stack traces)
4. SQL injection attempts don't corrupt database
5. Network/timeout errors are handled gracefully
6. Code coverage increased for error paths
