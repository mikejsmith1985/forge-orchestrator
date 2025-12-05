# Forge Orchestrator

Forge Orchestrator is a powerful developer tool designed to streamline your workflow by integrating LLM capabilities directly into your command-line and development environment. It acts as a bridge between your intent and execution, providing intelligent command generation, token usage tracking, and automated workflows.

## Features

-   **Architect**: Leverage LLMs to generate complex terminal commands from natural language descriptions.
-   **Ledger**: Keep track of your token usage and costs across different LLM providers.
-   **Commands**: Execute generated commands safely and maintain a history of your actions.
-   **Flows**: Visually design and execute complex automation workflows.
-   **Keyring**: Securely manage your API keys for various services.

## Setup Guide

### Prerequisites

-   **Go**: Version 1.23 or higher.
-   **Node.js**: Version 20 or higher.

### Installation

1.  Clone the repository:
    ```bash
    git clone https://github.com/mikejsmith1985/forge-orchestrator.git
    cd forge-orchestrator
    ```

2.  Install Backend Dependencies:
    ```bash
    go mod download
    ```

3.  Install Frontend Dependencies:
    ```bash
    cd frontend
    npm install
    cd ..
    ```

### Running Locally

To start the application, you need to run both the backend and the frontend.

1.  **Backend**:
    ```bash
    go run main.go
    ```
    The backend server will start on `http://localhost:8080`.

2.  **Frontend**:
    ```bash
    cd frontend
    npm run dev
    ```
    The frontend development server will start on `http://localhost:5173`.

## Architecture

Forge Orchestrator follows a modern client-server architecture. The frontend is built with React and TypeScript, communicating with a Go backend via a RESTful API. The backend handles business logic, interacts with a SQLite database for persistence, and manages integrations with LLM providers through a dedicated Gateway.

For a detailed deep dive into the system design, please refer to [docs/architecture.md](docs/architecture.md).

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `FORGE_ALLOWED_ORIGINS` | Comma-separated list of allowed origins for CORS | `http(s)://localhost:8080,http(s)://localhost:5173,http(s)://127.0.0.1:8080,http(s)://127.0.0.1:5173` |
| `FORGE_TLS_CERT` | Path to TLS certificate file for HTTPS | (none) |
| `FORGE_TLS_KEY` | Path to TLS private key file for HTTPS | (none) |

### TLS/HTTPS Configuration

Forge Orchestrator supports three TLS modes:

#### 1. Production HTTPS (Recommended)

Use your own TLS certificates for production deployments:

```bash
# Set certificate and key paths
export FORGE_TLS_CERT=/path/to/your/certificate.crt
export FORGE_TLS_KEY=/path/to/your/private.key

# Start the server
./forge-orchestrator
```

The server will log: `🔒 Starting HTTPS server on :8080`

#### 2. Development HTTPS (Self-Signed)

For local development with HTTPS, use the `--dev-tls` flag to auto-generate a self-signed certificate:

```bash
./forge-orchestrator --dev-tls
```

This will:
- Generate a self-signed certificate valid for `localhost` and `127.0.0.1`
- Save the certificate to `.forge/certs/` for inspection
- Certificate is valid for 1 year
- **⚠️ Not suitable for production use**

The server will log:
```
⚠️  Generating self-signed certificate for development
⚠️  This is NOT suitable for production use!
🔒 Starting HTTPS server on :8080 (self-signed)
```

#### 3. HTTP Mode (No TLS)

If no TLS configuration is provided, the server runs in HTTP mode:

```bash
./forge-orchestrator
```

The server will log: `⚠️  Running HTTP server (no TLS) - not recommended for production`

### Generating Production Certificates

For production, you can obtain certificates from:

1. **Let's Encrypt** (free, automated):
   ```bash
   certbot certonly --standalone -d yourdomain.com
   export FORGE_TLS_CERT=/etc/letsencrypt/live/yourdomain.com/fullchain.pem
   export FORGE_TLS_KEY=/etc/letsencrypt/live/yourdomain.com/privkey.pem
   ```

2. **Self-Signed (for internal use)**:
   ```bash
   openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
   export FORGE_TLS_CERT=cert.pem
   export FORGE_TLS_KEY=key.pem
   ```

#### CORS Configuration

By default, the server only allows connections from localhost development servers. For production deployments, you should set the `FORGE_ALLOWED_ORIGINS` environment variable to your domain:

```bash
# Development (default)
FORGE_ALLOWED_ORIGINS=http://localhost:8080,http://localhost:5173

# Production example
FORGE_ALLOWED_ORIGINS=https://myapp.example.com

# Multiple production origins
FORGE_ALLOWED_ORIGINS=https://myapp.example.com,https://admin.example.com
```

The server logs allowed origins on startup for verification.
