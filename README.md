# Forge Orchestrator

Forge Orchestrator is a terminal-centric developer tool that puts a real PTY (pseudo-terminal) at the center of your workflow. It integrates LLM capabilities directly into your development environment, providing intelligent command generation, token usage tracking, and automated workflows—all while keeping you in control of your shell.

## Features

-   **🖥️ Integrated Terminal**: Full PTY terminal powered by xterm.js with WebSocket streaming. Type commands, see real output, and maintain persistent shell sessions.
-   **🤖 Prompt Watcher**: Automatically responds to confirmation prompts (y/n) when enabled, perfect for unattended automation.
-   **🧠 Architect**: Leverage LLMs to brainstorm and plan tasks with live token estimation.
-   **📊 Ledger**: Track token usage and costs across different LLM providers with Primary Cost Unit support (TOKEN or PROMPT billing).
-   **⚡ Commands**: Save and execute frequently-used shell commands with one click.
-   **🔀 Flows**: Visually design automation workflows with two distinct node types:
    - **Shell Command Nodes** (Zero-Token): Execute local scripts without consuming LLM budget
    - **LLM Prompt Nodes** (Premium): Run AI-powered tasks with confirmation gating
-   **🔐 Secure Keyring**: API keys are encrypted and stored in your OS native keyring (macOS Keychain, Windows Credential Manager, Linux Secret Service).

## Quick Start

### Prerequisites

-   **Go**: Version 1.24 or higher
-   **Node.js**: Version 20 or higher

### Installation

```bash
# Clone the repository
git clone https://github.com/mikejsmith1985/forge-orchestrator.git
cd forge-orchestrator

# Install backend dependencies
go mod download

# Install frontend dependencies
cd frontend && npm install && cd ..

# Build the application
go build -o forge-orchestrator .
```

### Running

```bash
# Start the backend
./forge-orchestrator

# In another terminal, start the frontend (for development)
cd frontend && npm run dev
```

The application opens with the **Terminal** as the default view, establishing the app's identity as a terminal-first developer tool.

## Usage

### Terminal View (Default)

The Terminal is your primary workspace:
- Full shell access with PTY streaming
- Connection status indicator (green = connected)
- **Prompt Watcher** toggle to auto-respond to y/n confirmations
- Commands typed here execute in your actual shell

### Flow Editor

Create visual automation workflows:

1. **Drag nodes** from the sidebar onto the canvas:
   - **Shell Command** (⚡ Zero-Token): Free local execution
   - **LLM Prompt** (💎 Premium): AI-powered with token cost

2. **Configure nodes** by clicking them:
   - Enter your command or prompt string
   - Shell nodes execute immediately
   - LLM nodes show a confirmation modal before consuming budget

3. **Connect nodes** to define execution order

4. **Execute** the flow to run all nodes in sequence

### Token Economy

Forge Orchestrator tracks your AI spending with precision:

- **TOKEN-based billing**: For traditional LLM providers (OpenAI, Anthropic)
- **PROMPT-based billing**: For per-request pricing models
- **Dynamic Budget Meter**: Shows remaining budget in the correct currency
- **Ledger**: Full history with cost breakdown by Primary Cost Unit

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     React Frontend                          │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────────────┐│
│  │Terminal │  │ Flows   │  │ Ledger  │  │   Architect     ││
│  │(xterm)  │  │ Editor  │  │ View    │  │   View          ││
│  └────┬────┘  └────┬────┘  └────┬────┘  └────────┬────────┘│
│       │            │            │                 │         │
│       └────────────┴────────────┴─────────────────┘         │
│                         │ WebSocket + REST                  │
└─────────────────────────┼───────────────────────────────────┘
                          │
┌─────────────────────────┼───────────────────────────────────┐
│                     Go Backend                              │
│  ┌──────────────┐  ┌────────────┐  ┌───────────────────────┐│
│  │ PTY Manager  │  │ WebSocket  │  │ REST API Handlers     ││
│  │ (pty_manager)│  │ Hub        │  │ /api/*, /ws/*         ││
│  └──────┬───────┘  └─────┬──────┘  └───────────┬───────────┘│
│         │                │                      │            │
│  ┌──────┴────────────────┴──────────────────────┴──────────┐│
│  │              Core Services                               ││
│  │  Executor │ LLM Gateway │ Tokenizer │ Keyring │ Ledger  ││
│  └──────────────────────────────────────────────────────────┘│
│                          │                                   │
│  ┌───────────────────────┴──────────────────────────────────┐│
│  │                    SQLite Database                       ││
│  └──────────────────────────────────────────────────────────┘│
└──────────────────────────────────────────────────────────────┘
```

For detailed architecture documentation, see [docs/architecture.md](docs/architecture.md).

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `FORGE_ALLOWED_ORIGINS` | Comma-separated list of allowed CORS origins | `http(s)://localhost:*` |
| `FORGE_TLS_CERT` | Path to TLS certificate file | (none) |
| `FORGE_TLS_KEY` | Path to TLS private key file | (none) |

### TLS/HTTPS

```bash
# Production with certificates
export FORGE_TLS_CERT=/path/to/cert.crt
export FORGE_TLS_KEY=/path/to/key.pem
./forge-orchestrator

# Development with self-signed certificate
./forge-orchestrator --dev-tls

# HTTP only (not recommended for production)
./forge-orchestrator
```

## API Endpoints

### Terminal/PTY
- `WS /ws/pty` - WebSocket for PTY streaming
- `POST /api/command/execute` - Inject command into active PTY session

### Core
- `GET /api/health` - Health check
- `POST /api/execute` - Execute command via Executor interface
- `POST /api/tokens/estimate` - Estimate token count

### Flows
- `GET/POST /api/flows` - List/create flows
- `GET/PUT/DELETE /api/flows/{id}` - CRUD operations
- `POST /api/flows/{id}/execute` - Execute a flow

### Ledger
- `GET/POST /api/ledger` - Token usage records
- `GET /api/ledger/optimizations` - Cost optimization suggestions

### Keys
- `POST /api/keys` - Save API key to keyring
- `GET /api/keys/status` - Check which providers are configured

## Testing

```bash
# Backend tests
go test ./...

# Frontend tests
cd frontend
npm run test           # Unit tests
npm run test:e2e       # Playwright E2E tests
```

## License

MIT License - see [LICENSE](LICENSE) for details.
