import React, { useState, useEffect } from 'react';
import { Plus, Play, Edit, Trash2, Loader2, RefreshCw } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

interface Flow {
    id: number;
    name: string;
    data: string;
    status: string;
    created_at: string;
}

const FlowList: React.FC = () => {
    const navigate = useNavigate();
    const [flows, setFlows] = useState<Flow[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [deletingId, setDeletingId] = useState<number | null>(null);

    useEffect(() => {
        fetchFlows();
    }, []);

    const fetchFlows = async () => {
        try {
            setLoading(true);
            setError(null);
            const response = await fetch('/api/flows');
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            const data = await response.json();
            setFlows(data || []);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to fetch flows');
            console.error('Error fetching flows:', err);
        } finally {
            setLoading(false);
        }
    };

    const handleDelete = async (id: number) => {
        if (!confirm('Are you sure you want to delete this flow?')) {
            return;
        }
        
        try {
            setDeletingId(id);
            const response = await fetch(`/api/flows/${id}`, {
                method: 'DELETE',
            });
            
            if (!response.ok) {
                throw new Error('Failed to delete flow');
            }
            
            // Refetch to update the list
            await fetchFlows();
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to delete flow');
            console.error('Error deleting flow:', err);
        } finally {
            setDeletingId(null);
        }
    };

    const getFlowDescription = (flow: Flow): string => {
        try {
            const data = JSON.parse(flow.data);
            const nodeCount = data.nodes?.length || 0;
            const edgeCount = data.edges?.length || 0;
            return `${nodeCount} nodes • ${edgeCount} connections`;
        } catch {
            return 'Flow configuration';
        }
    };

    if (loading) {
        return (
            <div className="p-6 flex items-center justify-center h-64">
                <div className="flex items-center gap-3 text-slate-400">
                    <Loader2 className="animate-spin" size={24} />
                    <span>Loading flows...</span>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="p-6">
                <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-4 flex items-center justify-between">
                    <div>
                        <p className="text-red-400 font-medium">Error loading flows</p>
                        <p className="text-red-400/70 text-sm">{error}</p>
                    </div>
                    <button
                        onClick={fetchFlows}
                        className="flex items-center gap-2 px-4 py-2 bg-red-500/20 hover:bg-red-500/30 text-red-400 rounded-lg transition-colors"
                    >
                        <RefreshCw size={18} />
                        Retry
                    </button>
                </div>
            </div>
        );
    }

    return (
        <div className="p-6 space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-2xl font-bold text-white">Flows</h1>
                    <p className="text-slate-400">Manage your agent orchestration pipelines</p>
                </div>
                <button
                    onClick={() => navigate('/flows/new')}
                    className="flex items-center gap-2 bg-blue-600 hover:bg-blue-500 text-white px-4 py-2 rounded-lg transition-colors"
                >
                    <Plus size={20} />
                    Create New Flow
                </button>
            </div>

            {flows.length === 0 ? (
                <div className="text-center py-12">
                    <p className="text-slate-400 mb-4">No flows yet. Create your first flow to get started.</p>
                    <button
                        onClick={() => navigate('/flows/new')}
                        className="inline-flex items-center gap-2 bg-blue-600 hover:bg-blue-500 text-white px-4 py-2 rounded-lg transition-colors"
                    >
                        <Plus size={20} />
                        Create New Flow
                    </button>
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {flows.map((flow) => (
                        <div
                            key={flow.id}
                            className="bg-slate-800 border border-slate-700 rounded-xl p-5 hover:border-slate-600 transition-all group"
                        >
                            <div className="flex justify-between items-start mb-4">
                                <div className="p-2 bg-blue-500/10 rounded-lg">
                                    <Play size={24} className="text-blue-400" />
                                </div>
                                <div className="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                                    <button
                                        onClick={() => navigate(`/flows/${flow.id}`)}
                                        className="p-2 hover:bg-slate-700 rounded-lg text-slate-400 hover:text-white transition-colors"
                                        title="Edit"
                                    >
                                        <Edit size={18} />
                                    </button>
                                    <button
                                        onClick={() => handleDelete(flow.id)}
                                        disabled={deletingId === flow.id}
                                        className="p-2 hover:bg-red-500/10 rounded-lg text-slate-400 hover:text-red-400 transition-colors disabled:opacity-50"
                                        title="Delete"
                                    >
                                        {deletingId === flow.id ? (
                                            <Loader2 size={18} className="animate-spin" />
                                        ) : (
                                            <Trash2 size={18} />
                                        )}
                                    </button>
                                </div>
                            </div>

                            <h3 className="text-lg font-semibold text-white mb-2">{flow.name}</h3>
                            <p className="text-slate-400 text-sm mb-4 line-clamp-2">
                                {getFlowDescription(flow)}
                            </p>

                            <div className="flex items-center justify-between text-xs text-slate-500 border-t border-slate-700 pt-4">
                                <span>Created {new Date(flow.created_at).toLocaleDateString()}</span>
                                <span
                                    className={`px-2 py-1 rounded-full ${
                                        flow.status === 'active'
                                            ? 'bg-green-500/10 text-green-400'
                                            : 'bg-slate-700 text-slate-400'
                                    }`}
                                >
                                    {flow.status}
                                </span>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default FlowList;
