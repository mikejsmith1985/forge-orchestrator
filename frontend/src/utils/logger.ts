/**
 * Application logger for capturing logs to include in feedback reports.
 * Wraps console methods to store logs in memory while preserving original behavior.
 */

const MAX_LOGS = 200;
const logs: string[] = [];

const originalLog = console.log;
const originalWarn = console.warn;
const originalError = console.error;

function formatLog(level: string, args: unknown[]): string {
    const timestamp = new Date().toISOString();
    const message = args.map(arg => {
        if (typeof arg === 'object') {
            try {
                return JSON.stringify(arg);
            } catch {
                return '[Circular/Object]';
            }
        }
        return String(arg);
    }).join(' ');
    return `[${timestamp}] [${level}] ${message}`;
}

function addLog(level: string, args: unknown[]): void {
    const logEntry = formatLog(level, args);
    logs.push(logEntry);
    if (logs.length > MAX_LOGS) {
        logs.shift();
    }
}

// Override console methods
console.log = (...args: unknown[]) => {
    addLog('INFO', args);
    originalLog.apply(console, args);
};

console.warn = (...args: unknown[]) => {
    addLog('WARN', args);
    originalWarn.apply(console, args);
};

console.error = (...args: unknown[]) => {
    addLog('ERROR', args);
    originalError.apply(console, args);
};

/**
 * Get all captured logs as a single string.
 */
export function getLogs(): string {
    return logs.join('\n');
}

/**
 * Clear all captured logs.
 */
export function clearLogs(): void {
    logs.length = 0;
}

/**
 * Structured logger for specific application components.
 */
export const logger = {
    feedback: (action: string, data: Record<string, unknown> = {}) => {
        console.log(`[Feedback] ${action}`, data);
    },
    api: (action: string, data: Record<string, unknown> = {}) => {
        console.log(`[API] ${action}`, data);
    },
    flow: (action: string, data: Record<string, unknown> = {}) => {
        console.log(`[Flow] ${action}`, data);
    },
    toast: (action: string, data: Record<string, unknown> = {}) => {
        console.log(`[Toast] ${action}`, data);
    },
};
