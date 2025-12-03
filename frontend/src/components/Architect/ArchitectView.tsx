import React, { useState, useEffect } from 'react';
import { TokenMeter } from './TokenMeter';

/**
 * ArchitectView provides an interface for the "Forge Architect" input.
 * It includes a textarea for "Brain Dump" and a live TokenMeter.
 * Token count is approximated as 4 characters per token.
 */
export const ArchitectView: React.FC = () => {
    const [input, setInput] = useState('');
    const [tokenCount, setTokenCount] = useState(0);

    useEffect(() => {
        // Approximate token count: 4 characters = 1 token
        const count = Math.ceil(input.length / 4);
        setTokenCount(count);
    }, [input]);

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

                <TokenMeter tokenCount={tokenCount} />
            </div>
        </div>
    );
};
