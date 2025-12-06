import React, { useState } from 'react';
import { Key, Terminal } from 'lucide-react';
import { KeyManagement } from './KeyManagement';
import { TerminalSettings } from './TerminalSettings';

/**
 * Settings Component - Main settings page with tabs
 */
export const Settings: React.FC = () => {
    const [activeTab, setActiveTab] = useState<'terminal' | 'api-keys'>('terminal');

    return (
        <div className="h-full flex flex-col bg-slate-900">
            {/* Tabs */}
            <div className="border-b border-slate-700 bg-slate-800">
                <div className="flex gap-2 px-6">
                    <button
                        onClick={() => setActiveTab('terminal')}
                        className={`flex items-center gap-2 px-4 py-3 border-b-2 transition-colors ${
                            activeTab === 'terminal'
                                ? 'border-blue-500 text-blue-400'
                                : 'border-transparent text-slate-400 hover:text-slate-300'
                        }`}
                    >
                        <Terminal className="w-4 h-4" />
                        <span className="font-medium">Terminal</span>
                    </button>
                    <button
                        onClick={() => setActiveTab('api-keys')}
                        className={`flex items-center gap-2 px-4 py-3 border-b-2 transition-colors ${
                            activeTab === 'api-keys'
                                ? 'border-blue-500 text-blue-400'
                                : 'border-transparent text-slate-400 hover:text-slate-300'
                        }`}
                    >
                        <Key className="w-4 h-4" />
                        <span className="font-medium">API Keys</span>
                    </button>
                </div>
            </div>

            {/* Tab Content */}
            <div className="flex-1 overflow-auto">
                {activeTab === 'terminal' && <TerminalSettings />}
                {activeTab === 'api-keys' && <KeyManagement />}
            </div>
        </div>
    );
};
