# Forge Orchestrator v1.0.0 - Gap Analysis & Remediation Plan

**Analysis Date:** December 4, 2025  
**Analyst:** AI Engineering Review  
**Document Version:** 1.0

---

## Table of Contents

1. [What is This Document?](#what-is-this-document)
2. [Executive Summary](#executive-summary)
3. [Requirements Recap (From Project Charter)](#requirements-recap-from-project-charter)
4. [Gap Analysis](#gap-analysis)
   - [Critical Gaps](#critical-gaps)
   - [Security Gaps](#security-gaps)
   - [Testing Gaps](#testing-gaps)
   - [Code Quality Issues](#code-quality-issues)
5. [GitHub Issue Contracts](#github-issue-contracts)
6. [Appendix: Judgment Calls](#appendix-judgment-calls)

---

## What is This Document?

Think of this document like a home inspection report before you buy a house. Even if the house looks great from the outside, an inspector checks for problems in the plumbing, electrical, foundation, and more. This document does the same thing for our software.

**The Forge Orchestrator** is a tool that helps developers work with AI coding assistants (like Claude or GPT-4). It tracks how much those AI calls cost (in "tokens" - the currency AI models use), helps organize workflows (called "Flows"), and provides a nice visual interface to manage everything.

We built version 1.0.0, but before we call it "done," we need to check if we actually built everything the Project Charter (our blueprint) said we needed to build. This document:

1. **Lists what was promised** (from the Project Charter)
2. **Checks what was actually built** (by reading all the code)
3. **Identifies gaps** (things missing or broken)
4. **Creates a fix plan** (GitHub Issues ready to be worked on)

---

## Executive Summary

### What's Working Well ✅

| Feature | Status | Notes |
|---------|--------|-------|
| Go BFF + React Architecture | ✅ Complete | Single binary deployment working |
| Token Ledger Storage | ✅ Complete | SQLite schema and API functional |
| Token Preview Meter | ✅ Complete | Color-coded meter in Architect view |
| Command Cards CRUD | ✅ Complete | Create, Read, Delete working |
| Flows Editor UI | ✅ Complete | React Flow canvas with drag/drop |
| Keyring Security | ✅ Complete | OS-native keyring via `go-keyring` |
| LLM Gateway | ✅ Complete | OpenAI + Anthropic providers |
| Agent Prompts | ✅ Complete | 4 agent roles defined |
| Optimization Suggestions | ✅ Complete | Analysis algorithms working |

### What's Missing or Broken 🚨

| Issue | Severity | Category |
|-------|----------|----------|
| WebSocket is just echo, no real signaling | **CRITICAL** | Architecture |
| No GitHub fallback for WebSocket failures | **CRITICAL** | Reliability |
| CORS allows ALL origins (`CheckOrigin: true`) | **CRITICAL** | Security |
| No HTTPS/WSS enforcement | **HIGH** | Security |
| "Click to Apply" optimization is placeholder | **HIGH** | Feature |
| Command execution doesn't capture latency | **MEDIUM** | Accuracy |
| Flow Editor doesn't persist to database | **MEDIUM** | Feature |
| FlowList uses mock data, not real API | **MEDIUM** | Feature |
| No mobile responsiveness test for Architect | **LOW** | Testing |
| Token estimation is naive (len/4) | **LOW** | Accuracy |

### Bottom Line

The application has about **70% of the core requirements** implemented correctly. However, there are **3 critical security/architecture issues** that must be fixed before this can be considered production-ready. The testing coverage is decent but has gaps, particularly around edge cases and error handling.

---

## Requirements Recap (From Project Charter)

Here's what the Project Charter promised we'd build. I'm using simple language so anyone can understand.

### Core UX Features

| Requirement | Description (Simple) |
|-------------|----------------------|
| **Token Preview Meter** | As you type, show how many "tokens" (AI currency) your message will cost. Use colors: Green = cheap, Yellow = getting expensive, Red = too expensive! |
| **Audit Flow Playback** | Show suggestions for saving money on AI calls, with a button you can click to automatically apply the savings. |

### Security Requirements

| Requirement | Description (Simple) |
|-------------|----------------------|
| **Token Ledger Integrity** | Parse messy AI output correctly so we count tokens accurately. |
| **WebSocket Signaling** | Use real-time connections for communication, but have a backup plan (check GitHub commits) if the connection fails. |
| **Secure API Key Storage** | Never show API keys in the browser. Store them safely in the computer's secure vault (Keychain on Mac, Credential Manager on Windows). |
| **Transport & Access Control** | Use encrypted connections (HTTPS). Don't let random websites run commands on our app. |

### Agent Roles (4 Simulated AI Personas)

| Agent | Job |
|-------|-----|
| **Orchestrator (Planner)** | Breaks big goals into small tasks |
| **Implementation (Dev)** | Writes the actual code |
| **Test (QA)** | Checks if the code works |
| **Token Optimizer (Auditor)** | Finds ways to save money |

---

## Gap Analysis

### Critical Gaps

These are the "fix immediately" problems. Like finding your house has no smoke detectors.

---

#### GAP-001: WebSocket Implementation is Non-Functional

**Where:** `internal/server/websocket.go`

**What's Wrong:**  
The WebSocket handler just echoes messages back. It doesn't do any real work like:
- Broadcasting flow execution status to all connected clients
- Sending real-time updates when token usage changes
- Notifying when optimizations are available

**Current Code (Problem):**
```go
// For now, we just echo messages back
for {
    mt, message, err := conn.ReadMessage()
    // ... just echoes it back
    err = conn.WriteMessage(mt, message)
}
```

**Why It Matters:**  
The Project Charter specifically says: "Go Flow Engine uses WebSockets for primary signaling." Without real signaling, the frontend has no way to know when:
- A flow finishes executing
- An agent completes its work
- Token counts are updated

**Impact:** Users have to manually refresh to see updates. That's like having a phone that only works if you keep hanging up and calling back.

---

#### GAP-002: No GitHub Commit Fallback for WebSocket Failures

**Where:** Missing from codebase entirely

**What's Wrong:**  
The Charter says: "implement a GitHub Commit Check as a fallback if the signal fails." This code doesn't exist anywhere.

**Why It Matters:**  
If the WebSocket connection drops (internet hiccup, browser tab sleeps), the system has no backup plan. Flows could get "stuck" waiting for a signal that will never come.

**Impact:** System reliability is compromised. In a real multi-agent workflow, a stuck flow means wasted money and time.

---

#### GAP-003: CORS Security is Completely Disabled

**Where:** `internal/server/websocket.go`, line 14

**What's Wrong:**
```go
CheckOrigin: func(r *http.Request) bool {
    return true // Allow all origins for now ← DANGER!
}
```

This allows ANY website in the world to connect to your WebSocket and potentially run commands.

**Why It Matters:**  
The Charter specifically says: "Go enforces strict CORS restrictions." This does the opposite - it enforces NO restrictions.

**Attack Scenario:**
1. You visit a malicious website
2. That website secretly connects to your local Forge Orchestrator
3. It runs AI commands using YOUR API keys
4. You pay the bill

**Impact:** This is a security vulnerability that could lead to:
- Unauthorized command execution
- API key theft (indirect)
- Financial loss from unauthorized AI API usage

---

### Security Gaps

---

#### GAP-004: No HTTPS/WSS Enforcement

**Where:** `main.go`, line 67

**Current Code:**
```go
if err := http.ListenAndServe(":8080", mux); err != nil {
```

**What's Wrong:**  
The server only runs on HTTP (unencrypted). The Charter says: "HTTPS/WSS used locally."

**Why It Matters:**  
Without HTTPS:
- Anyone on your network can see your API keys if transmitted
- Man-in-the-middle attacks are possible
- Sensitive prompts and responses travel in plain text

**Note:** For localhost development, this is lower severity. But the architecture should support HTTPS for production or remote usage.

---

#### GAP-005: HTTP Endpoints Missing CORS Headers

**Where:** `internal/server/routes.go`

**What's Wrong:**  
While WebSocket has CORS (badly configured), the HTTP API endpoints have NO CORS configuration at all.

**Current Situation:**
- No `Access-Control-Allow-Origin` headers
- No preflight request handling
- Frontend works only because it's served from same origin

**Impact:** If someone deploys the frontend separately from the backend, it won't work. Also, no protection against cross-origin attacks on HTTP endpoints.

---

### Testing Gaps

---

#### GAP-006: Missing Error Handling Tests

**Where:** Multiple test files

**What's Missing:**

| File | Missing Tests |
|------|---------------|
| `commands_test.go` | Test for malformed JSON body |
| `commands_test.go` | Test for SQL injection attempts |
| `ledger_test.go` | Test for invalid timestamp format |
| `flows/engine_test.go` | Test for flow with no nodes |
| `flows/engine_test.go` | Test for flow with invalid provider |

**Why It Matters:**  
Tests only cover "happy path" (everything works). Real apps need to handle when things go wrong.

---

#### GAP-007: Frontend E2E Tests Don't Cover Real API Integration

**Where:** `frontend/tests/e2e/*.spec.ts`

**What's Wrong:**  
Almost every E2E test mocks the API responses:

```typescript
await page.route('/api/commands', async route => {
    await route.fulfill({
        status: 200,
        // ... fake data
    });
});
```

**Why It Matters:**  
These tests prove the UI works with fake data. They don't prove the UI works with the REAL backend. If there's a mismatch between what the frontend expects and what the backend sends, these tests won't catch it.

**Judgment Call Made:** I've classified this as MEDIUM severity because the backend has its own unit tests. But ideally, you'd have some integration tests that use the real API.

---

#### GAP-008: No Mobile Viewport Test for Architect View

**Where:** Missing from `tests/e2e/architect.spec.ts`

**What's Wrong:**  
The Charter says: "Confirmed responsiveness on a mobile viewport (360px width check)."

The `layout.spec.ts` tests mobile sidebar hiding, but the Architect view (the most important input area) isn't tested on mobile.

---

### Code Quality Issues

---

#### GAP-009: Flow Editor Doesn't Actually Save to Database

**Where:** `frontend/src/components/Flows/FlowEditor.tsx`, line 82

**Current Code:**
```typescript
const handleSave = () => {
    if (reactFlowInstance) {
        const flow = reactFlowInstance.toObject();
        console.log('Saving flow:', flow);
        // TODO: Call API to save flow ← Never implemented!
        alert('Flow saved! (Check console for object)');
    }
};
```

**Impact:** Users think they saved their flow, but it's just logged to console and lost on page refresh.

---

#### GAP-010: FlowList Uses Hardcoded Mock Data

**Where:** `frontend/src/components/Flows/FlowList.tsx`, lines 12-28

**Current Code:**
```typescript
const MOCK_FLOWS: Flow[] = [
    {
        id: '1',
        name: 'Customer Onboarding',
        // ... hardcoded fake data
    },
];
```

**Impact:** The Flows list doesn't show real flows from the database. The backend API exists (`GET /api/flows`) but the frontend doesn't call it.

---

#### GAP-011: "Apply Optimization" is Placeholder

**Where:** `internal/server/optimizer.go`, line 36-58

**Current Code:**
```go
// In a real implementation, we would fetch the suggestion by ID and apply it.
// Here, we just log that it was applied.
```

**Impact:** The "Click to Apply" button in the UI does record that you clicked it, but it doesn't actually change any flow configurations or settings.

---

#### GAP-012: Token Estimation is Naive

**Where:** `internal/server/ledger.go`, line 61

**Current Code:**
```go
// Simple approximation: len(text) / 4
count := len(req.Text) / 4
```

**Why It's Not Great:**  
Different languages and content types have different token densities:
- Code has different tokenization than prose
- Non-ASCII characters may use more tokens
- The ratio varies between OpenAI and Anthropic tokenizers

**Impact:** Token estimates shown in the UI may be off by 20-50% from actual usage.

---

#### GAP-013: Latency Not Captured in Command Execution

**Where:** `internal/server/commands.go`, line 170-171

**Current Code:**
```go
// Note: Latency is not captured here yet, could add timing around ExecutePrompt
```

**Impact:** The ledger shows 0 for latency on command executions, making performance analysis incomplete.

---

#### GAP-014: Agent Role Names Inconsistent

**Where:** `internal/agents/agent_prompts.go` vs usage in flows

**Issue:**  
- Prompts define: `"Architect"`, `"Implementation"`, `"Test"`, `"Optimizer"`
- Flow data uses: `"role": "coder"`, `"role": "planner"`, `"role": "tester"`

**Impact:** If someone creates a flow with `role: "coder"`, the `GetAgentPrompt` function returns an error because it doesn't recognize that role.

---

## GitHub Issue Contracts

Below are ready-to-use GitHub Issue templates in the format specified by the Project Charter.

---

### Issue #032: Implement Real WebSocket Signaling Hub

**Priority:** 🔴 CRITICAL

**1. 🎫 Contract Summary**

Implement a proper WebSocket hub pattern to enable real-time bidirectional communication between the Go backend and React frontend. The current echo-only implementation must be replaced with:
- A Hub that manages multiple client connections
- Broadcast capability for flow execution updates
- Per-client message routing for agent status updates

**2. 📊 Acceptance Criteria**

- [ ] Create `Hub` struct to manage connected clients
- [ ] Create `Client` struct with write pump and read pump goroutines
- [ ] Implement `Broadcast()` method to send messages to all clients
- [ ] Implement `SendToClient()` method for targeted messages
- [ ] Update flow execution to broadcast status updates: `STARTED`, `NODE_COMPLETE`, `FINISHED`, `FAILED`
- [ ] Add reconnection handling (client-side) with exponential backoff
- [ ] Write unit tests for Hub registration/unregistration
- [ ] Write E2E test verifying real-time updates in Ledger view

**3. 🏞️ Technical Specification**

Message format:
```json
{
    "type": "FLOW_STATUS",
    "payload": {
        "flowId": 1,
        "status": "NODE_COMPLETE",
        "nodeId": "agent-1",
        "timestamp": "2025-12-04T10:00:00Z"
    }
}
```

**4. 🔗 Related Files**
- `internal/server/websocket.go` (rewrite)
- `internal/flows/engine.go` (add broadcast calls)
- `frontend/src/hooks/useWebSocket.ts` (new file)

---

### Issue #033: Add GitHub Commit Fallback for WebSocket Failures

**Priority:** 🔴 CRITICAL

**1. 🎫 Contract Summary**

Per Project Charter requirement: "implement a GitHub Commit Check as a fallback if the signal fails."

When WebSocket connection is unavailable for >30 seconds, the flow engine should poll a GitHub commit status or use a file-based heartbeat as a fallback signaling mechanism.

**2. 📊 Acceptance Criteria**

- [ ] Create `FallbackSignaler` interface with `CheckFlowStatus(flowID int) (string, error)`
- [ ] Implement file-based fallback: write status to `.forge/flow_{id}_status.json`
- [ ] Implement GitHub-based fallback (optional): check commit status API
- [ ] Flow engine attempts WebSocket first, falls back after timeout
- [ ] Frontend polls `/api/flows/{id}/status` when WebSocket disconnected
- [ ] Write integration test simulating WebSocket failure

**3. 🔗 Related Files**
- `internal/flows/signaler.go` (new file)
- `internal/flows/engine.go` (add fallback logic)

---

### Issue #034: Fix CORS Security Vulnerability

**Priority:** 🔴 CRITICAL

**1. 🎫 Contract Summary**

The WebSocket upgrader currently allows all origins (`CheckOrigin: return true`). This must be replaced with a strict whitelist that only allows:
- `http://localhost:*` (development)
- Configurable production origins via environment variable

**2. 📊 Acceptance Criteria**

- [ ] Create `AllowedOrigins` configuration via environment variable `FORGE_ALLOWED_ORIGINS`
- [ ] Default to `http://localhost:8080,http://localhost:5173` for development
- [ ] Implement origin checking function that validates against whitelist
- [ ] Return 403 Forbidden for unauthorized origins
- [ ] Add HTTP CORS middleware for API routes with same restrictions
- [ ] Write test for blocked cross-origin request
- [ ] Document configuration in README

**3. 🔗 Related Files**
- `internal/server/websocket.go`
- `internal/server/middleware.go` (new file)
- `internal/server/routes.go`

---

### Issue #035: Connect Flow Editor to Backend API

**Priority:** 🟡 HIGH

**1. 🎫 Contract Summary**

The Flow Editor UI has working drag-and-drop but doesn't persist to the database. The FlowList shows hardcoded mock data instead of calling the API.

**2. 📊 Acceptance Criteria**

- [ ] FlowList: Replace `MOCK_FLOWS` with `fetch('/api/flows')`
- [ ] FlowList: Add loading and error states
- [ ] FlowEditor: Implement `handleSave()` to call `POST /api/flows` (new) or `PUT /api/flows/{id}` (existing)
- [ ] FlowEditor: Load existing flow data when `id` param is present
- [ ] FlowEditor: Implement `handleExecute()` to call `POST /api/flows/{id}/execute`
- [ ] Add success/error toast notifications
- [ ] Write E2E test for full CRUD flow

**3. 🔗 Related Files**
- `frontend/src/components/Flows/FlowList.tsx`
- `frontend/src/components/Flows/FlowEditor.tsx`

---

### Issue #036: Implement Real Optimization Apply Logic

**Priority:** 🟡 HIGH

**1. 🎫 Contract Summary**

The "Apply" button for optimizations logs to the ledger but doesn't actually modify flow configurations. Implement the actual apply logic for each optimization type.

**2. 📊 Acceptance Criteria**

- [ ] Parse `ApplyAction` JSON field to determine action type
- [ ] For `model_switch`: Update flow node's provider field in `forge_flows.data`
- [ ] For `prompt_optimization`: Flag the prompt for user review (can't auto-modify prompts)
- [ ] For `retry_strategy`: Add retry configuration to flow metadata
- [ ] Return updated flow configuration in response
- [ ] Frontend: Show confirmation modal with "before/after" preview
- [ ] Write unit test for each optimization type

**3. 🔗 Related Files**
- `internal/server/optimizer.go`
- `internal/optimizer/applier.go` (new file)
- `frontend/src/components/Ledger/OptimizationCard.tsx`

---

### Issue #037: Add HTTPS Support

**Priority:** 🟡 HIGH

**1. 🎫 Contract Summary**

Add optional HTTPS/TLS support via configuration. The server should support both HTTP (for easy development) and HTTPS (for production/remote access).

**2. 📊 Acceptance Criteria**

- [ ] Add `FORGE_TLS_CERT` and `FORGE_TLS_KEY` environment variables
- [ ] If both are set, use `http.ListenAndServeTLS()`
- [ ] If neither is set, fall back to HTTP with warning log
- [ ] Auto-generate self-signed cert for development (optional flag `--dev-tls`)
- [ ] Update frontend WebSocket to use `wss://` when on HTTPS
- [ ] Document TLS setup in README

**3. 🔗 Related Files**
- `main.go`
- `frontend/src/hooks/useWebSocket.ts`

---

### Issue #038: Fix Agent Role Name Mapping

**Priority:** 🟡 MEDIUM

**1. 🎫 Contract Summary**

The agent prompt system uses formal role names (`Architect`, `Implementation`) while flow JSON uses informal names (`planner`, `coder`). Add a mapping layer.

**2. 📊 Acceptance Criteria**

- [ ] Create role alias map:
  - `planner` → `Architect`
  - `coder` → `Implementation`  
  - `tester` → `Test`
  - `auditor` → `Optimizer`
- [ ] Update `GetAgentPrompt()` to check alias map before returning error
- [ ] Update Flow Editor UI to use consistent role names
- [ ] Write unit test for alias resolution

**3. 🔗 Related Files**
- `internal/agents/agent_prompts.go`
- `frontend/src/components/Flows/FlowEditor.tsx`

---

### Issue #039: Add Command Execution Latency Tracking

**Priority:** 🟢 MEDIUM

**1. 🎫 Contract Summary**

Command execution doesn't record latency in the ledger, making performance analysis incomplete.

**2. 📊 Acceptance Criteria**

- [ ] Wrap `ExecutePrompt` call with `time.Now()` before and `time.Since()` after
- [ ] Store `latency_ms` in ledger entry
- [ ] Add latency column to Ledger UI table
- [ ] Add latency chart to a new "Analytics" section (optional)

**3. 🔗 Related Files**
- `internal/server/commands.go`
- `frontend/src/components/Ledger/LedgerView.tsx`

---

### Issue #040: Improve Token Estimation Accuracy

**Priority:** 🟢 LOW

**1. 🎫 Contract Summary**

The current `len(text)/4` estimation is inaccurate. Implement proper tokenization or improve the heuristic.

**2. 📊 Acceptance Criteria**

- [ ] Research: Determine if we want to add tiktoken (Go port) as dependency
- [ ] Option A: Use `tiktoken-go` for accurate OpenAI token counting
- [ ] Option B: Improve heuristic with word-based + special character adjustments
- [ ] Add provider selection to estimation endpoint for provider-specific counts
- [ ] Write benchmark test comparing estimated vs actual tokens

**3. 🔗 Related Files**
- `internal/server/ledger.go`
- `frontend/src/components/Architect/TokenMeter.tsx`

---

### Issue #041: Add Integration Test Suite

**Priority:** 🟢 LOW

**1. 🎫 Contract Summary**

Current E2E tests mock all API calls. Add a separate integration test suite that tests frontend + backend together.

**2. 📊 Acceptance Criteria**

- [ ] Create `tests/integration/` directory
- [ ] Add Go test that starts real server on random port
- [ ] Add Playwright config for integration tests pointing to real server
- [ ] Write 3 critical path tests:
  - Create command → Run command → Verify ledger entry
  - Create flow → Execute flow → Verify ledger entries
  - Set API key → Verify status shows "Configured"
- [ ] Add to CI pipeline (run after unit tests)

**3. 🔗 Related Files**
- `tests/integration/full_flow_test.go` (new)
- `frontend/playwright.integration.config.ts` (new)

---

### Issue #042: Add Mobile Responsiveness Tests for Architect View

**Priority:** 🟢 LOW

**1. 🎫 Contract Summary**

Per Charter requirement: "Confirmed responsiveness on a mobile viewport (360px width check)."

**2. 📊 Acceptance Criteria**

- [ ] Add test in `architect.spec.ts` with viewport 360x640
- [ ] Verify textarea is visible and usable
- [ ] Verify TokenMeter is visible below textarea
- [ ] Test typing and paste operations work on mobile viewport

**3. 🔗 Related Files**
- `frontend/tests/e2e/architect.spec.ts`

---

## Appendix: Judgment Calls

During this analysis, I made the following judgment calls where the requirements were ambiguous:

1. **WebSocket fallback scope:** The Charter mentions "GitHub Commit Check" but doesn't specify the exact mechanism. I interpreted this as any reliable out-of-band status check, including file-based fallback as a simpler alternative.

2. **HTTPS severity:** Classified as HIGH not CRITICAL because localhost development typically doesn't require TLS. However, any remote/production deployment absolutely needs it.

3. **E2E test mocking:** I classified "tests only use mocks" as MEDIUM because the backend has its own unit tests. A stricter interpretation would call this HIGH since end-to-end should mean *end to end*.

4. **Token estimation:** Classified as LOW because users can see actual token counts after execution. The preview is informational, not transactional.

5. **Flow Editor persistence:** Classified as MEDIUM because flows can still be created via API. But the UX expectation is that the UI Save button works, so this is deceptive.

---

*End of Analysis Document*
