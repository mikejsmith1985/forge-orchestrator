# Issue #046: Add Flow Node Configuration UI

**Priority:** 🟢 MEDIUM  
**Estimated Tokens:** ~1,800 (Medium complexity)  
**Agent Role:** Implementation

---

## 1. 🎫 Related Issue Context

Per Project Charter Feature Evolution: "Flows Editor: Build multi-step, time-aware, multi-agent Orchestration Pipelines."

The current Flow Editor allows drag/drop of nodes, but nodes cannot be configured with:
- Agent role selection
- Provider selection
- Prompt input
- Retry settings

---

## 2. 📋 Acceptance Criteria

### Node Configuration Panel
- [ ] Click on node opens configuration sidebar/modal
- [ ] Form fields: Label (text), Role (dropdown), Provider (dropdown), Prompt (textarea)
- [ ] Save button updates node data
- [ ] Cancel button closes without saving

### Role Dropdown Options
- [ ] Planner / Architect
- [ ] Developer / Coder
- [ ] QA / Tester
- [ ] Auditor / Optimizer

### Provider Dropdown Options
- [ ] Show only configured providers (from `/api/keys/status`)
- [ ] Disabled state if no providers configured

### Prompt Field
- [ ] Multiline textarea
- [ ] Token count preview (reuse TokenMeter component)
- [ ] Placeholder with example

### Tests
- [ ] E2E test: Add node → Configure → Save → Node shows updated label
- [ ] E2E test: Unconfigured node shows warning indicator

---

## 3. 📊 Token Efficiency Strategy

- Create reusable NodeConfigPanel component
- Reuse existing TokenMeter component
- Leverage shadcn/ui form components

---

## 4. 🏗️ Technical Specification

### Node Data Structure
```typescript
interface AgentNodeData {
    label: string;
    role: 'Architect' | 'Implementation' | 'Test' | 'Optimizer';
    provider: 'Anthropic' | 'OpenAI';
    prompt: string;
    retryConfig?: {
        maxRetries: number;
        backoffMs: number;
    };
}
```

### Custom Node Component
```typescript
// frontend/src/components/Flows/AgentNode.tsx
import { Handle, Position } from 'reactflow';

function AgentNode({ data, selected }: NodeProps<AgentNodeData>) {
    const isConfigured = data.role && data.provider && data.prompt;
    
    return (
        <div className={`
            p-4 rounded-lg border-2 
            ${selected ? 'border-blue-500' : 'border-slate-600'}
            ${isConfigured ? 'bg-slate-800' : 'bg-yellow-900/20'}
        `}>
            <Handle type="target" position={Position.Top} />
            
            <div className="flex items-center gap-2">
                <div className={`w-3 h-3 rounded-full ${
                    isConfigured ? 'bg-green-500' : 'bg-yellow-500'
                }`} />
                <span className="text-white font-medium">{data.label}</span>
            </div>
            
            {data.role && (
                <div className="text-xs text-slate-400 mt-1">{data.role}</div>
            )}
            
            <Handle type="source" position={Position.Bottom} />
        </div>
    );
}
```

### Configuration Panel
```typescript
// frontend/src/components/Flows/NodeConfigPanel.tsx
interface Props {
    node: Node<AgentNodeData>;
    onSave: (data: AgentNodeData) => void;
    onClose: () => void;
}

function NodeConfigPanel({ node, onSave, onClose }: Props) {
    const [formData, setFormData] = useState<AgentNodeData>(node.data);
    const [providers, setProviders] = useState<string[]>([]);
    
    useEffect(() => {
        fetch('/api/keys/status')
            .then(res => res.json())
            .then(data => {
                const configured = Object.entries(data)
                    .filter(([_, isSet]) => isSet)
                    .map(([provider]) => provider);
                setProviders(configured);
            });
    }, []);
    
    return (
        <div className="w-80 bg-slate-800 border-l border-slate-700 p-4">
            <h3 className="text-lg font-semibold text-white mb-4">
                Configure Node
            </h3>
            
            <div className="space-y-4">
                <div>
                    <label className="text-sm text-slate-400">Label</label>
                    <input
                        type="text"
                        value={formData.label}
                        onChange={(e) => setFormData({...formData, label: e.target.value})}
                        className="w-full bg-slate-900 border border-slate-700 rounded px-3 py-2 text-white"
                    />
                </div>
                
                <div>
                    <label className="text-sm text-slate-400">Agent Role</label>
                    <select
                        value={formData.role}
                        onChange={(e) => setFormData({...formData, role: e.target.value})}
                        className="w-full bg-slate-900 border border-slate-700 rounded px-3 py-2 text-white"
                    >
                        <option value="">Select role...</option>
                        <option value="Architect">Planner / Architect</option>
                        <option value="Implementation">Developer / Coder</option>
                        <option value="Test">QA / Tester</option>
                        <option value="Optimizer">Auditor / Optimizer</option>
                    </select>
                </div>
                
                <div>
                    <label className="text-sm text-slate-400">Provider</label>
                    <select
                        value={formData.provider}
                        onChange={(e) => setFormData({...formData, provider: e.target.value})}
                        className="w-full bg-slate-900 border border-slate-700 rounded px-3 py-2 text-white"
                        disabled={providers.length === 0}
                    >
                        <option value="">Select provider...</option>
                        {providers.map(p => (
                            <option key={p} value={p}>{p}</option>
                        ))}
                    </select>
                    {providers.length === 0 && (
                        <p className="text-xs text-yellow-400 mt-1">
                            Configure API keys in Settings first
                        </p>
                    )}
                </div>
                
                <div>
                    <label className="text-sm text-slate-400">Prompt</label>
                    <textarea
                        value={formData.prompt}
                        onChange={(e) => setFormData({...formData, prompt: e.target.value})}
                        placeholder="Enter the task for this agent..."
                        rows={4}
                        className="w-full bg-slate-900 border border-slate-700 rounded px-3 py-2 text-white resize-none"
                    />
                    <TokenMeter tokenCount={Math.ceil(formData.prompt.length / 4)} maxTokens={4000} />
                </div>
            </div>
            
            <div className="flex gap-2 mt-6">
                <button
                    onClick={() => onSave(formData)}
                    className="flex-1 bg-blue-600 hover:bg-blue-500 text-white py-2 rounded"
                >
                    Save
                </button>
                <button
                    onClick={onClose}
                    className="flex-1 bg-slate-700 hover:bg-slate-600 text-white py-2 rounded"
                >
                    Cancel
                </button>
            </div>
        </div>
    );
}
```

---

## 5. 📁 Files to Create/Modify

| Action | File |
|--------|------|
| CREATE | `frontend/src/components/Flows/AgentNode.tsx` |
| CREATE | `frontend/src/components/Flows/NodeConfigPanel.tsx` |
| MODIFY | `frontend/src/components/Flows/FlowEditor.tsx` |
| CREATE | `frontend/tests/e2e/flow-config.spec.ts` |

---

## 6. ✅ Definition of Done

1. Clicking a node opens configuration panel
2. All form fields save correctly to node data
3. Provider dropdown only shows configured providers
4. Unconfigured nodes show yellow warning indicator
5. Token count updates as prompt is typed
6. E2E tests pass
