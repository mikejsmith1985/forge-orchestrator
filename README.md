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
