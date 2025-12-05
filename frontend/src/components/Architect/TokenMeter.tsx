import React from 'react';

interface TokenMeterProps {
    tokenCount: number;
    maxTokens?: number;
    method?: string;
    provider?: string;
}

/**
 * TokenMeter component visualizes the current token usage against a maximum limit.
 * It displays a progress bar that changes color based on usage percentage:
 * - Green: < 75%
 * - Yellow: 75% - 90%
 * - Red: > 90%
 */
export const TokenMeter: React.FC<TokenMeterProps> = ({
    tokenCount,
    maxTokens = 8000,
    method,
    provider
}) => {
    const percentage = Math.min((tokenCount / maxTokens) * 100, 100);

    let colorClass = 'bg-green-500';
    if (percentage > 90) {
        colorClass = 'bg-red-500';
    } else if (percentage >= 75) {
        colorClass = 'bg-yellow-500';
    }

    // Format method badge
    const methodBadge = method ? (
        <span className={`px-1.5 py-0.5 text-xs rounded ${
            method === 'tiktoken' ? 'bg-green-900 text-green-300' : 'bg-gray-700 text-gray-300'
        }`}>
            {method}
        </span>
    ) : null;

    // Format provider info
    const providerInfo = provider ? (
        <span className="text-gray-500 text-xs ml-1">
            ({provider})
        </span>
    ) : null;

    return (
        <div className="w-full space-y-2" data-testid="token-meter">
            <div className="flex justify-between text-sm text-gray-400">
                <span className="flex items-center gap-2">
                    Context Usage
                    {methodBadge}
                    {providerInfo}
                </span>
                <span>{tokenCount.toLocaleString()} / {maxTokens.toLocaleString()} tokens</span>
            </div>
            <div className="w-full h-2 bg-gray-800 rounded-full overflow-hidden">
                <div
                    className={`h-full transition-all duration-300 ease-out ${colorClass}`}
                    style={{ width: `${percentage}%` }}
                />
            </div>
        </div>
    );
};
