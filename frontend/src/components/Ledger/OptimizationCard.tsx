import { useState } from 'react';

// Educational Comment: Defining the interface for the suggestion object ensures we know what data to expect.
export interface Suggestion {
    id: string;
    title: string;
    description: string;
    estimatedSavings: number;
    status: 'pending' | 'applied';
}

interface OptimizationCardProps {
    suggestion: Suggestion;
    onApply: (id: string) => Promise<void>;
}

export function OptimizationCard({ suggestion, onApply }: OptimizationCardProps) {
    // Educational Comment: Local state to handle the loading status of the apply action.
    // This provides immediate feedback to the user while the async operation completes.
    const [isApplying, setIsApplying] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleApply = async () => {
        setIsApplying(true);
        setError(null);
        try {
            // Educational Comment: We await the parent's onApply handler.
            // This allows the parent to manage the actual data update and API call.
            await onApply(suggestion.id);
            // Note: We don't set isApplying(false) here if successful, because the parent
            // will likely update the suggestion prop to 'applied', which disables the button.
            // However, strictly speaking, we might want to reset it if the parent doesn't unmount us.
            // In this case, the button becomes disabled/changed based on the new 'status' prop.
        } catch (err) {
            setError('Failed to apply');
            setIsApplying(false);
        }
    };

    const isApplied = suggestion.status === 'applied';

    return (
        <div
            className="bg-gray-900/50 border border-white/10 rounded-lg p-4 mb-4 flex justify-between items-center"
            data-testid="optimization-card"
        >
            <div>
                <h3 className="text-lg font-semibold text-white">{suggestion.title}</h3>
                <p className="text-gray-400 text-sm">{suggestion.description}</p>
                <p className="text-green-400 text-sm mt-1 font-mono">
                    Estimated Savings: ${suggestion.estimatedSavings.toFixed(4)}
                </p>
                {error && <p className="text-red-400 text-xs mt-1">{error}</p>}
            </div>
            <button
                onClick={handleApply}
                disabled={isApplied || isApplying}
                className={`px-4 py-2 rounded text-sm font-medium transition-colors ${isApplied
                        ? 'bg-green-500/20 text-green-400 border border-green-500/20 cursor-default'
                        : 'bg-blue-600 hover:bg-blue-500 text-white disabled:opacity-50 disabled:cursor-not-allowed'
                    }`}
                data-testid="apply-btn"
            >
                {isApplying ? 'Applying...' : isApplied ? 'Applied' : 'Apply'}
            </button>
        </div>
    );
}
