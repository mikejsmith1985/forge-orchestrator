# Issue #035: Connect Flow Editor UI to Backend API

**Priority:** 🟡 HIGH  
**Estimated Tokens:** ~2,000 (Medium complexity)  
**Agent Role:** Implementation

---

## 1. 🎫 Related Issue Context

**Gap References:** GAP-009, GAP-010 from v1-analysis.md

**Problem 1 - FlowEditor doesn't save:**
```typescript
const handleSave = () => {
    console.log('Saving flow:', flow);
    // TODO: Call API to save flow ← Never implemented!
    alert('Flow saved! (Check console for object)');
};
```

**Problem 2 - FlowList uses mock data:**
```typescript
const MOCK_FLOWS: Flow[] = [
    { id: '1', name: 'Customer Onboarding', ... } // Hardcoded!
];
```

---

## 2. 📋 Acceptance Criteria

### FlowList Component
- [ ] Remove `MOCK_FLOWS` constant
- [ ] Add `useEffect` to fetch from `GET /api/flows` on mount
- [ ] Add loading state with spinner
- [ ] Add error state with retry button
- [ ] Delete button calls `DELETE /api/flows/{id}` and refetches list

### FlowEditor Component
- [ ] On mount: if `id` param exists, fetch `GET /api/flows/{id}` and load nodes/edges
- [ ] `handleSave()`: Call `POST /api/flows` (new) or `PUT /api/flows/{id}` (existing)
- [ ] `handleExecute()`: Call `POST /api/flows/{id}/execute`
- [ ] Add loading states during save/execute
- [ ] Show success toast on save, navigate back to list
- [ ] Show error toast on failure

### Data Transformation
- [ ] Transform ReactFlow format to backend format on save
- [ ] Transform backend format to ReactFlow format on load

### Tests
- [ ] E2E test: Create new flow → Save → Appears in list
- [ ] E2E test: Load existing flow → Modify → Save → Changes persist
- [ ] E2E test: Delete flow → Removed from list

---

## 3. 📊 Token Efficiency Strategy

- Focus only on FlowList.tsx and FlowEditor.tsx
- Backend API already exists and works (verified in analysis)
- Reuse existing fetch patterns from CommandDeck.tsx

---

## 4. 🏗️ Technical Specification

### Backend Flow Data Structure (existing)
```json
{
    "id": 1,
    "name": "My Flow",
    "data": "{\"nodes\":[...],\"edges\":[...]}",
    "status": "active",
    "created_at": "2025-12-04T10:00:00Z"
}
```

### ReactFlow to Backend Transform
```typescript
const saveFlow = async () => {
    const flowData = {
        name: flowName,
        data: JSON.stringify({
            nodes: nodes,
            edges: edges
        }),
        status: 'active'
    };
    
    if (id) {
        await fetch(`/api/flows/${id}`, {
            method: 'PUT',
            body: JSON.stringify(flowData)
        });
    } else {
        await fetch('/api/flows', {
            method: 'POST',
            body: JSON.stringify(flowData)
        });
    }
};
```

### Backend to ReactFlow Transform
```typescript
const loadFlow = async (id: string) => {
    const response = await fetch(`/api/flows/${id}`);
    const flow = await response.json();
    const graphData = JSON.parse(flow.data);
    
    setNodes(graphData.nodes);
    setEdges(graphData.edges);
    setFlowName(flow.name);
};
```

---

## 5. 📁 Files to Create/Modify

| Action | File |
|--------|------|
| MODIFY | `frontend/src/components/Flows/FlowList.tsx` |
| MODIFY | `frontend/src/components/Flows/FlowEditor.tsx` |
| CREATE | `frontend/tests/e2e/flows-crud.spec.ts` |

---

## 6. ✅ Definition of Done

1. FlowList displays flows from real API (no mock data)
2. New flows can be created and appear in list after save
3. Existing flows load their saved state when opened
4. Changes to flows persist after save
5. Delete removes flow from database and list
6. All E2E tests pass
