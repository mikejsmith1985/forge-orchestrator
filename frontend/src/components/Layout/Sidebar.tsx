import { useState } from 'react';
import { LayoutDashboard, GitGraph, Settings, Menu, X, BrainCircuit, Workflow, Download, MessageSquare, TerminalSquare } from 'lucide-react';
import { cn } from '../../lib/utils';
import { useNavigate, useLocation } from 'react-router-dom';

interface SidebarProps {
    currentVersion?: string;
    hasUpdate?: boolean;
    onUpdateClick?: () => void;
    onFeedbackClick?: () => void;
}

export function Sidebar({ currentVersion, hasUpdate, onUpdateClick, onFeedbackClick }: SidebarProps) {
    const [isOpen, setIsOpen] = useState(false);
    const navigate = useNavigate();
    const location = useLocation();

    const navItems = [
        { icon: TerminalSquare, label: 'Terminal', path: '/terminal' },
        { icon: BrainCircuit, label: 'Architect', path: '/architect' },
        { icon: LayoutDashboard, label: 'Dashboard', path: '/ledger' },
        { icon: GitGraph, label: 'Commands', path: '/commands' },
        { icon: Workflow, label: 'Flows', path: '/flows' },
        { icon: Settings, label: 'Settings', path: '/settings' },
    ];

    const isActive = (path: string) => {
        if (path === '/flows' && location.pathname.startsWith('/flows')) return true;
        return location.pathname === path;
    };

    return (
        <>
            {/* Mobile Menu Button */}
            <button
                data-testid="mobile-menu-btn"
                className="fixed top-4 left-4 z-50 p-2 rounded-md bg-gray-800 text-white md:hidden"
                onClick={() => setIsOpen(!isOpen)}
            >
                {isOpen ? <X size={24} /> : <Menu size={24} />}
            </button>

            {/* Sidebar Container */}
            <div
                data-testid="sidebar"
                className={cn(
                    "fixed inset-y-0 left-0 z-40 w-64 bg-gray-900/95 backdrop-blur-xl border-r border-white/10 text-white transition-transform duration-300 ease-in-out md:translate-x-0 md:static md:h-screen md:visible",
                    isOpen ? "translate-x-0" : "-translate-x-full invisible"
                )}
            >
                <div className="flex flex-col h-full">
                    {/* Header */}
                    <div className="flex items-center justify-center h-16 border-b border-white/10">
                        <h1 className="text-xl font-bold bg-gradient-to-r from-blue-400 to-purple-500 bg-clip-text text-transparent">
                            Forge Orchestrator
                        </h1>
                    </div>

                    {/* Navigation */}
                    <nav className="flex-1 px-4 py-6 space-y-2">
                        {navItems.map((item) => (
                            <button
                                key={item.label}
                                onClick={() => {
                                    navigate(item.path);
                                    setIsOpen(false);
                                }}
                                className={cn(
                                    "flex items-center w-full px-4 py-3 rounded-lg transition-colors group text-left",
                                    isActive(item.path)
                                        ? "bg-white/10 text-white"
                                        : "text-gray-300 hover:bg-white/5 hover:text-white"
                                )}
                            >
                                <item.icon className={cn(
                                    "w-5 h-5 mr-3 transition-colors",
                                    isActive(item.path)
                                        ? "text-blue-400"
                                        : "group-hover:text-blue-400"
                                )} />
                                <span className="font-medium">{item.label}</span>
                            </button>
                        ))}
                    </nav>

                    {/* Footer */}
                    <div className="p-4 border-t border-white/10 space-y-3">
                        {/* Feedback Button */}
                        <button
                            onClick={onFeedbackClick}
                            className="flex items-center w-full px-4 py-2 bg-gray-800 hover:bg-gray-700 rounded-lg text-gray-300 hover:text-white transition-colors text-sm"
                            aria-label="Send Feedback"
                        >
                            <MessageSquare size={16} className="mr-2" />
                            Send Feedback
                        </button>

                        {/* Update Available Button */}
                        {hasUpdate && (
                            <button
                                onClick={onUpdateClick}
                                className="flex items-center w-full px-4 py-2 bg-purple-600/30 border border-purple-500 rounded-lg text-purple-200 hover:bg-purple-600/50 transition-colors text-sm"
                            >
                                <Download size={16} className="mr-2" />
                                Update Available
                            </button>
                        )}

                        {/* Status and Version */}
                        <div className="flex items-center justify-between px-4 py-2 text-sm text-gray-400">
                            <div className="flex items-center">
                                <div className="w-2 h-2 mr-2 bg-green-500 rounded-full animate-pulse" />
                                Online
                            </div>
                            {currentVersion && (
                                <span 
                                    className="font-mono text-xs cursor-pointer hover:text-gray-300"
                                    onClick={onUpdateClick}
                                    title="Click to check for updates"
                                >
                                    v{currentVersion}
                                </span>
                            )}
                        </div>
                    </div>
                </div>
            </div>

            {/* Overlay for mobile */}
            {isOpen && (
                <div
                    className="fixed inset-0 z-30 bg-black/50 backdrop-blur-sm md:hidden"
                    onClick={() => setIsOpen(false)}
                />
            )}
        </>
    );
}
