import { useState, useEffect } from 'react';
import { Key, CheckCircle, XCircle, Loader2, Save } from 'lucide-react';
import { cn } from '../../lib/utils';

// Educational Comment: Defining the shape of our key status data
interface KeyStatus {
    provider: string;
    isSet: boolean;
}

export function KeyManagement() {
    const [statuses, setStatuses] = useState<KeyStatus[]>([]);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState<string | null>(null);
    const [inputs, setInputs] = useState<Record<string, string>>({});

    // Educational Comment: Fetch initial status on mount
    useEffect(() => {
        fetchStatus();
    }, []);

    const fetchStatus = async () => {
        try {
            const response = await fetch('/api/keys/status');
            if (response.ok) {
                const data = await response.json();
                setStatuses(data.keys || []);
            }
        } catch (error) {
            console.error('Failed to fetch key status:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleSave = async (provider: string) => {
        const key = inputs[provider];
        if (!key) return;

        setSaving(provider);
        try {
            // Educational Comment: We send the key to the backend via POST.
            // The backend is responsible for securely storing it (e.g., in the system keyring).
            // We never read the key back to the frontend for security reasons.
            const response = await fetch('/api/keys', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ provider, key }),
            });

            if (response.ok) {
                await fetchStatus();
                setInputs(prev => ({ ...prev, [provider]: '' })); // Clear input on success
            }
        } catch (error) {
            console.error('Failed to save key:', error);
        } finally {
            setSaving(null);
        }
    };

    // Educational Comment: Helper to get the status icon based on isSet state
    const getStatusIcon = (isSet: boolean) => {
        if (isSet) {
            return <CheckCircle className="w-5 h-5 text-green-500" />;
        }
        return <XCircle className="w-5 h-5 text-red-500" />;
    };

    // Default providers if API fails or returns empty (for initial UI state)
    const displayProviders = statuses.length > 0 ? statuses : [
        { provider: 'anthropic', isSet: false },
        { provider: 'openai', isSet: false }
    ];

    return (
        <div className="p-6 max-w-4xl mx-auto text-white">
            <div className="mb-8">
                <h2 className="text-3xl font-bold flex items-center gap-3 mb-2">
                    <Key className="w-8 h-8 text-blue-400" />
                    Key Management
                </h2>
                <p className="text-gray-400">
                    Securely manage your API keys. Keys are stored in your system's secure keyring and are never displayed.
                </p>
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
