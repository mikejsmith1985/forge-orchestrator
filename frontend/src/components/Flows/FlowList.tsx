import React, { useState } from 'react';
import { Plus, Play, Edit, Trash2 } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

// Mock data for flows
interface Flow {
    id: string;
    name: string;
    description: string;
    updatedAt: string;
    status: 'active' | 'inactive';
}

const MOCK_FLOWS: Flow[] = [
    {
        id: '1',
        name: 'Customer Onboarding',
        description: 'Orchestrate welcome email and database setup',
        updatedAt: '2023-10-27T10:00:00Z',
        status: 'active',
    },
    {
        id: '2',
        name: 'Data Processing',
        description: 'Process incoming CSV files and update metrics',
        updatedAt: '2023-10-26T14:30:00Z',
        status: 'inactive',
    },
];

const FlowList: React.FC = () => {
    const navigate = useNavigate();
    const [flows, setFlows] = useState<Flow[]>(MOCK_FLOWS);

    const handleDelete = (id: string) => {
        if (confirm('Are you sure you want to delete this flow?')) {
            setFlows(flows.filter((f) => f.id !== id));
        }
    };

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
                                    className="p-2 hover:bg-red-500/10 rounded-lg text-slate-400 hover:text-red-400 transition-colors"
                                    title="Delete"
                                >
                                    <Trash2 size={18} />
                                </button>
                            </div>
                        </div>

                        <h3 className="text-lg font-semibold text-white mb-2">{flow.name}</h3>
                        <p className="text-slate-400 text-sm mb-4 line-clamp-2">
                            {flow.description}
                        </p>

                        <div className="flex items-center justify-between text-xs text-slate-500 border-t border-slate-700 pt-4">
                            <span>Updated {new Date(flow.updatedAt).toLocaleDateString()}</span>
                            <span
                                className={`px-2 py-1 rounded-full ${flow.status === 'active'
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
        </div>
    );
};

export default FlowList;
