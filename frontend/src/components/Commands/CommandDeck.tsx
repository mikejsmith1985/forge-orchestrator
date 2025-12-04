import { useState, useEffect } from 'react';
import { Play, Trash2, Plus, X } from 'lucide-react';

// Educational Comment: This interface defines the shape of a command object
// as returned by the backend API. Using TypeScript interfaces ensures
// type safety and makes the code self-documenting.
interface Command {
    id: number;
    name: string;
    description: string;
    command: string;
    created_at?: string;
}

export function CommandDeck() {
    // Educational Comment: useState manages component state. When state changes,
    // React re-renders the component. We use separate state variables to
    // manage different aspects of the UI independently.
    const [commands, setCommands] = useState<Command[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [showForm, setShowForm] = useState(false);

    // Educational Comment: State for tracking command execution
    // runningCommandId tracks which specific command is currently executing to show loading state on the correct card
    const [runningCommandId, setRunningCommandId] = useState<number | null>(null);
    // executionResult stores the output from the backend to display in the modal
    const [executionResult, setExecutionResult] = useState<string | null>(null);
    const [showResultModal, setShowResultModal] = useState(false);

    // Educational Comment: Form state for creating new commands
    const [formData, setFormData] = useState({
        name: '',
        description: '',
        command: ''
    });

    // Educational Comment: useEffect runs side effects. The empty dependency array []
    // means this effect runs once when the component mounts, perfect for initial data fetching.
    useEffect(() => {
        fetchCommands();
    }, []);

    // Educational Comment: Async function to fetch commands from the backend API.
    // We use try/catch to handle potential network errors gracefully.
    const fetchCommands = async () => {
        try {
            setLoading(true);
            setError(null);
            const response = await fetch('/api/commands');

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            setCommands(data || []);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to fetch commands');
            console.error('Error fetching commands:', err);
        } finally {
            setLoading(false);
        }
    };

    // Educational Comment: Handler for creating a new command via POST request.
    // After successful creation, we refetch the command list to show the new command.
    const handleCreateCommand = async (e: React.FormEvent) => {
        e.preventDefault();

        try {
            const response = await fetch('/api/commands', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData),
            });

            if (!response.ok) {
                throw new Error('Failed to create command');
            }

            // Reset form and close modal
            setFormData({ name: '', description: '', command: '' });
            setShowForm(false);

            // Refetch to show new command
            await fetchCommands();
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to create command');
            console.error('Error creating command:', err);
        }
    };

    // Educational Comment: Handler for deleting a command via DELETE request.
    // We use optimistic UI update pattern here - we could immediately remove
    // from state, but refetching ensures consistency with backend.
    const handleDeleteCommand = async (id: number) => {
        if (!confirm('Are you sure you want to delete this command?')) {
            return;
        }

        try {
            const response = await fetch(`/api/commands/${id}`, {
                method: 'DELETE',
            });

            if (!response.ok) {
                throw new Error('Failed to delete command');
            }

            // Refetch to update the list
            await fetchCommands();
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to delete command');
            console.error('Error deleting command:', err);
        }
    };

    // Educational Comment: Handler for the "Run" button.
    // Executes the command via the backend API, which uses the LLM Gateway.
    // Educational Comment: Handler for the "Run" button.
    // Executes the command via the backend API, which uses the LLM Gateway.
    const handleRunCommand = async (command: Command) => {
        try {
            setRunningCommandId(command.id);
            const response = await fetch(`/api/commands/${command.id}/run`, {
                method: 'POST',
            });

            if (!response.ok) {
                throw new Error('Failed to run command');
            }

            const result = await response.json();
            console.log('Command executed:', result);

            // Educational Comment: Display the result to the user
            // We assume the backend returns a JSON object with an 'output' or similar field, 
            // or we just stringify the whole result if it's complex.
            // Adjust based on actual backend response structure.
            setExecutionResult(JSON.stringify(result, null, 2));
            setShowResultModal(true);
        } catch (err) {
            console.error('Error running command:', err);
            setError(err instanceof Error ? err.message : 'Failed to run command');
        } finally {
            setRunningCommandId(null);
        }
    };

    if (loading) {
        return (
            <div className="flex items-center justify-center h-full">
                <div className="text-gray-400 text-lg">Loading commands...</div>
            </div>
        );
    }

    return (
        <div className="p-8" data-testid="command-deck">
            {/* Header */}
            <div className="flex items-center justify-between mb-8">
                <div>
                    <h1 className="text-3xl font-bold text-white mb-2">Command Deck</h1>
                    <p className="text-gray-400">Manage and execute your command cards</p>
                </div>
                <button
                    onClick={() => setShowForm(true)}
                    className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                    data-testid="add-command-btn"
                >
                    <Plus size={20} />
                    Add Command
                </button>
            </div>

            {/* Error Display */}
            {error && (
                <div className="mb-6 p-4 bg-red-500/10 border border-red-500/50 rounded-lg text-red-400">
                    {error}
                </div>
            )}

            {/* Command Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {commands.length === 0 ? (
                    <div className="col-span-full text-center py-12 text-gray-500">
                        No commands yet. Create your first command to get started!
                    </div>
                ) : (
                    commands.map((cmd) => (
                        <div
                            key={cmd.id}
                            className="bg-gray-800/50 backdrop-blur-sm border border-white/10 rounded-lg p-6 hover:border-blue-500/50 transition-all group"
                            data-testid="command-card"
                        >
                            <h3 className="text-xl font-semibold text-white mb-2">{cmd.name}</h3>
                            <p className="text-gray-400 text-sm mb-4">{cmd.description}</p>
                            <div className="bg-gray-900/50 rounded p-3 mb-4">
                                <code className="text-green-400 text-xs font-mono break-all">
                                    {cmd.command}
                                </code>
                            </div>

                            <div className="flex gap-2">
                                <button
                                    onClick={() => handleRunCommand(cmd)}
                                    disabled={runningCommandId === cmd.id}
                                    className={`flex-1 flex items-center justify-center gap-2 px-4 py-2 rounded transition-colors ${runningCommandId === cmd.id
                                        ? 'bg-gray-600 cursor-not-allowed text-gray-300'
                                        : 'bg-green-600 hover:bg-green-700 text-white'
                                        }`}
                                    data-testid="run-command-btn"
                                >
                                    {runningCommandId === cmd.id ? (
                                        <div className="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent" />
                                    ) : (
                                        <Play size={16} />
                                    )}
                                    {runningCommandId === cmd.id ? 'Running...' : 'Run'}
                                </button>
                                <button
                                    onClick={() => handleDeleteCommand(cmd.id)}
                                    className="px-4 py-2 bg-red-600/20 hover:bg-red-600 text-red-400 hover:text-white rounded transition-colors"
                                    data-testid="delete-command-btn"
                                >
                                    <Trash2 size={16} />
                                </button>
                            </div>
                        </div>
                    ))
                )}
            </div>

            {/* Add Command Modal */}
            {showForm && (
                <>
                    {/* Educational Comment: This overlay provides a modal backdrop
                        that closes the form when clicked outside */}
                    <div
                        className="fixed inset-0 bg-black/50 backdrop-blur-sm z-40"
                        onClick={() => setShowForm(false)}
                    />

                    <div className="fixed inset-0 z-50 flex items-center justify-center p-4" data-testid="add-command-modal">
                        <div className="bg-gray-800 rounded-lg p-6 w-full max-w-md border border-white/10">
                            <div className="flex items-center justify-between mb-4">
                                <h2 className="text-2xl font-bold text-white">Add Command</h2>
                                <button
                                    onClick={() => setShowForm(false)}
                                    className="text-gray-400 hover:text-white transition-colors"
                                >
                                    <X size={24} />
                                </button>
                            </div>

                            <form onSubmit={handleCreateCommand} className="space-y-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-300 mb-2">
                                        Name
                                    </label>
                                    <input
                                        type="text"
                                        required
                                        value={formData.name}
                                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                        className="w-full px-4 py-2 bg-gray-900 border border-white/10 rounded-lg text-white focus:outline-none focus:border-blue-500"
                                        placeholder="e.g., Deploy to Production"
                                        data-testid="command-name-input"
                                    />
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-gray-300 mb-2">
                                        Description
                                    </label>
                                    <input
                                        type="text"
                                        required
                                        value={formData.description}
                                        onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                                        className="w-full px-4 py-2 bg-gray-900 border border-white/10 rounded-lg text-white focus:outline-none focus:border-blue-500"
                                        placeholder="Briefly describe what this command does"
                                        data-testid="command-description-input"
                                    />
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-gray-300 mb-2">
                                        Command
                                    </label>
                                    <textarea
                                        required
                                        value={formData.command}
                                        onChange={(e) => setFormData({ ...formData, command: e.target.value })}
                                        className="w-full px-4 py-2 bg-gray-900 border border-white/10 rounded-lg text-white focus:outline-none focus:border-blue-500 font-mono text-sm"
                                        placeholder="e.g., npm run build && npm run deploy"
                                        rows={4}
                                        data-testid="command-input"
                                    />
                                </div>

                                <div className="flex gap-3 pt-2">
                                    <button
                                        type="button"
                                        onClick={() => setShowForm(false)}
                                        className="flex-1 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
                                    >
                                        Cancel
                                    </button>
                                    <button
                                        type="submit"
                                        className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                                        data-testid="submit-command-btn"
                                    >
                                        Create
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                </>
            )}
            {/* Result Modal */}
            {showResultModal && (
                <>
                    <div
                        className="fixed inset-0 bg-black/50 backdrop-blur-sm z-40"
                        onClick={() => setShowResultModal(false)}
                    />
                    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
                        <div className="bg-gray-800 rounded-lg p-6 w-full max-w-2xl border border-white/10 max-h-[80vh] flex flex-col">
                            <div className="flex items-center justify-between mb-4">
                                <h2 className="text-2xl font-bold text-white">Execution Result</h2>
                                <button
                                    onClick={() => setShowResultModal(false)}
                                    className="text-gray-400 hover:text-white transition-colors"
                                >
                                    <X size={24} />
                                </button>
                            </div>
                            <div className="flex-1 overflow-auto bg-gray-900 rounded p-4">
                                <pre className="text-green-400 font-mono text-sm whitespace-pre-wrap">
                                    {executionResult}
                                </pre>
                            </div>
                            <div className="mt-4 flex justify-end">
                                <button
                                    onClick={() => setShowResultModal(false)}
                                    className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                                >
                                    Close
                                </button>
                            </div>
                        </div>
                    </div>
                </>
            )}
        </div>
    );
}
