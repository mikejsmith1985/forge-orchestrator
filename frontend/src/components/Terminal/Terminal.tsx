import React, { useEffect, useRef, useState, useCallback } from 'react';
import { Terminal as XTerm } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import { WebLinksAddon } from '@xterm/addon-web-links';
import '@xterm/xterm/css/xterm.css';
import { Eye, EyeOff } from 'lucide-react';

/**
 * Terminal Component - Task 1.2: React Terminal UI Implementation
 * 
 * This component integrates xterm.js with the Go backend's PTY WebSocket.
 * It provides a fully functional terminal that streams real shell output
 * from the local machine, confirming the application's identity as a
 * terminal-centric tool.
 */

interface TerminalProps {
    /** CSS class name for additional styling */
    className?: string;
    /** Callback when terminal connects */
    onConnect?: () => void;
    /** Callback when terminal disconnects */
    onDisconnect?: () => void;
}

export const Terminal: React.FC<TerminalProps> = ({
    className = '',
    onConnect,
    onDisconnect
}) => {
    const terminalRef = useRef<HTMLDivElement>(null);
    const xtermRef = useRef<XTerm | null>(null);
    const fitAddonRef = useRef<FitAddon | null>(null);
    const wsRef = useRef<WebSocket | null>(null);
    const [isConnected, setIsConnected] = useState(false);
    const [promptWatcherEnabled, setPromptWatcherEnabled] = useState(false);

    // Build WebSocket URL dynamically based on current location
    const getWebSocketUrl = useCallback(() => {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        // In development, backend is on port 8080
        const host = window.location.hostname;
        const port = window.location.port === '5173' ? '8080' : window.location.port;
        return `${protocol}//${host}:${port}/ws/pty`;
    }, []);

    // Toggle prompt watcher and notify backend
    const togglePromptWatcher = useCallback(() => {
        const newState = !promptWatcherEnabled;
        setPromptWatcherEnabled(newState);
        
        if (wsRef.current?.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify({
                type: 'prompt_watcher',
                data: newState ? 'enable' : 'disable'
            }));
        }
    }, [promptWatcherEnabled]);

    // Initialize terminal and WebSocket connection
    useEffect(() => {
        if (!terminalRef.current) return;

        // Create xterm instance with dark theme
        const term = new XTerm({
            cursorBlink: true,
            fontSize: 14,
            fontFamily: 'JetBrains Mono, Menlo, Monaco, "Courier New", monospace',
            theme: {
                background: '#0f172a',  // slate-900
                foreground: '#e2e8f0',  // slate-200
                cursor: '#3b82f6',       // blue-500
                cursorAccent: '#0f172a',
                selectionBackground: '#3b82f680',
                black: '#1e293b',
                red: '#ef4444',
                green: '#22c55e',
                yellow: '#eab308',
                blue: '#3b82f6',
                magenta: '#a855f7',
                cyan: '#06b6d4',
                white: '#f1f5f9',
                brightBlack: '#475569',
                brightRed: '#f87171',
                brightGreen: '#4ade80',
                brightYellow: '#facc15',
                brightBlue: '#60a5fa',
                brightMagenta: '#c084fc',
                brightCyan: '#22d3ee',
                brightWhite: '#ffffff',
            },
            allowProposedApi: true,
        });

        // Add fit addon for responsive sizing
        const fitAddon = new FitAddon();
        term.loadAddon(fitAddon);
        fitAddonRef.current = fitAddon;

        // Add web links addon for clickable URLs
        const webLinksAddon = new WebLinksAddon();
        term.loadAddon(webLinksAddon);

        // Open terminal in container
        term.open(terminalRef.current);
        xtermRef.current = term;

        // Fit terminal to container
        setTimeout(() => fitAddon.fit(), 0);

        // Connect to PTY WebSocket
        const wsUrl = getWebSocketUrl();
        const ws = new WebSocket(wsUrl);
        wsRef.current = ws;

        ws.onopen = () => {
            setIsConnected(true);
            term.writeln('\x1b[32m✓ Connected to terminal\x1b[0m\r\n');
            onConnect?.();

            // Send initial resize
            const { rows, cols } = term;
            ws.send(JSON.stringify({ type: 'resize', rows, cols }));
        };

        ws.onmessage = (event) => {
            // Write received data to terminal
            term.write(event.data);
        };

        ws.onclose = () => {
            setIsConnected(false);
            term.writeln('\r\n\x1b[31m✗ Disconnected from terminal\x1b[0m');
            onDisconnect?.();
        };

        ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            term.writeln('\r\n\x1b[31m✗ Connection error\x1b[0m');
        };

        // Send user input to PTY
        term.onData((data) => {
            if (ws.readyState === WebSocket.OPEN) {
                ws.send(JSON.stringify({ type: 'input', data }));
            }
        });

        // Handle terminal resize
        const handleResize = () => {
            fitAddon.fit();
            if (ws.readyState === WebSocket.OPEN) {
                const { rows, cols } = term;
                ws.send(JSON.stringify({ type: 'resize', rows, cols }));
            }
        };

        // Debounced resize handler
        let resizeTimeout: ReturnType<typeof setTimeout>;
        const debouncedResize = () => {
            clearTimeout(resizeTimeout);
            resizeTimeout = setTimeout(handleResize, 100);
        };

        window.addEventListener('resize', debouncedResize);

        // Cleanup
        return () => {
            window.removeEventListener('resize', debouncedResize);
            clearTimeout(resizeTimeout);
            ws.close();
            term.dispose();
        };
    }, [getWebSocketUrl, onConnect, onDisconnect]);

    return (
        <div className={`flex flex-col h-full ${className}`}>
            {/* Terminal header with controls */}
            <div className="flex items-center justify-between px-4 py-2 bg-slate-800 border-b border-slate-700">
                <div className="flex items-center gap-3">
                    <div className="flex gap-1.5">
                        <div className="w-3 h-3 rounded-full bg-red-500" />
                        <div className="w-3 h-3 rounded-full bg-yellow-500" />
                        <div className="w-3 h-3 rounded-full bg-green-500" />
                    </div>
                    <span className="text-sm text-slate-400 font-mono">Terminal</span>
                    <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`} />
                </div>
                
                {/* Prompt Watcher Toggle - Task 2.3 */}
                <button
                    onClick={togglePromptWatcher}
                    className={`flex items-center gap-2 px-3 py-1 rounded text-sm transition-colors ${
                        promptWatcherEnabled
                            ? 'bg-blue-600 text-white'
                            : 'bg-slate-700 text-slate-400 hover:bg-slate-600'
                    }`}
                    title="Auto-respond to confirmation prompts (y/n)"
                    data-testid="prompt-watcher-toggle"
                >
                    {promptWatcherEnabled ? <Eye size={16} /> : <EyeOff size={16} />}
                    <span>Prompt Watcher</span>
                </button>
            </div>
            
            {/* Terminal container */}
            <div
                ref={terminalRef}
                className="flex-1 bg-slate-900 p-2"
                data-testid="terminal-container"
            />
        </div>
    );
};

export default Terminal;
