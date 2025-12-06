import React, { useState, useCallback } from 'react';
import { Eye, Code2, GitBranch, FileText, ChevronRight, Lightbulb, Copy, Check } from 'lucide-react';

/**
 * V2.1 Core Feature: Forge Vision
 * 
 * Visualizes terminal output to make it more actionable.
 * Features:
 * - Git status visualization as clickable dashboard
 * - Explain Code button for complex output
 * - Copy to context for AI prompts
 */

interface ForgeVisionProps {
    output: string;
    onExplainCode?: (code: string) => void;
    onAddToContext?: (content: string) => void;
}

interface ParsedGitStatus {
    staged: string[];
    modified: string[];
    untracked: string[];
    branch: string;
}

export const ForgeVision: React.FC<ForgeVisionProps> = ({
    output,
    onExplainCode,
    onAddToContext
}) => {
    const [copied, setCopied] = useState(false);
    const [isExpanded, setIsExpanded] = useState(true);

    // Detect and parse git status output
    const parseGitStatus = useCallback((text: string): ParsedGitStatus | null => {
        if (!text.includes('On branch') && !text.includes('Changes') && !text.includes('Untracked')) {
            return null;
        }

        const lines = text.split('\n');
        const result: ParsedGitStatus = {
            staged: [],
            modified: [],
            untracked: [],
            branch: ''
        };

        let section = '';
        for (const line of lines) {
            if (line.startsWith('On branch ')) {
                result.branch = line.replace('On branch ', '').trim();
            } else if (line.includes('Changes to be committed')) {
                section = 'staged';
            } else if (line.includes('Changes not staged')) {
                section = 'modified';
            } else if (line.includes('Untracked files')) {
                section = 'untracked';
            } else if (line.trim().startsWith('modified:')) {
                const file = line.replace(/.*modified:\s*/, '').trim();
                if (section === 'staged') {
                    result.staged.push(file);
                } else {
                    result.modified.push(file);
                }
            } else if (line.trim().startsWith('new file:')) {
                const file = line.replace(/.*new file:\s*/, '').trim();
                result.staged.push(file);
            } else if (section === 'untracked' && line.trim() && !line.includes('(') && !line.includes('git add')) {
                result.untracked.push(line.trim());
            }
        }

        return result;
    }, []);

    // Detect if output is code/error that needs explanation
    const isCodeOrError = useCallback((text: string): boolean => {
        const codeIndicators = [
            'error:', 'Error:', 'ERROR',
            'warning:', 'Warning:', 'WARN',
            'exception', 'Exception',
            'stack trace', 'Traceback',
            'undefined', 'null reference',
            'at line', 'syntax error'
        ];
        return codeIndicators.some(indicator => text.includes(indicator));
    }, []);

    const handleCopy = async () => {
        await navigator.clipboard.writeText(output);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    const handleExplain = () => {
        onExplainCode?.(output);
    };

    const handleAddToContext = () => {
        onAddToContext?.(output);
    };

    const gitStatus = parseGitStatus(output);
    const needsExplanation = isCodeOrError(output);

    if (!output.trim()) {
        return null;
    }

    return (
        <div 
            className="bg-slate-800/50 rounded-lg border border-slate-700 overflow-hidden"
            data-testid="forge-vision"
        >
            {/* Header */}
            <button
                onClick={() => setIsExpanded(!isExpanded)}
                className="w-full flex items-center justify-between p-3 hover:bg-slate-700/50 transition-colors"
            >
                <div className="flex items-center gap-2">
                    <Eye className="w-4 h-4 text-purple-400" />
                    <span className="text-sm font-medium text-slate-300">Forge Vision</span>
                    {gitStatus && (
                        <span className="text-xs bg-green-500/20 text-green-400 px-2 py-0.5 rounded-full">
                            Git Status
                        </span>
                    )}
                    {needsExplanation && (
                        <span className="text-xs bg-yellow-500/20 text-yellow-400 px-2 py-0.5 rounded-full">
                            Needs Attention
                        </span>
                    )}
                </div>
                <ChevronRight className={`w-4 h-4 text-slate-500 transition-transform ${isExpanded ? 'rotate-90' : ''}`} />
            </button>

            {/* Expanded Content */}
            {isExpanded && (
                <div className="border-t border-slate-700 p-4">
                    {/* Git Status Visualization */}
                    {gitStatus && (
                        <div className="space-y-3 mb-4" data-testid="git-status-viz">
                            {/* Branch info */}
                            <div className="flex items-center gap-2 text-sm">
                                <GitBranch className="w-4 h-4 text-purple-400" />
                                <span className="text-slate-400">Branch:</span>
                                <span className="text-white font-mono">{gitStatus.branch}</span>
                            </div>

                            {/* Staged files */}
                            {gitStatus.staged.length > 0 && (
                                <div className="space-y-1">
                                    <div className="text-xs text-green-400 font-medium">
                                        ✓ Staged ({gitStatus.staged.length})
                                    </div>
                                    <div className="pl-4 space-y-1">
                                        {gitStatus.staged.map((file, i) => (
                                            <button
                                                key={i}
                                                className="flex items-center gap-2 text-sm text-green-300 hover:bg-green-500/10 px-2 py-1 rounded w-full text-left"
                                                onClick={() => onAddToContext?.(file)}
                                                data-testid={`staged-file-${i}`}
                                            >
                                                <FileText className="w-3 h-3" />
                                                <span className="font-mono text-xs">{file}</span>
                                            </button>
                                        ))}
                                    </div>
                                </div>
                            )}

                            {/* Modified files */}
                            {gitStatus.modified.length > 0 && (
                                <div className="space-y-1">
                                    <div className="text-xs text-yellow-400 font-medium">
                                        ● Modified ({gitStatus.modified.length})
                                    </div>
                                    <div className="pl-4 space-y-1">
                                        {gitStatus.modified.map((file, i) => (
                                            <button
                                                key={i}
                                                className="flex items-center gap-2 text-sm text-yellow-300 hover:bg-yellow-500/10 px-2 py-1 rounded w-full text-left"
                                                onClick={() => onAddToContext?.(file)}
                                                data-testid={`modified-file-${i}`}
                                            >
                                                <FileText className="w-3 h-3" />
                                                <span className="font-mono text-xs">{file}</span>
                                            </button>
                                        ))}
                                    </div>
                                </div>
                            )}

                            {/* Untracked files */}
                            {gitStatus.untracked.length > 0 && (
                                <div className="space-y-1">
                                    <div className="text-xs text-slate-400 font-medium">
                                        ? Untracked ({gitStatus.untracked.length})
                                    </div>
                                    <div className="pl-4 space-y-1">
                                        {gitStatus.untracked.map((file, i) => (
                                            <button
                                                key={i}
                                                className="flex items-center gap-2 text-sm text-slate-400 hover:bg-slate-500/10 px-2 py-1 rounded w-full text-left"
                                                onClick={() => onAddToContext?.(file)}
                                                data-testid={`untracked-file-${i}`}
                                            >
                                                <FileText className="w-3 h-3" />
                                                <span className="font-mono text-xs">{file}</span>
                                            </button>
                                        ))}
                                    </div>
                                </div>
                            )}
                        </div>
                    )}

                    {/* Raw output preview */}
                    {!gitStatus && (
                        <div className="bg-slate-900 rounded p-3 mb-3 max-h-32 overflow-y-auto">
                            <pre className="text-xs text-slate-400 font-mono whitespace-pre-wrap">
                                {output.slice(0, 500)}
                                {output.length > 500 && '...'}
                            </pre>
                        </div>
                    )}

                    {/* Action buttons */}
                    <div className="flex gap-2 flex-wrap">
                        <button
                            onClick={handleCopy}
                            className="flex items-center gap-1.5 px-3 py-1.5 bg-slate-700 hover:bg-slate-600 rounded text-xs text-slate-300 transition-colors"
                            data-testid="copy-output"
                        >
                            {copied ? <Check className="w-3 h-3" /> : <Copy className="w-3 h-3" />}
                            {copied ? 'Copied!' : 'Copy'}
                        </button>

                        {needsExplanation && onExplainCode && (
                            <button
                                onClick={handleExplain}
                                className="flex items-center gap-1.5 px-3 py-1.5 bg-purple-600 hover:bg-purple-500 rounded text-xs text-white transition-colors"
                                data-testid="explain-code"
                            >
                                <Lightbulb className="w-3 h-3" />
                                Explain This
                            </button>
                        )}

                        {onAddToContext && (
                            <button
                                onClick={handleAddToContext}
                                className="flex items-center gap-1.5 px-3 py-1.5 bg-blue-600 hover:bg-blue-500 rounded text-xs text-white transition-colors"
                                data-testid="add-to-context"
                            >
                                <Code2 className="w-3 h-3" />
                                Add to Context
                            </button>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
};

export default ForgeVision;
