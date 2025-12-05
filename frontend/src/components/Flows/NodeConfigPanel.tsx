import React, { useState, useEffect } from 'react';
import type { Node } from 'reactflow';
import { X } from 'lucide-react';
import { TokenMeter } from '../Architect/TokenMeter';
import type { AgentNodeData } from './AgentNode';

interface Props {
    node: Node<AgentNodeData>;
    onSave: (nodeId: string, data: AgentNodeData) => void;
    onClose: () => void;
}

const agentRoles = [
    { value: 'Architect', label: 'Planner / Architect' },
    { value: 'Implementation', label: 'Developer / Coder' },
    { value: 'Test', label: 'QA / Tester' },
    { value: 'Optimizer', label: 'Auditor / Optimizer' },
];

const NodeConfigPanel: React.FC<Props> = ({ node, onSave, onClose }) => {
    const [formData, setFormData] = useState<AgentNodeData>({
        label: node.data.label || '',
        role: node.data.role,
        provider: node.data.provider,
        prompt: node.data.prompt || '',
    });
    const [providers, setProviders] = useState<string[]>([]);
    const [loadingProviders, setLoadingProviders] = useState(true);

    useEffect(() => {
        fetch('/api/keys/status')
            .then((res) => res.json())
            .then((data) => {
                const configured = Object.entries(data)
                    .filter(([_, isSet]) => isSet)
                    .map(([provider]) => provider);
                setProviders(configured);
            })
            .catch((err) => {
                console.error('Failed to fetch providers:', err);
            })
            .finally(() => {
                setLoadingProviders(false);
            });
    }, []);

    const handleSave = () => {
        onSave(node.id, formData);
        onClose();
    };

    // Estimate token count (simple char/4 approximation)
    const estimatedTokens = Math.ceil((formData.prompt || '').length / 4);

    return (
        <div
            className="w-80 bg-slate-800 border-l border-slate-700 p-4 flex flex-col h-full"
            data-testid="node-config-panel"
        >
            <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-semibold text-white">Configure Node</h3>
                <button
                    onClick={onClose}
                    className="p-1 hover:bg-slate-700 rounded text-slate-400 hover:text-white"
                    data-testid="config-close-btn"
                >
                    <X size={20} />
                </button>
            </div>

            <div className="space-y-4 flex-1 overflow-y-auto">
                {/* Label Field */}
                <div>
                    <label className="block text-sm text-slate-400 mb-1">Label</label>
                    <input
                        type="text"
                        value={formData.label}
                        onChange={(e) => setFormData({ ...formData, label: e.target.value })}
                        className="w-full bg-slate-900 border border-slate-700 rounded px-3 py-2 text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                        placeholder="Node name"
                        data-testid="config-label-input"
                    />
                </div>

                {/* Role Dropdown */}
                <div>
                    <label className="block text-sm text-slate-400 mb-1">Agent Role</label>
                    <select
                        value={formData.role || ''}
                        onChange={(e) =>
                            setFormData({
                                ...formData,
                                role: e.target.value as AgentNodeData['role'],
                            })
                        }
                        className="w-full bg-slate-900 border border-slate-700 rounded px-3 py-2 text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                        data-testid="config-role-select"
                    >
                        <option value="">Select role...</option>
                        {agentRoles.map((role) => (
                            <option key={role.value} value={role.value}>
                                {role.label}
                            </option>
                        ))}
                    </select>
                </div>

                {/* Provider Dropdown */}
                <div>
                    <label className="block text-sm text-slate-400 mb-1">Provider</label>
                    <select
                        value={formData.provider || ''}
                        onChange={(e) => setFormData({ ...formData, provider: e.target.value })}
                        className="w-full bg-slate-900 border border-slate-700 rounded px-3 py-2 text-white focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
                        disabled={providers.length === 0 && !loadingProviders}
                        data-testid="config-provider-select"
                    >
                        <option value="">Select provider...</option>
                        {providers.map((p) => (
                            <option key={p} value={p}>
                                {p}
                            </option>
                        ))}
                    </select>
                    {providers.length === 0 && !loadingProviders && (
                        <p className="text-xs text-yellow-400 mt-1" data-testid="no-providers-warning">
                            Configure API keys in Settings first
                        </p>
                    )}
                </div>

                {/* Prompt Textarea */}
                <div>
                    <label className="block text-sm text-slate-400 mb-1">Prompt</label>
                    <textarea
                        value={formData.prompt || ''}
                        onChange={(e) => setFormData({ ...formData, prompt: e.target.value })}
                        placeholder="Enter the task for this agent...&#10;Example: Analyze the provided code and suggest improvements for performance and readability."
                        rows={6}
                        className="w-full bg-slate-900 border border-slate-700 rounded px-3 py-2 text-white resize-none focus:outline-none focus:ring-2 focus:ring-blue-500"
                        data-testid="config-prompt-textarea"
                    />
                    <div className="mt-2">
                        <TokenMeter tokenCount={estimatedTokens} maxTokens={4000} />
                    </div>
                </div>
            </div>

            {/* Action Buttons */}
            <div className="flex gap-2 mt-4 pt-4 border-t border-slate-700">
                <button
                    onClick={handleSave}
                    className="flex-1 bg-blue-600 hover:bg-blue-500 text-white py-2 rounded font-medium transition-colors"
                    data-testid="config-save-btn"
                >
                    Save
                </button>
                <button
                    onClick={onClose}
                    className="flex-1 bg-slate-700 hover:bg-slate-600 text-white py-2 rounded font-medium transition-colors"
                    data-testid="config-cancel-btn"
                >
                    Cancel
                </button>
            </div>
        </div>
    );
};

export default NodeConfigPanel;
