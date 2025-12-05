import { useState, useEffect, useCallback } from 'react';
import { OptimizationCard, type Suggestion } from './OptimizationCard';
import { useWebSocket } from '../../hooks/useWebSocket';

// Educational Comment: Defining the shape of our data ensures type safety throughout the component.
interface LedgerEntry {
    id: string;
    timestamp: string;
    flowId: string;
    model: string;
    inputTokens: number;
    outputTokens: number;
    cost: number;
    latencyMs: number;
    status: 'success' | 'failed' | 'pending';
}

// API response shape (snake_case from Go backend)
interface ApiLedgerEntry {
    id: number;
    timestamp: string;
    flow_id: string;
    model_used: string;
    input_tokens: number;
    output_tokens: number;
    total_cost_usd: number;
    latency_ms: number;
    status: string;
}

// Transform API response to frontend format
const mapLedgerEntry = (entry: ApiLedgerEntry): LedgerEntry => ({
    id: String(entry.id),
    timestamp: entry.timestamp,
    flowId: entry.flow_id,
    model: entry.model_used,
    inputTokens: entry.input_tokens,
    outputTokens: entry.output_tokens,
    cost: entry.total_cost_usd,
    latencyMs: entry.latency_ms ?? 0,
    status: entry.status.toLowerCase() as 'success' | 'failed' | 'pending',
});

// Toast notification component
function Toast({ message, type, onClose }: { message: string; type: 'info' | 'success' | 'error'; onClose: () => void }) {
    useEffect(() => {
        const timer = setTimeout(onClose, 5000);
        return () => clearTimeout(timer);
    }, [onClose]);

    const bgColor = type === 'success' ? 'bg-green-500' : type === 'error' ? 'bg-red-500' : 'bg-blue-500';

    return (
        <div className={`${bgColor} text-white px-4 py-2 rounded-lg shadow-lg flex items-center gap-2`} data-testid="toast">
            <span>{message}</span>
            <button onClick={onClose} className="ml-2 hover:opacity-80">&times;</button>
        </div>
    );
}

export function LedgerView() {
    const [entries, setEntries] = useState<LedgerEntry[]>([]);
    const [optimizations, setOptimizations] = useState<Suggestion[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [toasts, setToasts] = useState<Array<{ id: number; message: string; type: 'info' | 'success' | 'error' }>>([]);

    // WebSocket connection for real-time updates
    const { lastMessage, isConnected } = useWebSocket('/ws');

    const addToast = useCallback((message: string, type: 'info' | 'success' | 'error') => {
        const id = Date.now();
        setToasts(prev => [...prev, { id, message, type }]);
    }, []);

    const removeToast = useCallback((id: number) => {
        setToasts(prev => prev.filter(t => t.id !== id));
    }, []);

    const fetchData = useCallback(async () => {
        try {
            const [ledgerRes, optimizationsRes] = await Promise.all([
                fetch('/api/ledger'),
                fetch('/api/ledger/optimizations')
            ]);

            if (!ledgerRes.ok) throw new Error('Failed to fetch ledger data');
            if (!optimizationsRes.ok) throw new Error('Failed to fetch optimizations');

            const ledgerData: ApiLedgerEntry[] = await ledgerRes.json();
            const optimizationsData = await optimizationsRes.json();

            setEntries(ledgerData.map(mapLedgerEntry));
            setOptimizations(optimizationsData);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An unknown error occurred');
        } finally {
            setLoading(false);
        }
    }, []);

    // Initial data fetch
    useEffect(() => {
        fetchData();
    }, [fetchData]);

    // Handle WebSocket messages
    useEffect(() => {
        if (!lastMessage) return;

        const { type, payload } = lastMessage;

        switch (type) {
            case 'FLOW_STARTED':
                addToast(`Flow ${payload.flowId} started`, 'info');
                break;
            case 'FLOW_COMPLETED':
                addToast(`Flow ${payload.flowId} completed`, 'success');
                fetchData(); // Refresh ledger
                break;
            case 'FLOW_FAILED':
                addToast(`Flow ${payload.flowId} failed: ${payload.error}`, 'error');
                fetchData(); // Refresh ledger to show failed entries
                break;
            case 'LEDGER_UPDATE':
                fetchData(); // Refresh ledger on any update
                break;
        }
    }, [lastMessage, addToast, fetchData]);

    // Educational Comment: This function handles the optimistic UI update pattern.
    // We update the local state immediately to reflect the change, while the API call happens.
    // In this specific case, we wait for the API call to succeed before updating the status to 'applied'.
    const handleApplyOptimization = async (id: number) => {
        const response = await fetch(`/api/ledger/optimizations/${id}/apply`, {
            method: 'POST',
        });

        if (!response.ok) {
            throw new Error('Failed to apply optimization');
        }

        // Update the local state to reflect the applied status
        setOptimizations(prev => prev.map(opt =>
            opt.id === id ? { ...opt, status: 'applied' as const } : opt
        ));

        // Optionally refetch ledger to show updated costs/entries if the optimization affects past entries immediately
        // or just to ensure consistency. For this requirement, updating the card status is the primary visual feedback.
    };

    if (loading) {
        return <div className="p-8 text-white">Loading ledger data...</div>;
    }

    if (error) {
        return <div className="p-8 text-red-400">Error: {error}</div>;
    }

    return (
        <div className="p-8 h-full overflow-auto" data-testid="ledger-view">
            {/* Toast notifications */}
            <div className="fixed top-4 right-4 z-50 flex flex-col gap-2" data-testid="toast-container">
                {toasts.map(toast => (
                    <Toast
                        key={toast.id}
                        message={toast.message}
                        type={toast.type}
                        onClose={() => removeToast(toast.id)}
                    />
                ))}
            </div>

            {/* WebSocket connection indicator */}
            <div className="flex items-center gap-2 mb-4">
                <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`} />
                <span className="text-xs text-gray-400">
                    {isConnected ? 'Live updates enabled' : 'Connecting...'}
                </span>
            </div>

            <h2 className="text-2xl font-bold mb-6 text-white">Token Ledger</h2>

            {/* Optimization Suggestions Section */}
            <div className="mb-8">
                <h3 className="text-xl font-semibold mb-4 text-white">Optimizations</h3>
                {optimizations.length > 0 ? (
                    optimizations.map(opt => (
                        <OptimizationCard
                            key={opt.id}
                            suggestion={opt}
                            onApply={handleApplyOptimization}
                        />
                    ))
                ) : (
                    <div className="text-gray-500 italic">No optimization suggestions yet</div>
                )}
            </div>

            <div className="bg-gray-900/50 rounded-lg border border-white/10 overflow-hidden">
                <table className="w-full text-left text-sm text-gray-400" data-testid="ledger-table">
                    <thead className="bg-white/5 text-gray-200 uppercase font-medium">
                        <tr>
                            <th className="px-6 py-4">Timestamp</th>
                            <th className="px-6 py-4">Flow ID</th>
                            <th className="px-6 py-4">Model</th>
                            <th className="px-6 py-4 text-right">Input Tokens</th>
                            <th className="px-6 py-4 text-right">Output Tokens</th>
                            <th className="px-6 py-4 text-right">Latency</th>
                            <th className="px-6 py-4 text-right">Cost ($)</th>
                            <th className="px-6 py-4">Status</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-white/5">
                        {entries.map((entry) => (
                            <tr key={entry.id} className="hover:bg-white/5 transition-colors">
                                <td className="px-6 py-4 whitespace-nowrap">
                                    {new Date(entry.timestamp).toLocaleString()}
                                </td>
                                <td className="px-6 py-4 font-mono text-xs">{entry.flowId}</td>
                                <td className="px-6 py-4">
                                    <span className="px-2 py-1 rounded-full bg-blue-500/10 text-blue-400 text-xs border border-blue-500/20">
                                        {entry.model}
                                    </span>
                                </td>
                                <td className="px-6 py-4 text-right font-mono">{entry.inputTokens}</td>
                                <td className="px-6 py-4 text-right font-mono">{entry.outputTokens}</td>
                                <td className="px-6 py-4 text-right font-mono">
                                    <span className={
                                        entry.latencyMs < 1000 ? 'text-green-400' :
                                        entry.latencyMs < 5000 ? 'text-yellow-400' : 'text-red-400'
                                    }>
                                        {entry.latencyMs < 1000 
                                            ? `${entry.latencyMs}ms` 
                                            : `${(entry.latencyMs / 1000).toFixed(2)}s`}
                                    </span>
                                </td>
                                <td className="px-6 py-4 text-right font-mono text-green-400">
                                    ${entry.cost.toFixed(4)}
                                </td>
                                <td className="px-6 py-4">
                                    <span className={`px-2 py-1 rounded-full text-xs border ${entry.status === 'success'
                                        ? 'bg-green-500/10 text-green-400 border-green-500/20'
                                        : 'bg-red-500/10 text-red-400 border-red-500/20'
                                        }`}>
                                        {entry.status}
                                    </span>
                                </td>
                            </tr>
                        ))}
                        {entries.length === 0 && (
                            <tr>
                                <td colSpan={8} className="px-6 py-8 text-center text-gray-500">
                                    No ledger entries found
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
