# Contract: Production Build & Stability (Issue 031)

**Goal**: Transition the application from a fragile "Dev Mode" (Vite + Go Proxy) to a robust "Single Binary" architecture. This resolves the stability issues where the frontend server crashes or disconnects.

## Scope
1.  **Frontend Build**:
    -   Execute `npm run build` in `frontend/` to generate the static assets in `dist/`.
    -   Ensure the build is error-free.

2.  **Backend SPA Handler (`main.go`)**:
    -   Update `main.go` to correctly serve the Single Page Application (SPA).
    -   **Problem**: `http.FileServer` returns 404 for client-side routes (e.g., `/architect`, `/ledger`) on refresh.
    -   **Fix**: Implement a custom handler that serves `index.html` for any path that doesn't match a static file (and isn't an API route).

3.  **Verification**:
    -   Stop all running dev servers (`vite`, `go run`).
    -   Run `go run main.go` (or `go build` + `./forge-orchestrator`).
    -   Verify the app loads at `http://localhost:8080`.
    -   Verify navigation to `/architect` and refreshing the page works (no 404).

## Success Criteria
-   Application runs as a single process on port 8080.
-   No dependency on `npm run dev` or `vite` for runtime.
-   Client-side routing works (refreshing sub-pages loads the app).

## Handoff
-   **Signal File**: `handoffs/issue_031_production_build.json`
-   **Git Branch**: `feat/issue-031-production-build`
