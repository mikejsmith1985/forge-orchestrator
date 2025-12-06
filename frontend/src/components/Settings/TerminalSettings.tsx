import React, { useState, useEffect } from 'react';
import { Terminal, Settings as SettingsIcon, AlertCircle, CheckCircle } from 'lucide-react';

interface ShellConfig {
    type: 'bash' | 'cmd' | 'powershell' | 'wsl';
    wsl_distro?: string;
    wsl_user?: string;
    root_dir?: string;
}

interface Config {
    shell: ShellConfig;
}

/**
 * TerminalSettings Component
 * Allows users to configure their preferred terminal shell
 */
export const TerminalSettings: React.FC = () => {
    const [config, setConfig] = useState<Config | null>(null);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [message, setMessage] = useState<{ type: 'success' | 'error', text: string } | null>(null);
    const [isWindows, setIsWindows] = useState(false);

    useEffect(() => {
        loadConfig();
        detectPlatform();
    }, []);

    const detectPlatform = async () => {
        try {
            await fetch('/api/config');
            // Detect platform from user agent
            const platform = navigator.platform.toLowerCase();
            setIsWindows(platform.includes('win'));
        } catch (err) {
            console.error('Failed to detect platform:', err);
        }
    };

    const loadConfig = async () => {
        try {
            const res = await fetch('/api/config');
            if (!res.ok) throw new Error('Failed to load config');
            const data = await res.json();
            setConfig(data);
        } catch (err) {
            console.error('Failed to load config:', err);
            setMessage({ type: 'error', text: 'Failed to load configuration' });
        } finally {
            setLoading(false);
        }
    };

    const saveConfig = async () => {
        if (!config) return;

        setSaving(true);
        setMessage(null);

        try {
            const res = await fetch('/api/config', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(config),
            });

            if (!res.ok) throw new Error('Failed to save config');

            setMessage({ type: 'success', text: 'Configuration saved! Reload the terminal to apply changes.' });
        } catch (err) {
            console.error('Failed to save config:', err);
            setMessage({ type: 'error', text: 'Failed to save configuration' });
        } finally {
            setSaving(false);
        }
    };

    const updateShellType = (type: ShellConfig['type']) => {
        if (!config) return;
        setConfig({
            ...config,
            shell: { ...config.shell, type },
        });
    };

    const updateWSLDistro = (distro: string) => {
        if (!config) return;
        setConfig({
            ...config,
            shell: { ...config.shell, wsl_distro: distro },
        });
    };

    const updateRootDir = (dir: string) => {
        if (!config) return;
        setConfig({
            ...config,
            shell: { ...config.shell, root_dir: dir },
        });
    };

    if (loading) {
        return (
            <div className="flex items-center justify-center h-full">
                <div className="text-slate-400">Loading configuration...</div>
            </div>
        );
    }

    if (!config) {
        return (
            <div className="flex items-center justify-center h-full">
                <div className="text-red-400">Failed to load configuration</div>
            </div>
        );
    }

    return (
        <div className="p-6 max-w-4xl mx-auto">
            <div className="flex items-center gap-3 mb-6">
                <Terminal className="w-8 h-8 text-blue-500" />
                <h1 className="text-3xl font-bold text-slate-100">Terminal Settings</h1>
            </div>

            {/* Shell Selection */}
            <div className="bg-slate-800 rounded-lg p-6 mb-6">
                <h2 className="text-xl font-semibold text-slate-100 mb-4 flex items-center gap-2">
                    <SettingsIcon className="w-5 h-5" />
                    Shell Type
                </h2>
                <p className="text-slate-400 mb-4">
                    Select your preferred terminal shell. The terminal will restart with the new shell on next connection.
                </p>

                <div className="space-y-3">
                    {/* Bash Option (Unix/Linux) */}
                    {!isWindows && (
                        <label className={`flex items-center p-4 rounded-lg border-2 cursor-pointer transition-all ${
                            config.shell.type === 'bash'
                                ? 'border-blue-500 bg-blue-500/10'
                                : 'border-slate-700 bg-slate-900 hover:border-slate-600'
                        }`}>
                            <input
                                type="radio"
                                name="shell"
                                value="bash"
                                checked={config.shell.type === 'bash'}
                                onChange={() => updateShellType('bash')}
                                className="w-4 h-4 text-blue-500"
                            />
                            <div className="ml-3">
                                <div className="text-slate-100 font-medium">Bash</div>
                                <div className="text-sm text-slate-400">Standard Unix/Linux shell</div>
                            </div>
                        </label>
                    )}

                    {/* CMD Option (Windows) */}
                    {isWindows && (
                        <label className={`flex items-center p-4 rounded-lg border-2 cursor-pointer transition-all ${
                            config.shell.type === 'cmd'
                                ? 'border-blue-500 bg-blue-500/10'
                                : 'border-slate-700 bg-slate-900 hover:border-slate-600'
                        }`}>
                            <input
                                type="radio"
                                name="shell"
                                value="cmd"
                                checked={config.shell.type === 'cmd'}
                                onChange={() => updateShellType('cmd')}
                                className="w-4 h-4 text-blue-500"
                            />
                            <div className="ml-3">
                                <div className="text-slate-100 font-medium">Command Prompt (CMD)</div>
                                <div className="text-sm text-slate-400">Windows command line</div>
                            </div>
                        </label>
                    )}

                    {/* PowerShell Option (Windows) */}
                    {isWindows && (
                        <label className={`flex items-center p-4 rounded-lg border-2 cursor-pointer transition-all ${
                            config.shell.type === 'powershell'
                                ? 'border-blue-500 bg-blue-500/10'
                                : 'border-slate-700 bg-slate-900 hover:border-slate-600'
                        }`}>
                            <input
                                type="radio"
                                name="shell"
                                value="powershell"
                                checked={config.shell.type === 'powershell'}
                                onChange={() => updateShellType('powershell')}
                                className="w-4 h-4 text-blue-500"
                            />
                            <div className="ml-3">
                                <div className="text-slate-100 font-medium">PowerShell</div>
                                <div className="text-sm text-slate-400">Modern Windows shell with scripting capabilities</div>
                            </div>
                        </label>
                    )}

                    {/* WSL Option (Windows) */}
                    {isWindows && (
                        <label className={`flex items-center p-4 rounded-lg border-2 cursor-pointer transition-all ${
                            config.shell.type === 'wsl'
                                ? 'border-blue-500 bg-blue-500/10'
                                : 'border-slate-700 bg-slate-900 hover:border-slate-600'
                        }`}>
                            <input
                                type="radio"
                                name="shell"
                                value="wsl"
                                checked={config.shell.type === 'wsl'}
                                onChange={() => updateShellType('wsl')}
                                className="w-4 h-4 text-blue-500"
                            />
                            <div className="ml-3 flex-1">
                                <div className="text-slate-100 font-medium">WSL (Windows Subsystem for Linux)</div>
                                <div className="text-sm text-slate-400 mb-2">Linux environment on Windows</div>
                                
                                {config.shell.type === 'wsl' && (
                                    <div className="mt-3 space-y-3">
                                        <div>
                                            <label className="block text-sm text-slate-300 mb-1">
                                                WSL Distribution (optional)
                                            </label>
                                            <input
                                                type="text"
                                                placeholder="e.g., Ubuntu-24.04"
                                                value={config.shell.wsl_distro || ''}
                                                onChange={(e) => updateWSLDistro(e.target.value)}
                                                className="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded text-slate-100 placeholder-slate-500 focus:outline-none focus:border-blue-500"
                                            />
                                            <p className="text-xs text-slate-500 mt-1">
                                                Leave empty to use default. Run 'wsl --list' in CMD to see installed distributions.
                                            </p>
                                        </div>
                                        <div>
                                            <label className="block text-sm text-slate-300 mb-1">
                                                Starting Directory (optional)
                                            </label>
                                            <input
                                                type="text"
                                                placeholder="e.g., C:\Users\mike\projects\forge-orchestrator"
                                                value={config.shell.root_dir || ''}
                                                onChange={(e) => updateRootDir(e.target.value)}
                                                className="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded text-slate-100 placeholder-slate-500 focus:outline-none focus:border-blue-500"
                                            />
                                            <p className="text-xs text-slate-500 mt-1">
                                                Leave empty to use current working directory. Use Windows path format - it will be converted automatically.
                                            </p>
                                        </div>
                                    </div>
                                )}
                            </div>
                        </label>
                    )}
                </div>
            </div>

            {/* Help Section */}
            <div className="bg-slate-800/50 rounded-lg p-6 mb-6 border border-slate-700">
                <h3 className="text-lg font-semibold text-slate-100 mb-3 flex items-center gap-2">
                    <AlertCircle className="w-5 h-5 text-blue-500" />
                    Troubleshooting
                </h3>
                <ul className="space-y-2 text-slate-400 text-sm">
                    <li>• If terminal fails to connect, try changing the shell type</li>
                    {isWindows && (
                        <>
                            <li>• For WSL: Ensure WSL is installed and configured (<code className="text-blue-400">wsl --install</code>)</li>
                            <li>• For PowerShell: Make sure PowerShell is in your system PATH</li>
                            <li>• WSL Starting Directory: Use Windows path format (e.g., C:\Users\mike\projects) - it will be auto-converted to WSL format</li>
                        </>
                    )}
                    <li>• Check browser console (F12) for detailed error messages</li>
                    <li>• Terminal will automatically reconnect after changing settings</li>
                </ul>
            </div>

            {/* Save Button */}
            <div className="flex flex-col gap-4">
                {/* Status Message - Now at bottom near save button */}
                {message && (
                    <div className={`flex items-center gap-2 p-4 rounded-lg ${
                        message.type === 'success' 
                            ? 'bg-green-500/20 text-green-400 border border-green-500/30'
                            : 'bg-red-500/20 text-red-400 border border-red-500/30'
                    }`}>
                        {message.type === 'success' ? (
                            <CheckCircle className="w-5 h-5" />
                        ) : (
                            <AlertCircle className="w-5 h-5" />
                        )}
                        <span>{message.text}</span>
                    </div>
                )}
                
                <div className="flex justify-end">
                    <button
                        onClick={saveConfig}
                        disabled={saving}
                        className="px-6 py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-slate-700 disabled:cursor-not-allowed text-white rounded-lg font-medium transition-colors"
                    >
                        {saving ? 'Saving...' : 'Save Configuration'}
                    </button>
                </div>
            </div>
        </div>
    );
};
