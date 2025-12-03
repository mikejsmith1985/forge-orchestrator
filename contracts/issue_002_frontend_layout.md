# Issue: Implement Application Layout & Forge Vision Container (Frontend)

**Role**: Implementation Agent (React/TS)
**Priority**: High
**Context**: Basic Vite app exists. We need the core layout structure to host the "Forge Vision" interface.

## 1. Requirements
1.  **Create `src/components/Layout/Sidebar.tsx`**:
    - Fixed width sidebar.
    - Placeholder navigation items (Dashboard, Flows, Settings).
    - Use Tailwind classes for styling (dark mode, glassmorphism).
2.  **Create `src/components/Layout/MainContent.tsx`**:
    - Flex-grow area for the main view.
3.  **Update `src/App.tsx`**:
    - Implement the Layout wrapper.
    - Ensure full viewport height (`h-screen`).
4.  **Responsiveness**:
    - Sidebar should collapse or hide on mobile (< 768px).

## 2. TDD & Verification Protocol
> [!IMPORTANT]
> You must follow this TDD workflow.
1.  **Create Test Contract**: Define what needs to be tested for the QE Agent (since you cannot run Playwright until the code exists, you define the *testability*).
    -   "Sidebar must have `data-testid='sidebar'`".
    -   "Mobile toggle must have `data-testid='mobile-menu-btn'`".
2.  **Implement Feature**: Write code to satisfy requirements and test IDs.
3.  **Visual Validation**:
    -   Verify Sidebar is visible on desktop.
    -   Verify Sidebar collapses on mobile.

## 3. Handoff & Deliverables
Upon completion, you must provide:
1.  **Committed Code**: `src/components/Layout/*`, `src/App.tsx`.
2.  **Token Efficiency Report**:
    -   Estimated Input Tokens: [Value]
    -   Actual Output Tokens: [Value]
    -   Optimization Strategy: [e.g., "Used Tailwind utility classes"]
3.  **WebSocket Signal**: Send signal `FRONTEND_LAYOUT_READY` (Simulated).

## 4. Acceptance Criteria
- [ ] `npm run build` passes.
- [ ] UI renders with Sidebar and Main Content.
- [ ] Responsive behavior confirmed (manual or via QE agent).
