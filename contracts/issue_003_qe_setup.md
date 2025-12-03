# Issue: Setup Playwright & Verify Layout (QE)

**Role**: Test Agent (QA)
**Priority**: High
**Context**: Frontend Layout has been implemented. We need to validate it using Playwright as per the Charter.

## 1. Requirements
1.  **Install Playwright**:
    - `npm init playwright@latest` (or manual install).
    - Configure for local testing (`http://localhost:8080`).
2.  **Create `tests/e2e/layout.spec.ts`**:
    - Test 1: Verify Sidebar is visible on desktop (`data-testid='sidebar'`).
    - Test 2: Verify Sidebar is hidden/collapsed on mobile viewport (`isMobile: true`).
    - Test 3: Verify "Forge Orchestrator" title/logo is present.

## 2. TDD & Verification Protocol
> [!IMPORTANT]
> You must follow this TDD workflow.
1.  **Run Tests**: `npx playwright test`.
2.  **Analyze Results**:
    -   If PASS: Proceed to Handoff.
    -   If FAIL: Report failure to Orchestrator (User) or fix test if it's a test issue.
3.  **Visual Validation**:
    -   Review Playwright traces/screenshots for visual correctness.

## 3. Handoff & Deliverables
Upon completion, you must provide:
1.  **Committed Test**: `tests/e2e/layout.spec.ts`.
2.  **Test Report**:
    -   Status: PASS/FAIL
    -   Screenshots/Traces path.
3.  **Token Efficiency Report**:
    -   Estimated Input Tokens: [Value]
    -   Actual Output Tokens: [Value]
    -   Optimization Strategy: [e.g., "Targeted specific spec file"]
4.  **WebSocket Signal**: Send signal `QE_VERIFIED` (Simulated).

## 4. Acceptance Criteria
- [ ] Playwright installs successfully.
- [ ] `npx playwright test` runs and passes.
- [ ] Tests use robust locators (data-testid, role).
