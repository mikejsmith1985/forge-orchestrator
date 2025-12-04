# Contract: Integration Verification (Smoke Tests)

**Goal**: Verify the running application's integrity by executing real HTTP requests against the backend, ensuring the "Integration Gap" is closed.

## Scope
1.  **Smoke Test Script**: Create `scripts/smoke_test.go` to:
    -   Ping `GET /api/health`.
    -   Create a Ledger Entry (`POST /api/ledger`).
    -   Fetch Ledger Entries (`GET /api/ledger`).
    -   Create a Command (`POST /api/commands`).
    -   Execute a Command (`POST /api/commands/{id}/run`).
    -   Create a Flow (`POST /api/flows`).
2.  **Execution**: Run this script against `http://localhost:8080` (Backend) and `http://localhost:8081` (Frontend Proxy) to verify both direct and proxied access.

## Success Criteria
-   All requests return `200 OK` (or `201 Created`).
-   Responses are valid JSON.
-   **CRITICAL**: No `Unexpected token <` errors (HTML fallbacks).

## Handoff
-   **Signal File**: `handoffs/task_integration_verification.json`
