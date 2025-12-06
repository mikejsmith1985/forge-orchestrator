import React, { useState, useEffect } from 'react';
import { FolderGit, FileCode, Filter, RefreshCw, X, ChevronDown, ChevronRight } from 'lucide-react';

/**
 * V2.1 Core Feature: Forge Context Navigator
 * 
 * Intelligent file picker with filters for efficient context selection.
 * Prevents token waste by helping users select only relevant files.
 * 
 * Filters:
 * - Uncommitted Changes: Files modified in git working tree
 * - Recent Files: Files modified in last 24h
 * - Source Code: Filter by file extension
 */

interface ContextNavigatorProps {
    selectedFiles: string[];
    onFilesSelected: (files: string[]) => void;
}

type FilterType = 'all' | 'uncommitted' | 'recent' | 'source';

interface FileEntry {
    path: string;
    status?: string; // For git status: M, A, D, ?
    type: 'file' | 'directory';
}

export const ContextNavigator: React.FC<ContextNavigatorProps> = ({
    selectedFiles,
    onFilesSelected
}) => {
    const [isExpanded, setIsExpanded] = useState(false);
    const [activeFilter, setActiveFilter] = useState<FilterType>('uncommitted');
    const [files, setFiles] = useState<FileEntry[]>([]);
    const [loading, setLoading] = useState(false);

    // Fetch files based on active filter
    useEffect(() => {
        const fetchFiles = async () => {
            setLoading(true);
            try {
                const response = await fetch(`/api/context/files?filter=${activeFilter}`);
                if (response.ok) {
                    const data = await response.json();
                    setFiles(data.files || []);
                } else {
                    // Fallback mock data for development
                    setFiles(getMockFiles(activeFilter));
                }
            } catch {
                // Fallback mock data
                setFiles(getMockFiles(activeFilter));
            } finally {
                setLoading(false);
            }
        };

        if (isExpanded) {
            fetchFiles();
        }
    }, [activeFilter, isExpanded]);

    const getMockFiles = (filter: FilterType): FileEntry[] => {
        switch (filter) {
            case 'uncommitted':
                return [
                    { path: 'src/components/Architect/ArchitectView.tsx', status: 'M', type: 'file' },
                    { path: 'src/components/Architect/BudgetMeter.tsx', status: 'A', type: 'file' },
                    { path: 'internal/server/ledger.go', status: 'M', type: 'file' },
                ];
            case 'recent':
                return [
                    { path: 'src/App.tsx', type: 'file' },
                    { path: 'src/components/Layout/Sidebar.tsx', type: 'file' },
                    { path: 'package.json', type: 'file' },
                ];
            case 'source':
                return [
                    { path: 'src/', type: 'directory' },
                    { path: 'internal/', type: 'directory' },
                    { path: 'tests/', type: 'directory' },
                ];
            default:
                return [];
        }
    };

    const toggleFile = (path: string) => {
        if (selectedFiles.includes(path)) {
            onFilesSelected(selectedFiles.filter(f => f !== path));
        } else {
            onFilesSelected([...selectedFiles, path]);
        }
    };

    const clearSelection = () => {
        onFilesSelected([]);
    };

    const getStatusBadge = (status?: string) => {
        if (!status) return null;
        const colors: Record<string, string> = {
            'M': 'bg-yellow-500/20 text-yellow-400 border-yellow-500/30',
            'A': 'bg-green-500/20 text-green-400 border-green-500/30',
            'D': 'bg-red-500/20 text-red-400 border-red-500/30',
            '?': 'bg-gray-500/20 text-gray-400 border-gray-500/30',
        };
        const labels: Record<string, string> = {
            'M': 'Modified',
            'A': 'Added',
            'D': 'Deleted',
            '?': 'Untracked',
        };
        return (
            <span className={`text-xs px-1.5 py-0.5 rounded border ${colors[status] || colors['?']}`}>
                {labels[status] || status}
            </span>
        );
    };

    const filters: { id: FilterType; label: string; icon: React.ReactNode }[] = [
        { id: 'uncommitted', label: 'Uncommitted', icon: <FolderGit className="w-3 h-3" /> },
        { id: 'recent', label: 'Recent', icon: <RefreshCw className="w-3 h-3" /> },
        { id: 'source', label: 'Source', icon: <FileCode className="w-3 h-3" /> },
    ];

    return (
        <div 
            className="bg-gray-900/50 rounded-lg border border-gray-800 overflow-hidden"
            data-testid="context-navigator"
        >
            {/* Header */}
            <button
                onClick={() => setIsExpanded(!isExpanded)}
                className="w-full flex items-center justify-between p-3 hover:bg-gray-800/50 transition-colors"
                data-testid="context-navigator-toggle"
            >
                <div className="flex items-center gap-2">
                    {isExpanded ? (
                        <ChevronDown className="w-4 h-4 text-gray-400" />
                    ) : (
                        <ChevronRight className="w-4 h-4 text-gray-400" />
                    )}
                    <FolderGit className="w-4 h-4 text-blue-400" />
                    <span className="text-sm font-medium text-gray-300">Context Navigator</span>
                    {selectedFiles.length > 0 && (
                        <span className="text-xs bg-blue-500/20 text-blue-400 px-2 py-0.5 rounded-full">
                            {selectedFiles.length} selected
                        </span>
                    )}
                </div>
                <Filter className="w-4 h-4 text-gray-500" />
            </button>

            {/* Expanded Content */}
            {isExpanded && (
                <div className="border-t border-gray-800">
                    {/* Filter Tabs */}
                    <div className="flex gap-1 p-2 bg-gray-900/30">
                        {filters.map(filter => (
                            <button
                                key={filter.id}
                                onClick={() => setActiveFilter(filter.id)}
                                className={`flex items-center gap-1.5 px-3 py-1.5 rounded text-xs font-medium transition-colors ${
                                    activeFilter === filter.id
                                        ? 'bg-blue-500/20 text-blue-400 border border-blue-500/30'
                                        : 'text-gray-400 hover:bg-gray-800 hover:text-gray-300'
                                }`}
                                data-testid={`filter-${filter.id}`}
                            >
                                {filter.icon}
                                {filter.label}
                            </button>
                        ))}
                    </div>

                    {/* File List */}
                    <div className="max-h-48 overflow-y-auto p-2">
                        {loading ? (
                            <div className="flex items-center justify-center py-4 text-gray-500">
                                <RefreshCw className="w-4 h-4 animate-spin mr-2" />
                                Loading...
                            </div>
                        ) : files.length === 0 ? (
                            <div className="text-center py-4 text-gray-500 text-sm">
                                No files found for this filter
                            </div>
                        ) : (
                            <div className="space-y-1">
                                {files.map(file => (
                                    <button
                                        key={file.path}
                                        onClick={() => toggleFile(file.path)}
                                        className={`w-full flex items-center justify-between p-2 rounded text-left text-sm transition-colors ${
                                            selectedFiles.includes(file.path)
                                                ? 'bg-blue-500/10 border border-blue-500/30'
                                                : 'hover:bg-gray-800/50'
                                        }`}
                                        data-testid={`file-${file.path.replace(/\//g, '-')}`}
                                    >
                                        <div className="flex items-center gap-2 min-w-0">
                                            <FileCode className="w-4 h-4 text-gray-500 flex-shrink-0" />
                                            <span className="text-gray-300 truncate font-mono text-xs">
                                                {file.path}
                                            </span>
                                        </div>
                                        {getStatusBadge(file.status)}
                                    </button>
                                ))}
                            </div>
                        )}
                    </div>

                    {/* Selected Files Footer */}
                    {selectedFiles.length > 0 && (
                        <div className="border-t border-gray-800 p-2 bg-gray-900/30">
                            <div className="flex items-center justify-between">
                                <span className="text-xs text-gray-500">
                                    {selectedFiles.length} file(s) will be included in context
                                </span>
                                <button
                                    onClick={clearSelection}
                                    className="flex items-center gap-1 text-xs text-red-400 hover:text-red-300"
                                    data-testid="clear-selection"
                                >
                                    <X className="w-3 h-3" />
                                    Clear
                                </button>
                            </div>
                        </div>
                    )}
                </div>
            )}
        </div>
    );
};
