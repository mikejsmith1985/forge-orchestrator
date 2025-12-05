import React, { useState, useEffect, useCallback } from 'react';
import { TokenMeter } from './TokenMeter';

interface TokenEstimate {
    count: number;
    method: string;
    provider: string;
    model?: string;
}

/**
 * ArchitectView provides an interface for the "Forge Architect" input.
 * It includes a textarea for "Brain Dump" and a live TokenMeter.
 * Token estimation uses the backend API with tiktoken for accuracy.
 */
export const ArchitectView: React.FC = () => {
    const [input, setInput] = useState('');
    const [tokenEstimate, setTokenEstimate] = useState<TokenEstimate>({
        count: 0,
        method: 'heuristic',
        provider: 'openai'
    });
    const [provider] = useState('openai'); // Could be made configurable

    // Debounced token estimation
    const estimateTokens = useCallback(async (text: string) => {
        if (!text.trim()) {
            setTokenEstimate({ count: 0, method: 'heuristic', provider });
            return;
        }

        try {
            const response = await fetch('/api/tokens/estimate', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ text, provider })
            });

            if (response.ok) {
                const data = await response.json();
                setTokenEstimate(data);
            } else {
                // Fallback to local estimation
                setTokenEstimate({
                    count: Math.ceil(text.length / 4),
                    method: 'fallback',
                    provider
                });
            }
        } catch {
            // Fallback to local estimation on network error
            setTokenEstimate({
                count: Math.ceil(text.length / 4),
                method: 'fallback',
                provider
            });
        }
    }, [provider]);

    // Debounce API calls
    useEffect(() => {
        const timeoutId = setTimeout(() => {
            estimateTokens(input);
        }, 300);

        return () => clearTimeout(timeoutId);
    }, [input, estimateTokens]);

    return (
        <div className="flex flex-col h-full p-6 space-y-6 max-w-4xl mx-auto w-full">
            <div className="space-y-2">
                <h1 className="text-2xl font-bold text-white">Forge Architect</h1>
                <p className="text-gray-400">
                    Describe your vision, requirements, or tasks. The architect will analyze and plan the implementation.
                </p>
            </div>

            <div className="flex-1 flex flex-col space-y-4">
                <div className="flex-1 bg-gray-900/50 rounded-lg border border-gray-800 p-4 focus-within:border-blue-500/50 transition-colors">
                    <textarea
                        data-testid="architect-input"
                        className="w-full h-full bg-transparent text-gray-200 placeholder-gray-600 resize-none focus:outline-none font-mono text-sm leading-relaxed"
                        placeholder="Start typing your requirements here..."
                        value={input}
                        onChange={(e) => setInput(e.target.value)}
                        spellCheck={false}
                    />
                </div>

                <TokenMeter 
                    tokenCount={tokenEstimate.count}
                    method={tokenEstimate.method}
                    provider={tokenEstimate.provider}
                />
            </div>
        </div>
    );
};
