import { useState } from 'react';
import { LayoutDashboard, GitGraph, Settings, Menu, X, BrainCircuit } from 'lucide-react';
import { cn } from '../../lib/utils';

// Educational Comment: Defining props interface for type safety and clarity on what data the component needs.
interface SidebarProps {
    currentView: 'architect' | 'ledger';
    onViewChange: (view: 'architect' | 'ledger') => void;
}

export function Sidebar({ currentView, onViewChange }: SidebarProps) {
    const [isOpen, setIsOpen] = useState(false);

    const navItems = [
        { icon: BrainCircuit, label: 'Architect', view: 'architect' as const },
        { icon: LayoutDashboard, label: 'Dashboard', view: 'ledger' as const },
        { icon: GitGraph, label: 'Flows', view: 'ledger' as const }, // Mapping Flows to ledger for now as per instructions
        { icon: Settings, label: 'Settings', view: 'architect' as const }, // Placeholder
    ];

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
                                    onViewChange(item.view);
                                    setIsOpen(false);
                                }}
                                className={cn(
                                    "flex items-center w-full px-4 py-3 rounded-lg transition-colors group text-left",
                                    currentView === item.view && (item.label === 'Architect' || item.label === 'Dashboard') // Simple active check
                                        ? "bg-white/10 text-white"
                                        : "text-gray-300 hover:bg-white/5 hover:text-white"
                                )}
                            >
                                <item.icon className={cn(
                                    "w-5 h-5 mr-3 transition-colors",
                                    currentView === item.view && (item.label === 'Architect' || item.label === 'Dashboard')
                                        ? "text-blue-400"
                                        : "group-hover:text-blue-400"
                                )} />
                                <span className="font-medium">{item.label}</span>
                            </button>
                        ))}
                    </nav>

                    {/* Footer */}
                    <div className="p-4 border-t border-white/10">
                        <div className="flex items-center px-4 py-3 text-sm text-gray-400">
                            <div className="w-2 h-2 mr-2 bg-green-500 rounded-full animate-pulse" />
                            System Online
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
