import React, { useCallback, useRef, useState } from 'react';
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
import { Save, Play, ArrowLeft } from 'lucide-react';
import { useNavigate, useParams } from 'react-router-dom';

// Initial nodes for a new flow
const initialNodes: Node[] = [
    {
        id: '1',
        type: 'input',
        data: { label: 'Start Node' },
        position: { x: 250, y: 5 },
    },
];

const FlowEditorContent: React.FC = () => {
    const reactFlowWrapper = useRef<HTMLDivElement>(null);
    const navigate = useNavigate();
    const { id } = useParams();

    // EDUCATIONAL COMMENT: React Flow State Management
    // React Flow manages nodes and edges state internally but exposes hooks to control them.
    // useNodesState and useEdgesState are wrappers around useState that handle internal updates
    // (like dragging nodes) automatically while keeping our local state in sync.
    const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
    const [edges, setEdges, onEdgesChange] = useEdgesState([]);
    const [reactFlowInstance, setReactFlowInstance] = useState<ReactFlowInstance | null>(null);

    // EDUCATIONAL COMMENT: Handling Connections
    // The onConnect callback is triggered when a user connects two handles.
    // We use the addEdge utility to create a new edge object and update the edges state.
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

            // check if the dropped element is valid
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
                data: { label: `${type} node` },
            };

            setNodes((nds) => nds.concat(newNode));
        },
        [reactFlowInstance, nodes, setNodes]
    );

    const handleSave = () => {
        if (reactFlowInstance) {
            const flow = reactFlowInstance.toObject();
            console.log('Saving flow:', flow);
            // TODO: Call API to save flow
            alert('Flow saved! (Check console for object)');
        }
    };

    const handleExecute = () => {
        console.log('Executing flow:', id);
        // TODO: Call API to execute flow
        alert(`Executing flow ${id || 'new'}!`);
    };

    return (
        <div className="flex flex-col h-full bg-slate-900">
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
                        <h1 className="text-xl font-bold text-white">
                            {id ? `Edit Flow: ${id}` : 'New Flow'}
                        </h1>
                        <p className="text-xs text-slate-400">
                            {nodes.length} nodes • {edges.length} connections
                        </p>
                    </div>
                </div>
                <div className="flex gap-2">
                    <button
                        onClick={handleSave}
                        className="flex items-center gap-2 px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors"
                    >
                        <Save size={18} />
                        Save
                    </button>
                    <button
                        onClick={handleExecute}
                        className="flex items-center gap-2 px-4 py-2 bg-green-600 hover:bg-green-500 text-white rounded-lg transition-colors"
                    >
                        <Play size={18} />
                        Execute
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

                    <div className="mt-auto p-4 bg-slate-700/50 rounded-lg border border-slate-700">
                        <p className="text-xs text-slate-400">
                            Drag components to the canvas to build your flow.
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
