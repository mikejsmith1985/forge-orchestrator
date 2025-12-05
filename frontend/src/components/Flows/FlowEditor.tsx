import React, { useCallback, useRef, useState, useEffect } from 'react';
import ReactFlow, {
    ReactFlowProvider,
    addEdge,
    useNodesState,
    useEdgesState,
    Controls,
    Background,
    type Connection,
    type Edge,
    type Node,
    type ReactFlowInstance,
} from 'reactflow';
import 'reactflow/dist/style.css';
import { Save, Play, ArrowLeft, Loader2, Check, AlertCircle } from 'lucide-react';
import { useNavigate, useParams } from 'react-router-dom';

// Agent role configuration with user-friendly labels
const agentRoles = [
    { value: 'Architect', label: 'Planner / Architect' },
    { value: 'Implementation', label: 'Developer / Coder' },
    { value: 'Test', label: 'QA / Tester' },
    { value: 'Optimizer', label: 'Auditor / Optimizer' },
];

// Initial nodes for a new flow
const initialNodes: Node[] = [
    {
        id: '1',
        type: 'input',
        data: { label: 'Start Node' },
        position: { x: 250, y: 5 },
    },
];

interface Toast {
    type: 'success' | 'error';
    message: string;
}

const FlowEditorContent: React.FC = () => {
    const reactFlowWrapper = useRef<HTMLDivElement>(null);
    const navigate = useNavigate();
    const { id } = useParams();

    const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
    const [edges, setEdges, onEdgesChange] = useEdgesState([]);
    const [reactFlowInstance, setReactFlowInstance] = useState<ReactFlowInstance | null>(null);
    
    // Flow metadata
    const [flowName, setFlowName] = useState('New Flow');
    
    // Selected node for configuration
    const [selectedNode, setSelectedNode] = useState<Node | null>(null);
    
    // Loading states
    const [loadingFlow, setLoadingFlow] = useState(false);
    const [saving, setSaving] = useState(false);
    const [executing, setExecuting] = useState(false);
    
    // Toast notification
    const [toast, setToast] = useState<Toast | null>(null);

    // Load existing flow data when editing
    useEffect(() => {
        if (id) {
            loadFlow(id);
        }
    }, [id]);

    // Auto-hide toast after 3 seconds
    useEffect(() => {
        if (toast) {
            const timer = setTimeout(() => setToast(null), 3000);
            return () => clearTimeout(timer);
        }
    }, [toast]);

    const loadFlow = async (flowId: string) => {
        try {
            setLoadingFlow(true);
            const response = await fetch(`/api/flows/${flowId}`);
            
            if (!response.ok) {
                throw new Error('Failed to load flow');
            }
            
            const flow = await response.json();
            setFlowName(flow.name);
            
            // Parse the graph data
            if (flow.data) {
                const graphData = JSON.parse(flow.data);
                if (graphData.nodes) setNodes(graphData.nodes);
                if (graphData.edges) setEdges(graphData.edges);
            }
        } catch (err) {
            console.error('Error loading flow:', err);
            setToast({ type: 'error', message: 'Failed to load flow' });
        } finally {
            setLoadingFlow(false);
        }
    };

    const onConnect = useCallback(
        (params: Connection | Edge) => setEdges((eds) => addEdge(params, eds)),
        [setEdges]
    );

    const onDragOver = useCallback((event: React.DragEvent) => {
        event.preventDefault();
        event.dataTransfer.dropEffect = 'move';
    }, []);

    const onDrop = useCallback(
        (event: React.DragEvent) => {
            event.preventDefault();

            const type = event.dataTransfer.getData('application/reactflow');

            if (typeof type === 'undefined' || !type || !reactFlowInstance) {
                return;
            }

            const position = reactFlowInstance.screenToFlowPosition({
                x: event.clientX,
                y: event.clientY,
            });

            const newNode: Node = {
                id: `${type}-${nodes.length + 1}`,
                type,
                position,
                data: { label: `${type} node`, role: 'Implementation' },
            };

            setNodes((nds) => nds.concat(newNode));
        },
        [reactFlowInstance, nodes, setNodes]
    );

    // Handle node selection
    const onNodeClick = useCallback((_: React.MouseEvent, node: Node) => {
        setSelectedNode(node);
    }, []);

    // Handle pane click to deselect
    const onPaneClick = useCallback(() => {
        setSelectedNode(null);
    }, []);

    // Update node role
    const handleRoleChange = useCallback((nodeId: string, role: string) => {
        setNodes((nds) =>
            nds.map((node) =>
                node.id === nodeId
                    ? { ...node, data: { ...node.data, role } }
                    : node
            )
        );
        // Update selected node state
        setSelectedNode((prev) => prev && prev.id === nodeId 
            ? { ...prev, data: { ...prev.data, role } }
            : prev
        );
    }, [setNodes]);

    // Update node label
    const handleLabelChange = useCallback((nodeId: string, label: string) => {
        setNodes((nds) =>
            nds.map((node) =>
                node.id === nodeId
                    ? { ...node, data: { ...node.data, label } }
                    : node
            )
        );
        setSelectedNode((prev) => prev && prev.id === nodeId 
            ? { ...prev, data: { ...prev.data, label } }
            : prev
        );
    }, [setNodes]);

    const handleSave = async () => {
        if (!flowName.trim()) {
            setToast({ type: 'error', message: 'Please enter a flow name' });
            return;
        }

        try {
            setSaving(true);
            
            const flowData = {
                name: flowName,
                data: JSON.stringify({
                    nodes: nodes,
                    edges: edges
                }),
                status: 'active'
            };

            const url = id ? `/api/flows/${id}` : '/api/flows';
            const method = id ? 'PUT' : 'POST';

            const response = await fetch(url, {
                method,
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(flowData),
            });

            if (!response.ok) {
                throw new Error('Failed to save flow');
            }

            setToast({ type: 'success', message: 'Flow saved successfully!' });
            
            // Navigate back to list after a brief delay to show success message
            setTimeout(() => navigate('/flows'), 1000);
        } catch (err) {
            console.error('Error saving flow:', err);
            setToast({ type: 'error', message: 'Failed to save flow' });
        } finally {
            setSaving(false);
        }
    };

    const handleExecute = async () => {
        if (!id) {
            setToast({ type: 'error', message: 'Save the flow first before executing' });
            return;
        }

        try {
            setExecuting(true);
            
            const response = await fetch(`/api/flows/${id}/execute`, {
                method: 'POST',
            });

            if (!response.ok) {
                throw new Error('Failed to execute flow');
            }

            setToast({ type: 'success', message: 'Flow execution started!' });
        } catch (err) {
            console.error('Error executing flow:', err);
            setToast({ type: 'error', message: 'Failed to execute flow' });
        } finally {
            setExecuting(false);
        }
    };

    if (loadingFlow) {
        return (
            <div className="flex items-center justify-center h-full bg-slate-900">
                <div className="flex items-center gap-3 text-slate-400">
                    <Loader2 className="animate-spin" size={24} />
                    <span>Loading flow...</span>
                </div>
            </div>
        );
    }

    return (
        <div className="flex flex-col h-full bg-slate-900">
            {/* Toast Notification */}
            {toast && (
                <div
                    className={`fixed top-4 right-4 z-50 flex items-center gap-2 px-4 py-3 rounded-lg shadow-lg ${
                        toast.type === 'success'
                            ? 'bg-green-500/90 text-white'
                            : 'bg-red-500/90 text-white'
                    }`}
                >
                    {toast.type === 'success' ? <Check size={18} /> : <AlertCircle size={18} />}
                    {toast.message}
                </div>
            )}

            {/* Toolbar */}
            <div className="flex items-center justify-between p-4 border-b border-slate-700 bg-slate-800">
                <div className="flex items-center gap-4">
                    <button
                        onClick={() => navigate('/flows')}
                        className="p-2 hover:bg-slate-700 rounded-lg text-slate-400 hover:text-white transition-colors"
                    >
                        <ArrowLeft size={20} />
                    </button>
                    <div>
                        <input
                            type="text"
                            value={flowName}
                            onChange={(e) => setFlowName(e.target.value)}
                            className="text-xl font-bold text-white bg-transparent border-none outline-none focus:ring-2 focus:ring-blue-500 rounded px-2 py-1"
                            placeholder="Flow Name"
                        />
                        <p className="text-xs text-slate-400 px-2">
                            {nodes.length} nodes • {edges.length} connections
                        </p>
                    </div>
                </div>
                <div className="flex gap-2">
                    <button
                        onClick={handleSave}
                        disabled={saving}
                        className="flex items-center gap-2 px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors disabled:opacity-50"
                    >
                        {saving ? (
                            <Loader2 size={18} className="animate-spin" />
                        ) : (
                            <Save size={18} />
                        )}
                        {saving ? 'Saving...' : 'Save'}
                    </button>
                    <button
                        onClick={handleExecute}
                        disabled={executing || !id}
                        className="flex items-center gap-2 px-4 py-2 bg-green-600 hover:bg-green-500 text-white rounded-lg transition-colors disabled:opacity-50"
                    >
                        {executing ? (
                            <Loader2 size={18} className="animate-spin" />
                        ) : (
                            <Play size={18} />
                        )}
                        {executing ? 'Executing...' : 'Execute'}
                    </button>
                </div>
            </div>

            <div className="flex flex-1 overflow-hidden">
                {/* Sidebar for Drag & Drop */}
                <div className="w-64 bg-slate-800 border-r border-slate-700 p-4 flex flex-col gap-4">
                    <h3 className="text-sm font-semibold text-slate-400 uppercase tracking-wider">
                        Components
                    </h3>
                    <div className="space-y-2">
                        <div
                            className="bg-slate-700 p-3 rounded-lg cursor-grab hover:bg-slate-600 transition-colors border border-slate-600 flex items-center gap-2"
                            onDragStart={(event) =>
                                event.dataTransfer.setData('application/reactflow', 'input')
                            }
                            draggable
                        >
                            <div className="w-3 h-3 rounded-full bg-blue-500" />
                            <span className="text-sm text-white">Input Node</span>
                        </div>
                        <div
                            className="bg-slate-700 p-3 rounded-lg cursor-grab hover:bg-slate-600 transition-colors border border-slate-600 flex items-center gap-2"
                            onDragStart={(event) =>
                                event.dataTransfer.setData('application/reactflow', 'default')
                            }
                            draggable
                        >
                            <div className="w-3 h-3 rounded-full bg-slate-400" />
                            <span className="text-sm text-white">Agent Node</span>
                        </div>
                        <div
                            className="bg-slate-700 p-3 rounded-lg cursor-grab hover:bg-slate-600 transition-colors border border-slate-600 flex items-center gap-2"
                            onDragStart={(event) =>
                                event.dataTransfer.setData('application/reactflow', 'output')
                            }
                            draggable
                        >
                            <div className="w-3 h-3 rounded-full bg-green-500" />
                            <span className="text-sm text-white">Output Node</span>
                        </div>
                    </div>

                    {/* Node Configuration Panel */}
                    {selectedNode && (
                        <div className="mt-4 p-4 bg-slate-700/50 rounded-lg border border-slate-600">
                            <h4 className="text-sm font-semibold text-slate-300 mb-3">
                                Node Configuration
                            </h4>
                            <div className="space-y-3">
                                <div>
                                    <label className="block text-xs text-slate-400 mb-1">
                                        Label
                                    </label>
                                    <input
                                        type="text"
                                        value={selectedNode.data.label || ''}
                                        onChange={(e) => handleLabelChange(selectedNode.id, e.target.value)}
                                        className="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded text-sm text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                                    />
                                </div>
                                <div>
                                    <label className="block text-xs text-slate-400 mb-1">
                                        Agent Role
                                    </label>
                                    <select
                                        value={selectedNode.data.role || 'Implementation'}
                                        onChange={(e) => handleRoleChange(selectedNode.id, e.target.value)}
                                        className="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded text-sm text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                                    >
                                        {agentRoles.map((role) => (
                                            <option key={role.value} value={role.value}>
                                                {role.label}
                                            </option>
                                        ))}
                                    </select>
                                </div>
                            </div>
                        </div>
                    )}

                    <div className="mt-auto p-4 bg-slate-700/50 rounded-lg border border-slate-700">
                        <p className="text-xs text-slate-400">
                            {selectedNode 
                                ? 'Configure the selected node above.'
                                : 'Drag components to the canvas to build your flow.'}
                        </p>
                    </div>
                </div>

                {/* React Flow Canvas */}
                <div className="flex-1 h-full" ref={reactFlowWrapper}>
                    <ReactFlow
                        nodes={nodes}
                        edges={edges}
                        onNodesChange={onNodesChange}
                        onEdgesChange={onEdgesChange}
                        onConnect={onConnect}
                        onInit={setReactFlowInstance}
                        onDrop={onDrop}
                        onDragOver={onDragOver}
                        onNodeClick={onNodeClick}
                        onPaneClick={onPaneClick}
                        fitView
                        className="bg-slate-900"
                    >
                        <Controls className="bg-white text-black" />
                        <Background color="#334155" gap={16} />
                    </ReactFlow>
                </div>
            </div>
        </div>
    );
};

const FlowEditor: React.FC = () => {
    return (
        <ReactFlowProvider>
            <FlowEditorContent />
        </ReactFlowProvider>
    );
};

export default FlowEditor;
