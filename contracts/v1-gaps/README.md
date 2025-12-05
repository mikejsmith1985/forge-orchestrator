# v1 Gap Remediation - Issue Index

## Overview

This directory contains GitHub Issue contracts for all gaps identified in the v1.0.0 analysis.
Each issue is scoped to be completable in a single chat session per the Project Charter methodology.

## Issue Summary by Priority

### 🔴 CRITICAL (Fix Immediately)

| Issue | Title | Est. Tokens | Depends On |
|-------|-------|-------------|------------|
| #032 | [WebSocket Hub Implementation](./issue_032_websocket_hub.md) | ~2,500 | - |
| #033 | [WebSocket Fallback Signaling](./issue_033_websocket_fallback.md) | ~1,800 | #032 |
| #034 | [Fix CORS Security Vulnerability](./issue_034_cors_security.md) | ~1,500 | - |

### 🟡 HIGH (Fix Before Next Release)

| Issue | Title | Est. Tokens | Depends On |
|-------|-------|-------------|------------|
| #035 | [Connect Flow Editor to Backend API](./issue_035_flow_editor_api.md) | ~2,000 | - |
| #036 | [Implement Real Optimization Apply](./issue_036_optimization_apply.md) | ~2,200 | - |
| #037 | [Add HTTPS/TLS Support](./issue_037_https_support.md) | ~1,200 | - |
| #043 | [HTTP CORS Middleware](./issue_043_http_cors.md) | ~800 | #034 |
| #045 | [WebSocket Hub + Flow Engine Integration](./issue_045_hub_integration.md) | ~1,500 | #032 |

### 🟡 MEDIUM

| Issue | Title | Est. Tokens | Depends On |
|-------|-------|-------------|------------|
| #038 | [Fix Agent Role Name Mapping](./issue_038_role_mapping.md) | ~800 | - |
| #039 | [Add Latency Tracking](./issue_039_latency_tracking.md) | ~1,000 | - |
| #044 | [Add Error Handling Tests](./issue_044_error_tests.md) | ~1,200 | - |
| #046 | [Flow Node Configuration UI](./issue_046_node_config_ui.md) | ~1,800 | #035 |

### 🟢 LOW

| Issue | Title | Est. Tokens | Depends On |
|-------|-------|-------------|------------|
| #040 | [Improve Token Estimation](./issue_040_token_estimation.md) | ~1,500 | - |
| #041 | [Add Integration Test Suite](./issue_041_integration_tests.md) | ~1,800 | - |
| #042 | [Mobile Responsiveness Tests](./issue_042_mobile_tests.md) | ~600 | - |

---

## Recommended Execution Order

Based on dependencies and priority, execute issues in this order:

### Phase 1: Critical Security & Architecture (Week 1)
1. **#034** - Fix CORS Security (no dependencies, critical security fix)
2. **#032** - WebSocket Hub (foundation for real-time features)
3. **#043** - HTTP CORS Middleware (builds on #034)
4. **#033** - WebSocket Fallback (requires #032)

### Phase 2: Core Feature Completion (Week 2)
5. **#035** - Flow Editor API Connection
6. **#045** - WebSocket + Flow Integration (requires #032)
7. **#036** - Optimization Apply Logic
8. **#037** - HTTPS Support

### Phase 3: Polish & Quality (Week 3)
9. **#038** - Role Name Mapping
10. **#039** - Latency Tracking
11. **#046** - Node Configuration UI (requires #035)
12. **#044** - Error Handling Tests

### Phase 4: Enhancement (Week 4+)
13. **#040** - Token Estimation Improvement
14. **#041** - Integration Test Suite
15. **#042** - Mobile Tests

---

## Token Budget Summary

| Priority | Issues | Total Est. Tokens |
|----------|--------|-------------------|
| CRITICAL | 3 | ~5,800 |
| HIGH | 5 | ~7,700 |
| MEDIUM | 4 | ~4,800 |
| LOW | 3 | ~3,900 |
| **TOTAL** | **15** | **~22,200** |

Assuming ~3,000 tokens per chat session average, this represents approximately **8-10 chat sessions** to complete all gaps.

---

## Contract Format

Each issue follows the Project Charter template:

1. **🎫 Related Issue Context** - Gap reference and problem description
2. **📋 Acceptance Criteria** - Checkboxes for completion
3. **📊 Token Efficiency Strategy** - How to minimize token usage
4. **🏗️ Technical Specification** - Code examples and architecture
5. **📁 Files to Create/Modify** - Clear file list
6. **✅ Definition of Done** - Verification checklist
