# **Forge Orchestrator: Foundational Implementation Contract (V2.1)**

**Goal:** Implement the core architectural foundation (Interfaces, Database, and Execution Logic) required for the Forge Orchestrator. This contract is designed to be worked through sequentially by a **Single Implementation Agent** (e.g., using Copilot CLI) with atomic commits.

**Project Status:** Execution is assumed to be local and serial. All required files must be created before their usage in subsequent steps.

## **I. Global TDD and Code Quality Rules (MANDATORY)**

* **Validation Contract (The No Mocks Rule):** All tests must be executed via Go unit tests at this stage. **MOCKING OF CORE LOGIC IS FORBIDDEN.** Tests must verify real function execution (database writes, command execution).  
* **Self-Documenting Code:** All code MUST include verbose, plain-English comments, understandable to a child.  
* **Commit Granularity:** Commits MUST be atomic (one conceptual change per commit).

## **II. Contract 1: Database Setup and Core Data Model**

**Focus:** Data Persistence Layer.

| Field | Value |
| :---- | :---- |
| **Branch Name** | feat/init-database-layer |
| **Files Created** | internal/data/sqlite\_schema.go, internal/data/db.go |
| **Task Description** | 1\. Create sqlite\_schema.go defining the token\_ledger, forge\_flows, and user\_secrets table schemas using standard SQLite SQL. 2\. Create db.go with functions to connect to and initialize the SQLite database using the schema. The initialization function must create the database file if it does not exist. |
| **Validation Contract** | **Commit:** Feature code. **Test:** Go unit test verifying connection can be established and all three tables (token\_ledger, forge\_flows, user\_secrets) exist after initialization. |

## **III. Contract 2: Execution Interface and Local Runner**

**Focus:** Abstracting Shell Commands.

| Field | Value |
| :---- | :---- |
| **Branch Name** | feat/execution-interface-and-local-runner |
| **Files Created** | internal/execution/executor.go, internal/execution/local\_runner.go, internal/execution/execution\_models.go |
| **Task Description** | 1\. Define the Executor interface and necessary data models (ExecutionContext, ExecutionResult). 2\. Implement the Executor interface in local\_runner.go. This implementation MUST use Go's os/exec to execute shell commands locally, capturing all stdout, stderr, and the numeric exit code. |
| **Validation Contract** | **Commit:** Feature code. **Test:** Go unit test verifying that: a) local\_runner.Execute("echo hello world") returns "hello world" in stdout and exit code 0\. b) local\_runner.Execute("exit 1") returns exit code 1 and captures any stderr. |

## **IV. Contract 3: LLM Gateway Interface and Stub Implementation**

**Focus:** Abstracting AI Vendor APIs.

| Field | Value |
| :---- | :---- |
| **Branch Name** | feat/llm-gateway-stub-implementation |
| **Files Created** | internal/llm/gateway.go, internal/llm/stub\_adapter.go, internal/llm/llm\_models.go |
| **Task Description** | 1\. Define the LLMGateway interface and necessary models (LLMConfig, LLMResponse, BudgetStatus), including the PrimaryCostUnit field ('TOKEN' or 'PROMPT') in the config. 2\. Implement a non-functional StubAdapter that satisfies the LLMGateway interface. The Generate method MUST return hardcoded, successful JSON data and set token/cost to zero. **This stub prevents front-end API errors.** |
| **Validation Contract** | **Commit:** Feature code. **Test:** Go unit test confirming that the stub\_adapter.Generate() method returns the expected zero-cost, hardcoded response and that the PrimaryCostUnit field can be read successfully from the configuration struct. |

## **V. Contract 4: Token Ledger Insertion Logic**

**Focus:** Core Business Logic.

| Field | Value |
| :---- | :---- |
| **Branch Name** | feat/token-ledger-insertion-logic |
| **Files Required** | internal/data/db.go, internal/llm/gateway.go |
| **Files Created** | internal/data/ledger\_service.go, internal/data/data\_models.go |
| **Task Description** | 1\. Define the necessary data models (TokenLedgerEntry). 2\. Create ledger\_service.go with a function LogUsage(entry TokenLedgerEntry) error. This function MUST use the database connection from db.go to securely insert a complete record into the token\_ledger table, ensuring all required fields are correctly mapped. |
| **Validation Contract** | **Commit:** Feature code. **Test:** Go unit test that inserts a complete, test TokenLedgerEntry into the SQLite database and verifies that the record can be retrieved and its content matches the insertion data. |

## **VI. Contract 5: Initial Server and Frontend Wiring**

**Focus:** Final Go Scaffolding and FE/BE Communication Test.

| Field | Value |
| :---- | :---- |
| **Branch Name** | feat/initial-server-api-wiring |
| **Files Required** | cmd/forge/main.go, UI assets (assumed built), internal/execution/executor.go |
| **Files Created** | internal/server/router.go, internal/server/websocket.go, internal/server/api\_handlers.go |
| **Task Description** | 1\. Refactor main.go to use router.go and start the HTTP server. 2\. Implement a non-functional PTY/WebSocket handler in websocket.go (stubbed). 3\. Create a basic /api/execute handler in api\_handlers.go that calls the stub Executor interface. |
| **Validation Contract** | **Commit:** Feature code and Playwright test. **Test:** Playwright test **MUST** validate the full stack (NO MOCKS): React UI button click $\\rightarrow$ Go BFF endpoint reached $\\rightarrow$ Go endpoint successfully calls the stub Executor interface $\\rightarrow$ UI receives and renders "Execution Request Received" message. **(Crucial: Verifies FE talks to BFF.)** |

This consolidated plan provides the exact contracts needed. You can now tell your Implementation Agent to read this document and proceed with the tasks sequentially.