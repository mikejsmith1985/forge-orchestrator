import React, { useState } from 'react';
import type { Node } from 'reactflow';
import { X, Terminal, Brain, AlertTriangle } from 'lucide-react';
import { TokenMeter } from '../Architect/TokenMeter';
import type { AgentNodeData } from './AgentNode';

/**
 * NodeConfigPanel - Redesigned for V2.1 Remediation Plan Task 3.1-3.3
 * 
 * CHANGES:
 * - Removed Agent Role and Provider dropdowns (Task 3.1)
 * - Simple text input for command/prompt string
 * - Two distinct node types: Shell Command and LLM Prompt (Task 3.2)
 * - Execution gating confirmation for Premium mode (Task 3.3)
 */

interface Props {
    node: Node<AgentNodeData>;
    onSave: (nodeId: string, data: AgentNodeData) => void;
    onClose: () => void;
}

const NodeConfigPanel: React.FC<Props> = ({ node, onSave, onClose }) => {
    const [formData, setFormData] = useState<AgentNodeData>({
        label: node.data.label || '',
        command: node.data.command || '',
        nodeType: node.data.nodeType || 'shell',
        premiumConfirmed: node.data.premiumConfirmed || false,
    });
    const [showPremiumWarning, setShowPremiumWarning] = useState(false);

    // Token estimate for LLM nodes
    const estimatedTokens = formData.nodeType === 'llm' 
        ? Math.ceil((formData.command || '').length / 4)
        : 0;

    const handleSave = () => {
        // Task 3.3: Show confirmation modal for Premium (LLM) nodes
        if (formData.nodeType === 'llm' && !formData.premiumConfirmed) {
            setShowPremiumWarning(true);
            return;
        }
        onSave(node.id, formData);
        onClose();
    };

    const confirmPremium = () => {
        const confirmedData = { ...formData, premiumConfirmed: true };
        onSave(node.id, confirmedData);
        setShowPremiumWarning(false);
        onClose();
    };

    return (
        <>
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

                    {/* Node Type Selection - Task 3.2 */}
                    <div>
                        <label className="block text-sm text-slate-400 mb-2">Node Type</label>
                        <div className="grid grid-cols-2 gap-2">
                            <button
                                type="button"
                                onClick={() => setFormData({ ...formData, nodeType: 'shell', premiumConfirmed: false })}
                                className={`flex items-center justify-center gap-2 p-3 rounded-lg border transition-colors ${
                                    formData.nodeType === 'shell'
                                        ? 'bg-green-500/20 border-green-500 text-green-300'
                                        : 'bg-slate-900 border-slate-700 text-slate-400 hover:border-slate-500'
                                }`}
                                data-testid="node-type-shell"
                            >
                                <Terminal size={18} />
                                <div className="text-left">
                                    <div className="font-medium text-sm">Shell</div>
                                    <div className="text-xs opacity-75">Zero-Token</div>
                                </div>
                            </button>
                            <button
                                type="button"
                                onClick={() => setFormData({ ...formData, nodeType: 'llm' })}
                                className={`flex items-center justify-center gap-2 p-3 rounded-lg border transition-colors ${
                                    formData.nodeType === 'llm'
                                        ? 'bg-purple-500/20 border-purple-500 text-purple-300'
                                        : 'bg-slate-900 border-slate-700 text-slate-400 hover:border-slate-500'
                                }`}
                                data-testid="node-type-llm"
                            >
                                <Brain size={18} />
                                <div className="text-left">
                                    <div className="font-medium text-sm">LLM</div>
                                    <div className="text-xs opacity-75">Premium</div>
                                </div>
                            </button>
                        </div>
                    </div>

                    {/* Command/Prompt Input - Task 3.1 */}
                    <div>
                        <label className="block text-sm text-slate-400 mb-1">
                            {formData.nodeType === 'llm' ? 'Prompt' : 'Command'}
                        </label>
                        <textarea
                            value={formData.command || ''}
                            onChange={(e) => setFormData({ ...formData, command: e.target.value, premiumConfirmed: false })}
                            placeholder={
                                formData.nodeType === 'llm'
                                    ? 'copilot -p "Refactor the authentication module..."'
                                    : 'git status && npm test'
                            }
                            rows={6}
                            className="w-full bg-slate-900 border border-slate-700 rounded px-3 py-2 text-white font-mono text-sm resize-none focus:outline-none focus:ring-2 focus:ring-blue-500"
                            data-testid="config-command-textarea"
                        />
                        
                        {/* Token Meter for LLM nodes */}
                        {formData.nodeType === 'llm' && (
                            <div className="mt-2">
                                <TokenMeter tokenCount={estimatedTokens} maxTokens={4000} />
                                <p className="text-xs text-purple-400 mt-1">
                                    💎 This node consumes tokens from your premium budget
                                </p>
                            </div>
                        )}
                        
                        {/* Zero-token indicator for shell nodes */}
                        {formData.nodeType === 'shell' && (
                            <p className="text-xs text-green-400 mt-2">
                                ⚡ Shell commands are free - no token consumption
                            </p>
                        )}
                    </div>
                </div>

                {/* Action Buttons */}
                <div className="flex gap-2 mt-4 pt-4 border-t border-slate-700">
                    <button
                        onClick={handleSave}
                        className={`flex-1 py-2 rounded font-medium transition-colors ${
                            formData.nodeType === 'llm'
                                ? 'bg-purple-600 hover:bg-purple-500 text-white'
                                : 'bg-blue-600 hover:bg-blue-500 text-white'
                        }`}
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

            {/* Premium Execution Confirmation Modal - Task 3.3 */}
            {showPremiumWarning && (
                <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50" data-testid="premium-confirm-modal">
                    <div className="bg-slate-800 rounded-lg p-6 max-w-md mx-4 border border-purple-500/50">
                        <div className="flex items-center gap-3 mb-4">
                            <div className="p-2 bg-purple-500/20 rounded-full">
                                <AlertTriangle className="text-purple-400" size={24} />
                            </div>
                            <h3 className="text-lg font-semibold text-white">Premium Resource Confirmation</h3>
                        </div>
                        
                        <p className="text-slate-300 mb-4">
                            This LLM node will consume tokens from your premium budget when executed.
                        </p>
                        
                        <div className="bg-slate-900 rounded p-3 mb-4">
                            <div className="flex justify-between text-sm mb-1">
                                <span className="text-slate-400">Estimated tokens:</span>
                                <span className="text-purple-400 font-mono">{estimatedTokens}</span>
                            </div>
                        </div>
                        
                        <div className="flex gap-3">
                            <button
                                onClick={() => setShowPremiumWarning(false)}
                                className="flex-1 bg-slate-700 hover:bg-slate-600 text-white py-2 rounded font-medium transition-colors"
                            >
                                Cancel
                            </button>
                            <button
                                onClick={confirmPremium}
                                className="flex-1 bg-purple-600 hover:bg-purple-500 text-white py-2 rounded font-medium transition-colors"
                                data-testid="confirm-premium-btn"
                            >
                                Confirm & Save
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </>
    );
};

export default NodeConfigPanel;
