import { useState, useEffect } from 'react';
import { Key, CheckCircle, XCircle, Loader2, Save } from 'lucide-react';
import { cn } from '../../lib/utils';
import { useToastContext } from '../../lib/ToastContext';

interface KeyStatus {
    provider: string;
    isSet: boolean;
}

interface KeyStatusResponse {
    keys: KeyStatus[];
}

interface SaveKeyResponse {
    status: string;
    message: string;
}

export function KeyManagement() {
    const [statuses, setStatuses] = useState<KeyStatus[]>([]);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState<string | null>(null);
    const [inputs, setInputs] = useState<Record<string, string>>({});
    const toast = useToastContext();

    useEffect(() => {
        fetchStatus();
    }, []);

    const fetchStatus = async () => {
        try {
            const response = await fetch('/api/keys/status');
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}`);
            }
            const data: KeyStatusResponse = await response.json();
            setStatuses(data.keys || []);
        } catch (error) {
            console.error('Failed to fetch key status:', error);
            toast.error('Failed to load API key status');
        } finally {
            setLoading(false);
        }
    };

    const handleSave = async (provider: string) => {
        const key = inputs[provider];
        if (!key) return;

        setSaving(provider);
        try {
            const response = await fetch('/api/keys', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ provider, key }),
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || `HTTP ${response.status}`);
            }

            const data: SaveKeyResponse = await response.json();
            toast.success(data.message || 'API key saved successfully');
            
            await fetchStatus();
            setInputs(prev => ({ ...prev, [provider]: '' }));
        } catch (error) {
            console.error('Failed to save key:', error);
            toast.error(`Failed to save API key: ${error instanceof Error ? error.message : 'Unknown error'}`);
        } finally {
            setSaving(null);
        }
    };

    const getStatusIcon = (isSet: boolean) => {
        if (isSet) {
            return <CheckCircle className="w-5 h-5 text-green-500" />;
        }
        return <XCircle className="w-5 h-5 text-red-500" />;
    };

    // Default providers if API fails or returns empty
    const displayProviders = statuses.length > 0 ? statuses : [
        { provider: 'anthropic', isSet: false },
        { provider: 'openai', isSet: false }
    ];

    return (
        <div className="p-6 max-w-4xl mx-auto text-white" data-testid="settings-view">
            <div className="mb-8">
                <h2 className="text-3xl font-bold flex items-center gap-3 mb-2">
                    <Key className="w-8 h-8 text-blue-400" />
                    Key Management
                </h2>
                <p className="text-gray-400">
                    Securely manage your API keys for LLM providers.
                </p>
                {/* Task 4.1: Security Assurance Text */}
                <div className="mt-4 p-4 bg-green-900/20 border border-green-500/30 rounded-lg" data-testid="security-assurance">
                    <div className="flex items-start gap-3">
                        <div className="p-1 bg-green-500/20 rounded">
                            <svg className="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                            </svg>
                        </div>
                        <div>
                            <h3 className="text-green-400 font-semibold mb-1">🔐 Secure Storage</h3>
                            <p className="text-sm text-green-300/80">
                                Your API keys are <strong>encrypted</strong> and stored securely in your operating system's native keyring 
                                (macOS Keychain, Windows Credential Manager, or Linux Secret Service). 
                                Keys are <strong>never</strong> exposed to the browser or stored in plain text.
                            </p>
                        </div>
                    </div>
                </div>
            </div>

            <div className="grid gap-6">
                {loading ? (
                    <div className="flex justify-center p-12">
                        <Loader2 className="w-8 h-8 animate-spin text-blue-400" />
                    </div>
                ) : (
                    displayProviders.map((status) => (
                        <div
                            key={status.provider}
                            data-testid={`provider-card-${status.provider}`}
                            className="bg-gray-900/50 border border-white/10 rounded-xl p-6 backdrop-blur-sm"
                        >
                            <div className="flex items-center justify-between mb-4">
                                <div className="flex items-center gap-3">
                                    <h3 className="text-xl font-semibold capitalize">
                                        {status.provider}
                                    </h3>
                                    <div className="flex items-center gap-2 text-sm bg-gray-800 px-3 py-1 rounded-full border border-white/5">
                                        {getStatusIcon(status.isSet)}
                                        <span className={status.isSet ? "text-green-400" : "text-red-400"}>
                                            {status.isSet ? "Configured" : "Not Configured"}
                                        </span>
                                    </div>
                                </div>
                            </div>

                            <div className="flex gap-4">
                                <div className="flex-1">
                                    <input
                                        type="password"
                                        placeholder={`Enter ${status.provider} API Key`}
                                        className="w-full bg-gray-950 border border-white/10 rounded-lg px-4 py-3 text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
                                        value={inputs[status.provider] || ''}
                                        onChange={(e) => setInputs(prev => ({
                                            ...prev,
                                            [status.provider]: e.target.value
                                        }))}
                                    />
                                </div>
                                <button
                                    onClick={() => handleSave(status.provider)}
                                    disabled={!inputs[status.provider] || saving === status.provider}
                                    className={cn(
                                        "flex items-center gap-2 px-6 py-3 rounded-lg font-medium transition-all",
                                        !inputs[status.provider] || saving === status.provider
                                            ? "bg-gray-800 text-gray-500 cursor-not-allowed"
                                            : "bg-blue-600 hover:bg-blue-500 text-white shadow-lg shadow-blue-500/20"
                                    )}
                                >
                                    {saving === status.provider ? (
                                        <Loader2 className="w-5 h-5 animate-spin" />
                                    ) : (
                                        <Save className="w-5 h-5" />
                                    )}
                                    Save Key
                                </button>
                            </div>
                            <p className="mt-3 text-sm text-gray-500">
                                {status.isSet
                                    ? "Key is currently set. Enter a new value to update it."
                                    : "Enter your API key to enable this provider."}
                            </p>
                        </div>
                    ))
                )}
            </div>
        </div>
    );
}
