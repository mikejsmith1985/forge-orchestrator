# Contract: Final Cleanup & Linting (Issue 030)

**Goal**: Polish the codebase for release by removing debug artifacts, ensuring consistent formatting, and fixing any lingering lint issues.

## Scope
1.  **Code Cleanup**:
    -   Remove `// TODO` comments that have been addressed.
    -   Remove `fmt.Println` debug logs (use `log.Println` or structured logging if needed, but keep console clean).
    -   Ensure "Educational Comments" are preserved (do NOT remove them).

2.  **Linting & Formatting**:
    -   Run `gofmt -s -w .` on backend code.
    -   Run `npm run lint` (if configured) or ensure consistent indentation in frontend.
    -   Check for unused imports or variables.

3.  **Final Sanity Check**:
    -   Verify `go.mod` and `package.json` are clean.
    -   Ensure `.gitignore` is correct (excludes `dist/`, `node_modules/`, `token_ledger.db`).

## Success Criteria
-   Codebase compiles without warnings.
-   No stray debug prints in the console during normal operation.
-   Project structure is clean and organized.

## Handoff
-   **Signal File**: `handoffs/issue_030_cleanup.json`
-   **Git Branch**: `feat/issue-030-cleanup`
