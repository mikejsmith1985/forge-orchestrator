# Contract: Documentation (Issue 029)

**Goal**: Create comprehensive documentation for the Forge Orchestrator to ensure maintainability and ease of use for future developers.

## Scope
1.  **Update `README.md`**:
    -   **Project Overview**: What is Forge Orchestrator?
    -   **Features**: Architect, Ledger, Commands, Flows, Keyring.
    -   **Setup Guide**:
        -   Prerequisites (Go, Node.js).
        -   Installation steps.
        -   Running locally (`go run main.go` + `npm run dev`).
    -   **Architecture**: Brief high-level overview.

2.  **Create `docs/architecture.md`**:
    -   **System Design**: Diagram (Mermaid) of Frontend <-> BFF <-> SQLite <-> LLM Gateway.
    -   **Data Flow**: How a command travels from UI to LLM and back to Ledger.
    -   **Key Components**:
        -   `internal/server`: API Handlers.
        -   `internal/llm`: Gateway & Providers.
        -   `internal/flows`: Execution Engine.
        -   `internal/security`: Keyring.
    -   **Database Schema**: Description of tables (`token_ledger`, `command_cards`, `forge_flows`).

## Success Criteria
-   `README.md` is clear, typo-free, and allows a new user to start the app.
-   `docs/architecture.md` accurately reflects the implemented codebase.

## Handoff
-   **Signal File**: `handoffs/issue_029_docs.json`
-   **Git Branch**: `feat/issue-029-docs`
