import React, { memo } from 'react';
import { Handle, Position, type NodeProps } from 'reactflow';
import { Terminal, Brain, Zap } from 'lucide-react';

/**
 * AgentNodeData - Redesigned for V2.1 Remediation Plan Task 3.1
 * 
 * CHANGE: Removed confusing Agent Role and Provider dropdowns.
 * The node now accepts a simple command/prompt string and distinguishes
 * between two node types: Shell Command (zero-token) and LLM Prompt (premium).
 */
export interface AgentNodeData {
    label: string;
    /** The command or prompt string to execute */
    command?: string;
    /** Node type: 'shell' for local commands, 'llm' for LLM prompts */
    nodeType?: 'shell' | 'llm';
    /** Whether this node has been confirmed for premium execution */
    premiumConfirmed?: boolean;
}

const AgentNode: React.FC<NodeProps<AgentNodeData>> = ({ data, selected }) => {
    const isConfigured = Boolean(data.command);
    const isLLMNode = data.nodeType === 'llm';
    const isShellNode = data.nodeType === 'shell';

    // Get icon based on node type
    const NodeIcon = isLLMNode ? Brain : isShellNode ? Terminal : Zap;

    return (
        <div
            className={`
                p-4 rounded-lg border-2 min-w-[180px] max-w-[250px]
                ${selected ? 'border-blue-500' : 'border-slate-600'}
                ${isConfigured 
                    ? isLLMNode 
                        ? 'bg-purple-900/30 border-purple-500/50' 
                        : 'bg-slate-800' 
                    : 'bg-yellow-900/20'}
            `}
            data-testid="agent-node"
        >
            <Handle type="target" position={Position.Top} className="w-3 h-3" />

            <div className="flex items-center gap-2">
                <NodeIcon 
                    size={16} 
                    className={isLLMNode ? 'text-purple-400' : 'text-blue-400'} 
                />
                <span className="text-white font-medium truncate">{data.label}</span>
            </div>

            {/* Node Type Badge */}
            {data.nodeType && (
                <div className={`text-xs mt-2 px-2 py-0.5 rounded-full inline-block ${
                    isLLMNode 
                        ? 'bg-purple-500/20 text-purple-300 border border-purple-500/30' 
                        : 'bg-green-500/20 text-green-300 border border-green-500/30'
                }`}>
                    {isLLMNode ? '💎 Premium' : '⚡ Zero-Token'}
                </div>
            )}

            {/* Command preview */}
            {data.command && (
                <div className="text-xs text-slate-400 mt-2 font-mono truncate" title={data.command}>
                    {data.command.length > 30 ? data.command.substring(0, 30) + '...' : data.command}
                </div>
            )}

            {!isConfigured && (
                <div className="text-xs text-yellow-400 mt-2" data-testid="node-warning">
                    ⚠ Click to configure
                </div>
            )}

            <Handle type="source" position={Position.Bottom} className="w-3 h-3" />
        </div>
    );
};

export default memo(AgentNode);
