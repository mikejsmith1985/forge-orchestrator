import React, { memo } from 'react';
import { Handle, Position, type NodeProps } from 'reactflow';

export interface AgentNodeData {
    label: string;
    role?: 'Architect' | 'Implementation' | 'Test' | 'Optimizer';
    provider?: string;
    prompt?: string;
    retryConfig?: {
        maxRetries: number;
        backoffMs: number;
    };
}

const AgentNode: React.FC<NodeProps<AgentNodeData>> = ({ data, selected }) => {
    const isConfigured = Boolean(data.role && data.provider && data.prompt);

    return (
        <div
            className={`
                p-4 rounded-lg border-2 min-w-[150px]
                ${selected ? 'border-blue-500' : 'border-slate-600'}
                ${isConfigured ? 'bg-slate-800' : 'bg-yellow-900/20'}
            `}
            data-testid="agent-node"
        >
            <Handle type="target" position={Position.Top} className="w-3 h-3" />

            <div className="flex items-center gap-2">
                <div
                    className={`w-3 h-3 rounded-full ${
                        isConfigured ? 'bg-green-500' : 'bg-yellow-500'
                    }`}
                    data-testid={isConfigured ? 'node-configured' : 'node-unconfigured'}
                />
                <span className="text-white font-medium">{data.label}</span>
            </div>

            {data.role && (
                <div className="text-xs text-slate-400 mt-1">{data.role}</div>
            )}

            {!isConfigured && (
                <div className="text-xs text-yellow-400 mt-2" data-testid="node-warning">
                    ⚠ Not configured
                </div>
            )}

            <Handle type="source" position={Position.Bottom} className="w-3 h-3" />
        </div>
    );
};

export default memo(AgentNode);
