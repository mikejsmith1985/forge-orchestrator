import { useState } from 'react';
import { AlertTriangle, Check, X } from 'lucide-react';

export interface Suggestion {
    id: number;
    type: string;
    title: string;
    description: string;
    estimated_savings: number;
    savings_unit: string;
    target_flow_id: string;
    apply_action: string;
    status: 'pending' | 'applied';
}

interface ApplyAction {
    action: string;
    from_model?: string;
    to_model?: string;
    flow_id?: string;
}

interface OptimizationCardProps {
    suggestion: Suggestion;
    onApply: (id: number) => Promise<void>;
}

export function OptimizationCard({ suggestion, onApply }: OptimizationCardProps) {
    const [isApplying, setIsApplying] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [showConfirmModal, setShowConfirmModal] = useState(false);

    const parseApplyAction = (): ApplyAction | null => {
        try {
            return JSON.parse(suggestion.apply_action);
        } catch {
            return null;
        }
    };

    const getChangeDescription = (): string => {
        const action = parseApplyAction();
        if (!action) return 'Apply this optimization';

        switch (action.action) {
            case 'switch_model':
                return `Change flow ${action.flow_id} from ${action.from_model} to ${action.to_model}`;
            case 'optimize_prompt':
                return `Flag flow ${action.flow_id} for prompt optimization review`;
            case 'implement_retry':
                return `Add exponential backoff retry logic to flow ${action.flow_id}`;
            default:
                return 'Apply this optimization';
        }
    };

    const handleApplyClick = () => {
        setShowConfirmModal(true);
    };

    const handleConfirmApply = async () => {
        setIsApplying(true);
        setError(null);
        setShowConfirmModal(false);
        try {
            await onApply(suggestion.id);
            setIsApplying(false);
        } catch {
            setError('Failed to apply');
            setIsApplying(false);
        }
    };

    const isApplied = suggestion.status === 'applied';

    return (
        <>
            <div
                className="bg-gray-900/50 border border-white/10 rounded-lg p-4 mb-4 flex justify-between items-center"
                data-testid="optimization-card"
            >
                <div className="flex-1">
                    <div className="flex items-center gap-2">
                        <h3 className="text-lg font-semibold text-white">{suggestion.title}</h3>
                        {suggestion.type === 'model_switch' && (
                            <span className="px-2 py-0.5 bg-blue-500/20 text-blue-400 text-xs rounded">Model Switch</span>
                        )}
                        {suggestion.type === 'prompt_optimization' && (
                            <span className="px-2 py-0.5 bg-purple-500/20 text-purple-400 text-xs rounded">Prompt</span>
                        )}
                        {suggestion.type === 'retry_strategy' && (
                            <span className="px-2 py-0.5 bg-orange-500/20 text-orange-400 text-xs rounded">Retry</span>
                        )}
                    </div>
                    <p className="text-gray-400 text-sm mt-1">{suggestion.description}</p>
                    <p className="text-green-400 text-sm mt-2 font-mono">
                        Estimated Savings: {suggestion.savings_unit === 'USD' ? '$' : ''}{suggestion.estimated_savings.toFixed(suggestion.savings_unit === 'USD' ? 4 : 0)} {suggestion.savings_unit === 'tokens' ? 'tokens' : ''}
                    </p>
                    {error && <p className="text-red-400 text-xs mt-1">{error}</p>}
                </div>
                <button
                    onClick={handleApplyClick}
                    disabled={isApplied || isApplying}
                    className={`px-4 py-2 rounded text-sm font-medium transition-colors ml-4 ${isApplied
                        ? 'bg-green-500/20 text-green-400 border border-green-500/20 cursor-default'
                        : 'bg-blue-600 hover:bg-blue-500 text-white disabled:opacity-50 disabled:cursor-not-allowed'
                        }`}
                    data-testid="apply-btn"
                >
                    {isApplying ? 'Applying...' : isApplied ? 'Applied' : 'Apply'}
                </button>
            </div>

            {/* Confirmation Modal */}
            {showConfirmModal && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
                    <div className="bg-slate-800 border border-slate-700 rounded-xl p-6 max-w-md w-full mx-4 shadow-xl">
                        <div className="flex items-center gap-3 mb-4">
                            <div className="p-2 bg-yellow-500/20 rounded-lg">
                                <AlertTriangle className="text-yellow-400" size={24} />
                            </div>
                            <h3 className="text-lg font-semibold text-white">Confirm Changes</h3>
                        </div>
                        
                        <p className="text-gray-300 mb-4">
                            This will make the following change:
                        </p>
                        
                        <div className="bg-slate-700/50 border border-slate-600 rounded-lg p-3 mb-6">
                            <p className="text-white font-mono text-sm">{getChangeDescription()}</p>
                        </div>
                        
                        <p className="text-gray-400 text-sm mb-6">
                            Are you sure you want to continue? This action will modify your flow configuration.
                        </p>
                        
                        <div className="flex gap-3 justify-end">
                            <button
                                onClick={() => setShowConfirmModal(false)}
                                className="flex items-center gap-2 px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors"
                            >
                                <X size={18} />
                                Cancel
                            </button>
                            <button
                                onClick={handleConfirmApply}
                                className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-lg transition-colors"
                            >
                                <Check size={18} />
                                Apply Changes
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </>
    );
}
