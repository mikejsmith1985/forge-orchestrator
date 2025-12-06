import React from 'react';
import { Wallet, Zap } from 'lucide-react';

/**
 * Task 4.2: Dynamic Budget Meter
 * 
 * Displays the remaining budget for the selected model.
 * Shows "X Prompts Remaining" for transparent resource gating.
 */

interface BudgetMeterProps {
    remainingPrompts: number;
    remainingBudget: number;
    totalBudget: number;
    costUnit: 'TOKEN' | 'PROMPT';
    model: string;
}

export const BudgetMeter: React.FC<BudgetMeterProps> = ({
    remainingPrompts,
    remainingBudget,
    totalBudget,
    costUnit,
    model
}) => {
    const percentage = (remainingBudget / totalBudget) * 100;
    
    // Color based on remaining budget
    let colorClass = 'text-green-400';
    let bgColorClass = 'bg-green-500';
    if (percentage < 25) {
        colorClass = 'text-red-400';
        bgColorClass = 'bg-red-500';
    } else if (percentage < 50) {
        colorClass = 'text-yellow-400';
        bgColorClass = 'bg-yellow-500';
    }

    const unitLabel = costUnit === 'PROMPT' ? 'Prompts' : 'Prompts';
    const unitIcon = costUnit === 'PROMPT' ? '💬' : '🪙';

    return (
        <div 
            className="bg-gray-900/70 rounded-lg border border-gray-700 p-4"
            data-testid="budget-meter"
        >
            <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-2">
                    <Wallet className="w-4 h-4 text-blue-400" />
                    <span className="text-sm font-medium text-gray-300">Daily Budget</span>
                    <span className="text-xs text-gray-500 bg-gray-800 px-2 py-0.5 rounded">
                        {model}
                    </span>
                </div>
                <div className="flex items-center gap-1">
                    <Zap className={`w-4 h-4 ${colorClass}`} />
                    <span className={`text-lg font-bold ${colorClass}`} data-testid="remaining-prompts">
                        {remainingPrompts.toLocaleString()}
                    </span>
                    <span className="text-sm text-gray-400">{unitLabel} Remaining</span>
                </div>
            </div>

            {/* Progress bar */}
            <div className="w-full h-2 bg-gray-800 rounded-full overflow-hidden">
                <div
                    className={`h-full transition-all duration-500 ease-out ${bgColorClass}`}
                    style={{ width: `${percentage}%` }}
                    data-testid="budget-meter-bar"
                />
            </div>

            {/* Budget details */}
            <div className="flex justify-between mt-2 text-xs text-gray-500">
                <span>{unitIcon} ${remainingBudget.toFixed(2)} remaining</span>
                <span>${totalBudget.toFixed(2)} daily limit</span>
            </div>
        </div>
    );
};
