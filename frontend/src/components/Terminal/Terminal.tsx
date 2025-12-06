import React, { useEffect, useRef, useState, useCallback } from 'react';
import { Terminal as XTerm } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import { WebLinksAddon } from '@xterm/addon-web-links';
import { SearchAddon } from '@xterm/addon-search';
import '@xterm/xterm/css/xterm.css';
import { Eye, EyeOff, ArrowDownToLine, Plus, Minus } from 'lucide-react';

/**
 * Enhanced Terminal Component with forge-terminal features
 * 
 * Features:
 * - Advanced CLI prompt detection and auto-respond
 * - Auto-reconnection with exponential backoff
 * - Search functionality
 * - Scroll-to-bottom button
 * - Better error handling and connection status
 */

// ============================================================================
// CLI Prompt Detection for Auto-Respond Feature
// ============================================================================

function stripAnsi(text: string): string {
    // eslint-disable-next-line no-control-regex
    return text.replace(/\x1b\[[0-9;]*[a-zA-Z]/g, '');
}

// Menu-style prompts where an option is already selected (just press Enter)
const MENU_SELECTION_PATTERNS = [
    /[›❯>]\s*1\.\s*Yes\b/i,
    /[›❯>]\s*Yes\b/i,
    /[›❯>]\s*Run\s+this\s+command/i,
    /[●◉✓✔]\s*Yes\b/i,
];

const MENU_CONTEXT_PATTERNS = [
    /Confirm with number keys or.*Enter/i,
    /use.*arrow.*keys.*select/i,
    /↑↓.*keys.*Enter/i,
    /Do you want to run this command\??/i,
    /Do you want to run\??/i,
    /Cancel with Esc/i,
];

const YN_PROMPT_PATTERNS = [
    /\(y\/n\)\s*$/i,
    /\[Y\/n\]\s*$/i,
    /\[y\/N\]\s*$/i,
    /\(yes\/no\)\s*$/i,
    /\[yes\/no\]\s*$/i,
    /\?\s*\(y\/n\)\s*$/i,
    /\?\s*\[Y\/n\]\s*$/i,
    /\?\s*\[y\/N\]\s*$/i,
    /\?\s*›?\s*\(Y\/n\)\s*$/i,
    /Are you sure.*\?\s*$/i,
];

const QUESTION_PATTERNS = [
    /Do you want to run this command\?/i,
    /Do you want to proceed\?/i,
    /Do you want to continue\?/i,
    /Would you like to proceed\?/i,
    /Proceed\?/i,
    /Continue\?/i,
    /Run this command\?/i,
];

const TUI_FRAME_INDICATORS = [
    /[╭╮╯╰│─┌┐└┘├┤┬┴┼]/,
    /Remaining requests:\s*[\d.]+%/i,
    /Ctrl\+c\s+Exit/i,
];

function detectMenuPrompt(cleanText: string): { detected: boolean; confidence: 'high' | 'medium' | 'low' } {
    const hasYesSelected = MENU_SELECTION_PATTERNS.some(p => p.test(cleanText));
    
    if (!hasYesSelected) {
        return { detected: false, confidence: 'low' };
    }
    
    const hasMenuContext = MENU_CONTEXT_PATTERNS.some(p => p.test(cleanText));
    const hasQuestion = QUESTION_PATTERNS.some(p => p.test(cleanText));
    const hasTuiFrame = TUI_FRAME_INDICATORS.some(p => p.test(cleanText));
    
    if (hasYesSelected && (hasMenuContext || hasTuiFrame)) {
        return { detected: true, confidence: 'high' };
    }
    
    if (hasYesSelected && hasQuestion) {
        return { detected: true, confidence: 'medium' };
    }
    
    if (hasYesSelected) {
        return { detected: true, confidence: 'low' };
    }
    
    return { detected: false, confidence: 'low' };
}

function detectYnPrompt(cleanText: string): { detected: boolean } {
    const lines = cleanText.split(/[\r\n]/).filter(l => l.trim());
    const lastLines = lines.slice(-3).join('\n');
    const hasYnPrompt = YN_PROMPT_PATTERNS.some(p => p.test(lastLines));
    return { detected: hasYnPrompt };
}

function detectCliPrompt(text: string): { waiting: boolean; responseType: 'enter' | 'y-enter' | null; confidence: string } {
    if (!text || text.length < 10) {
        return { waiting: false, responseType: null, confidence: 'none' };
    }
    
    const cleanText = stripAnsi(text);
    const bufferToCheck = cleanText.slice(-2000);
    
    const menuResult = detectMenuPrompt(bufferToCheck);
    if (menuResult.detected && menuResult.confidence !== 'low') {
        return { 
            waiting: true, 
            responseType: 'enter', 
            confidence: menuResult.confidence 
        };
    }
    
    const ynResult = detectYnPrompt(bufferToCheck);
    if (ynResult.detected) {
        return { 
            waiting: true, 
            responseType: 'y-enter', 
            confidence: 'high' 
        };
    }
    
    if (menuResult.detected && menuResult.confidence === 'low') {
        return { 
            waiting: true, 
            responseType: 'enter', 
            confidence: 'low' 
        };
    }
    
    return { waiting: false, responseType: null, confidence: 'none' };
}

interface TerminalProps {
    className?: string;
    onConnect?: () => void;
    onDisconnect?: () => void;
}

export const Terminal: React.FC<TerminalProps> = ({
    className = '',
    onConnect,
    onDisconnect
}) => {
    const terminalRef = useRef<HTMLDivElement>(null);
    const containerRef = useRef<HTMLDivElement>(null);
    const xtermRef = useRef<XTerm | null>(null);
    const fitAddonRef = useRef<FitAddon | null>(null);
    const searchAddonRef = useRef<SearchAddon | null>(null);
    const wsRef = useRef<WebSocket | null>(null);
    const connectFnRef = useRef<(() => void) | null>(null);
    
    const [isConnected, setIsConnected] = useState(false);
    const [reconnecting, setReconnecting] = useState(false);
    const [promptWatcherEnabled, setPromptWatcherEnabled] = useState(false);
    const [showScrollButton, setShowScrollButton] = useState(false);
    
    const reconnectAttemptsRef = useRef(0);
    const reconnectTimeoutRef = useRef<number | null>(null);
    const maxReconnectAttempts = 5;
    const lastOutputRef = useRef('');
    const waitingCheckTimeoutRef = useRef<number | null>(null);

    const getWebSocketUrl = useCallback(() => {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const host = window.location.hostname;
        // In development, backend is on port 9000, but vite proxy handles it
        const port = window.location.port === '5173' ? '5173' : window.location.port;
        return `${protocol}//${host}:${port}/ws/pty`;
    }, []);

    const [fontSize, setFontSize] = useState(() => {
        const saved = localStorage.getItem('forge_terminal_font_size');
        return saved ? parseInt(saved, 10) : 14;
    });

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

    const changeFontSize = useCallback((delta: number) => {
        setFontSize(prev => {
            const newSize = Math.max(8, Math.min(24, prev + delta));
            localStorage.setItem('forge_terminal_font_size', newSize.toString());
            if (xtermRef.current) {
                xtermRef.current.options.fontSize = newSize;
                if (fitAddonRef.current) {
                    setTimeout(() => fitAddonRef.current?.fit(), 0);
                }
            }
            return newSize;
        });
    }, []);
    
    const handleScrollToBottom = useCallback(() => {
        if (xtermRef.current) {
            xtermRef.current.scrollToBottom();
            setShowScrollButton(false);
        }
    }, []);

    // Initialize terminal and WebSocket connection
    useEffect(() => {
        if (!terminalRef.current) return;

        const term = new XTerm({
            cursorBlink: true,
            fontSize: fontSize,
            fontFamily: 'JetBrains Mono, Menlo, Monaco, "Courier New", monospace',
            theme: {
                background: '#0f172a',
                foreground: '#e2e8f0',
                cursor: '#3b82f6',
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
            scrollback: 5000,
        });

        const fitAddon = new FitAddon();
        term.loadAddon(fitAddon);
        fitAddonRef.current = fitAddon;

        const webLinksAddon = new WebLinksAddon();
        term.loadAddon(webLinksAddon);
        
        const searchAddon = new SearchAddon();
        term.loadAddon(searchAddon);
        searchAddonRef.current = searchAddon;

        term.open(terminalRef.current);
        xtermRef.current = term;

        setTimeout(() => fitAddon.fit(), 0);

        // WebSocket connection function
        const connectWebSocket = () => {
            const wsUrl = getWebSocketUrl();
            const ws = new WebSocket(wsUrl);
            wsRef.current = ws;

            ws.onopen = () => {
                console.log('[Terminal] WebSocket connected');
                reconnectAttemptsRef.current = 0;
                setReconnecting(false);
                setIsConnected(true);
                
                const reconnectLabel = reconnectAttemptsRef.current > 0 ? ' [Reconnected]' : '';
                term.write(`\r\n\x1b[38;2;249;115;22m[Forge Terminal]\x1b[0m Connected${reconnectLabel}.\r\n\r\n`);
                onConnect?.();

                const { rows, cols } = term;
                ws.send(JSON.stringify({ type: 'resize', rows, cols }));
            };

            ws.onmessage = (event) => {
                let textData = '';
                if (event.data instanceof ArrayBuffer) {
                    const data = new Uint8Array(event.data);
                    term.write(data);
                    textData = new TextDecoder().decode(data);
                } else {
                    term.write(event.data);
                    textData = event.data;
                }
                
                // Accumulate output for prompt detection
                lastOutputRef.current = (lastOutputRef.current + textData).slice(-3000);
                
                // Debounce waiting check
                if (waitingCheckTimeoutRef.current) {
                    clearTimeout(waitingCheckTimeoutRef.current);
                }
                waitingCheckTimeoutRef.current = setTimeout(() => {
                    const { waiting, responseType, confidence } = detectCliPrompt(lastOutputRef.current);
                    
                    const shouldAutoRespond = waiting && 
                        promptWatcherEnabled && 
                        ws.readyState === WebSocket.OPEN &&
                        (confidence === 'high' || confidence === 'medium');
                    
                    if (shouldAutoRespond) {
                        console.log('[Terminal] Auto-responding to CLI prompt', { responseType, confidence });
                        
                        if (responseType === 'enter') {
                            ws.send('\r');
                        } else {
                            ws.send('y\r');
                        }
                        
                        lastOutputRef.current = '';
                    }
                }, 500);
            };

            ws.onerror = (error) => {
                console.error('[Terminal] WebSocket error:', error);
                term.write('\r\n\x1b[1;31m[Error]\x1b[0m Connection error.\r\n');
            };

            ws.onclose = (event) => {
                console.log('[Terminal] WebSocket closed', { code: event.code, reason: event.reason });
                
                setIsConnected(false);
                onDisconnect?.();
                
                let disconnectMessage = 'Terminal session ended.';
                let messageColor = '1;33';
                let shouldReconnect = false;
                
                switch (event.code) {
                    case 1000:
                        disconnectMessage = 'Session closed normally.';
                        break;
                    case 1001:
                    case 1006:
                    case 1011:
                    case 1012:
                    case 1013:
                        disconnectMessage = 'Connection lost. Attempting to reconnect...';
                        shouldReconnect = true;
                        break;
                    case 4000:
                        disconnectMessage = 'Shell process exited.';
                        break;
                    default:
                        if (event.reason) {
                            disconnectMessage = event.reason;
                        }
                        shouldReconnect = true;
                }
                
                if (xtermRef.current) {
                    term.write(`\r\n\x1b[${messageColor}m[Disconnected]\x1b[0m ${disconnectMessage}\r\n`);
                }
                
                // Attempt reconnection with exponential backoff
                if (shouldReconnect && reconnectAttemptsRef.current < maxReconnectAttempts) {
                    const delay = Math.min(1000 * Math.pow(2, reconnectAttemptsRef.current), 30000);
                    reconnectAttemptsRef.current += 1;
                    setReconnecting(true);
                    
                    console.log('[Terminal] Scheduling reconnection', { 
                        attempt: reconnectAttemptsRef.current, 
                        delay 
                    });
                    
                    reconnectTimeoutRef.current = setTimeout(() => {
                        if (xtermRef.current && wsRef.current === ws) {
                            console.log('[Terminal] Attempting reconnection...', { attempt: reconnectAttemptsRef.current });
                            if (xtermRef.current) {
                                term.write(`\x1b[1;33m[Reconnecting...]\x1b[0m Attempt ${reconnectAttemptsRef.current}/${maxReconnectAttempts}\r\n`);
                            }
                            connectWebSocket();
                        }
                    }, delay);
                } else if (reconnectAttemptsRef.current >= maxReconnectAttempts) {
                    setReconnecting(false);
                    if (xtermRef.current) {
                        term.write(`\x1b[1;31m[Connection Failed]\x1b[0m Max reconnection attempts reached.\r\n`);
                    }
                }
            };

            term.onData((data) => {
                if (ws.readyState === WebSocket.OPEN) {
                    ws.send(data);
                }
            });

            term.onResize(({ cols, rows }) => {
                if (ws.readyState === WebSocket.OPEN) {
                    ws.send(JSON.stringify({ type: 'resize', cols, rows }));
                }
            });
            
            // Track scroll position
            const viewport = terminalRef.current?.querySelector('.xterm-viewport');
            if (viewport) {
                const checkScroll = () => {
                    const isAtBottom = viewport.scrollHeight - viewport.scrollTop - viewport.clientHeight < 50;
                    setShowScrollButton(!isAtBottom);
                };
                viewport.addEventListener('scroll', checkScroll);
            }

            return ws;
        };

        connectFnRef.current = connectWebSocket;
        connectWebSocket();

        // Handle window resize
        const debouncedFit = (() => {
            let timeoutId: number;
            return () => {
                clearTimeout(timeoutId);
                timeoutId = window.setTimeout(() => {
                    if (fitAddonRef.current) {
                        fitAddonRef.current.fit();
                    }
                }, 100);
            };
        })();

        window.addEventListener('resize', debouncedFit);

        const resizeObserver = new ResizeObserver(() => {
            debouncedFit();
        });
        resizeObserver.observe(terminalRef.current);

        return () => {
            window.removeEventListener('resize', debouncedFit);
            resizeObserver.disconnect();
            if (waitingCheckTimeoutRef.current) {
                clearTimeout(waitingCheckTimeoutRef.current);
            }
            if (reconnectTimeoutRef.current) {
                clearTimeout(reconnectTimeoutRef.current);
            }
            if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
                wsRef.current.onclose = null;
                wsRef.current.close();
            }
            xtermRef.current = null;
            term.dispose();
        };
    }, [getWebSocketUrl, onConnect, onDisconnect, promptWatcherEnabled, fontSize]);

    return (
        <div ref={containerRef} className={`flex flex-col h-full ${className}`} style={{ position: 'relative' }}>
            {/* Connection Status Overlay */}
            {!isConnected && (
                <div style={{
                    position: 'absolute',
                    top: 0,
                    left: 0,
                    right: 0,
                    bottom: 0,
                    background: 'rgba(15, 23, 42, 0.95)',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    zIndex: 10,
                }}>
                    <div style={{
                        textAlign: 'center',
                        padding: '20px',
                    }}>
                        {reconnecting ? (
                            <>
                                <div style={{
                                    width: '40px',
                                    height: '40px',
                                    border: '3px solid #3b82f6',
                                    borderTopColor: 'transparent',
                                    borderRadius: '50%',
                                    animation: 'spin 1s linear infinite',
                                    margin: '0 auto 15px',
                                }} />
                                <span style={{ color: '#e2e8f0' }}>
                                    Reconnecting... (Attempt {reconnectAttemptsRef.current}/{maxReconnectAttempts})
                                </span>
                            </>
                        ) : (
                            <>
                                <span style={{ color: '#ef4444', fontWeight: 600, display: 'block', marginBottom: '15px' }}>
                                    ⚠ Disconnected
                                </span>
                                <button 
                                    className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors"
                                    onClick={() => {
                                        if (xtermRef.current) {
                                            xtermRef.current.clear();
                                        }
                                        reconnectAttemptsRef.current = 0;
                                        if (connectFnRef.current) {
                                            connectFnRef.current();
                                        }
                                    }}
                                >
                                    Reconnect Terminal
                                </button>
                                <small style={{ display: 'block', marginTop: '10px', color: '#94a3b8' }}>
                                    The terminal connection was lost. Click to reconnect.
                                </small>
                            </>
                        )}
                    </div>
                </div>
            )}
            
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
                
                <div className="flex items-center gap-2">
                    {/* Font Size Controls */}
                    <div className="flex items-center gap-1 bg-slate-700 rounded px-2 py-1">
                        <button
                            onClick={() => changeFontSize(-1)}
                            className="text-slate-400 hover:text-white transition-colors p-1"
                            title="Decrease font size"
                            data-testid="font-size-decrease"
                        >
                            <Minus size={14} />
                        </button>
                        <span className="text-xs text-slate-400 min-w-[2rem] text-center" data-testid="font-size-display">
                            {fontSize}px
                        </span>
                        <button
                            onClick={() => changeFontSize(1)}
                            className="text-slate-400 hover:text-white transition-colors p-1"
                            title="Increase font size"
                            data-testid="font-size-increase"
                        >
                            <Plus size={14} />
                        </button>
                    </div>
                    
                    {/* Prompt Watcher Toggle */}
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
                        <span>Auto-Respond</span>
                    </button>
                </div>
            </div>
            
            {/* Terminal container */}
            <div
                ref={terminalRef}
                className="flex-1 bg-slate-900 p-2"
                data-testid="terminal-container"
                onClick={() => {
                    if (xtermRef.current) {
                        xtermRef.current.focus();
                    }
                }}
                style={{
                    width: '100%',
                    height: '100%',
                    backgroundColor: '#0f172a',
                    cursor: 'text',
                }}
            />
            
            {/* Scroll to bottom button */}
            {showScrollButton && isConnected && (
                <button
                    onClick={handleScrollToBottom}
                    title="Scroll to bottom (Ctrl+End)"
                    aria-label="Scroll to bottom"
                    style={{
                        position: 'absolute',
                        right: '20px',
                        bottom: '20px',
                        width: '40px',
                        height: '40px',
                        borderRadius: '50%',
                        background: '#3b82f6',
                        border: 'none',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        cursor: 'pointer',
                        boxShadow: '0 4px 12px rgba(59, 130, 246, 0.5)',
                        transition: 'all 0.2s',
                        zIndex: 5,
                    }}
                    onMouseEnter={(e) => {
                        e.currentTarget.style.transform = 'scale(1.1)';
                        e.currentTarget.style.background = '#2563eb';
                    }}
                    onMouseLeave={(e) => {
                        e.currentTarget.style.transform = 'scale(1)';
                        e.currentTarget.style.background = '#3b82f6';
                    }}
                >
                    <ArrowDownToLine size={20} color="white" />
                </button>
            )}
        </div>
    );
};

export default Terminal;
