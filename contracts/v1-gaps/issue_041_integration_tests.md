# Issue #041: Add Integration Test Suite

**Priority:** 🟢 LOW  
**Estimated Tokens:** ~1,800 (Medium complexity)  
**Agent Role:** Test

---

## 1. 🎫 Related Issue Context

**Gap Reference:** GAP-007 from v1-analysis.md

Current E2E tests mock all API responses:
```typescript
await page.route('/api/commands', async route => {
    await route.fulfill({ json: mockData }); // Fake data!
});
```

This doesn't test the actual frontend + backend integration. Per Charter: "The Test Agent focuses on whether the React UI renders correctly and handles user interaction."

---

## 2. 📋 Acceptance Criteria

### Test Infrastructure
- [ ] Create `tests/integration/` directory
- [ ] Create Go test that starts real server on random port
- [ ] Server uses in-memory SQLite database
- [ ] Create separate Playwright config for integration tests

### Critical Path Tests
- [ ] Test 1: Create command → Run command → Verify ledger entry appears
- [ ] Test 2: Create flow → Execute flow → Verify ledger entries created
- [ ] Test 3: Set API key → Verify status shows "Configured"
- [ ] Test 4: Add optimization data → Verify suggestions appear in UI

### CI Integration
- [ ] Add integration test job to GitHub Actions
- [ ] Run after unit tests pass
- [ ] Use mock LLM provider to avoid real API calls

---

## 3. 📊 Token Efficiency Strategy

- Reuse existing E2E test patterns
- Use Go test for server setup/teardown
- Minimal new code (mostly configuration)

---

## 4. 🏗️ Technical Specification

### Go Test Server Setup
```go
// tests/integration/setup_test.go
package integration

import (
    "database/sql"
    "net/http/httptest"
    "testing"
    
    "github.com/mikejsmith1985/forge-orchestrator/internal/server"
    "github.com/mikejsmith1985/forge-orchestrator/internal/data"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
    // Setup
    db, _ := sql.Open("sqlite3", ":memory:")
    db.Exec(data.SQLiteSchema)
    
    srv := server.NewServer(db)
    testServer = httptest.NewServer(srv.RegisterRoutes())
    
    // Run tests
    code := m.Run()
    
    // Teardown
    testServer.Close()
    db.Close()
    
    os.Exit(code)
}

func GetServerURL() string {
    return testServer.URL
}
```

### Playwright Integration Config
```typescript
// frontend/playwright.integration.config.ts
import { defineConfig } from '@playwright/test';

export default defineConfig({
    testDir: './tests/integration',
    use: {
        baseURL: process.env.INTEGRATION_SERVER_URL || 'http://localhost:8080',
    },
    webServer: undefined, // Don't start dev server, use Go test server
});
```

### Integration Test Example
```typescript
// frontend/tests/integration/command-flow.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Command Execution Flow', () => {
    test('complete command lifecycle', async ({ page }) => {
        // 1. Navigate to Commands
        await page.goto('/commands');
        
        // 2. Create a command (NO MOCKING!)
        await page.getByTestId('add-command-btn').click();
        await page.getByTestId('command-name-input').fill('Integration Test Cmd');
        await page.getByTestId('command-input').fill('echo hello');
        await page.getByTestId('submit-command-btn').click();
        
        // 3. Verify it appears in the list
        await expect(page.getByText('Integration Test Cmd')).toBeVisible();
        
        // 4. Navigate to Ledger (initially empty)
        await page.goto('/ledger');
        
        // Note: Running command requires API key, which we'd need to set up
        // For CI, we'd use a mock provider in the backend
    });
});
```

### Mock LLM Provider for CI
```go
// internal/llm/mock_provider.go
type MockProvider struct{}

func (m *MockProvider) Send(system, user, key string) (string, int, int, error) {
    return `{"result": "mock response"}`, 50, 100, nil
}
```

---

## 5. 📁 Files to Create/Modify

| Action | File |
|--------|------|
| CREATE | `tests/integration/setup_test.go` |
| CREATE | `tests/integration/command_flow_test.go` |
| CREATE | `frontend/playwright.integration.config.ts` |
| CREATE | `frontend/tests/integration/command-flow.spec.ts` |
| CREATE | `frontend/tests/integration/flow-execution.spec.ts` |
| CREATE | `frontend/tests/integration/key-management.spec.ts` |
| MODIFY | `.github/workflows/release.yml` (add integration test job) |

---

## 6. ✅ Definition of Done

1. `go test ./tests/integration/...` passes
2. `npx playwright test --config=playwright.integration.config.ts` passes
3. Tests use real database (in-memory SQLite)
4. Tests use real API endpoints (no mocking)
5. CI pipeline runs integration tests after unit tests
6. Documentation explains how to run integration tests locally
