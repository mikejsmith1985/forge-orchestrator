# System Architecture

This document provides a technical overview of the Forge Orchestrator's architecture, including system design, data flow, and key components.

## System Design

The application is structured as a Backend-for-Frontend (BFF) architecture.

```mermaid
graph TD
    User[User] --> Frontend[Frontend (React/TS)]
    Frontend --> Backend[Backend API (Go)]
    Backend --> SQLite[(SQLite Database)]
    Backend --> LLM[LLM Gateway]
    LLM --> OpenAI[OpenAI API]
    LLM --> Anthropic[Anthropic API]
```

## Data Flow

### Command Generation & Execution

1.  **Request**: User enters a natural language prompt in the "Architect" UI.
2.  **Processing**: The Frontend sends the prompt to the Backend (`POST /api/architect/generate`).
3.  **LLM Interaction**: The Backend's LLM Gateway constructs a prompt and calls the configured LLM provider.
4.  **Response**: The LLM returns a suggested command.
5.  **Storage**: The command is stored in the `command_cards` table in SQLite.
6.  **Display**: The Frontend receives the command and displays it as a card.
7.  **Execution**: User clicks "Run". Frontend calls `POST /api/commands/run`.
8.  **Ledger**: Token usage for the generation is recorded in the `token_ledger`.

## Key Components

### Backend (`internal/`)

-   **`internal/server`**: Contains the HTTP server setup, route definitions, and API handlers. This is the entry point for all client requests.
-   **`internal/llm`**: The LLM Gateway. It abstracts the differences between providers (OpenAI, Anthropic) and handles API communication.
-   **`internal/flows`**: The execution engine for workflows. It parses flow definitions and orchestrates the execution of steps.
-   **`internal/security`**: The Keyring component. It handles the secure storage and retrieval of sensitive API keys using the operating system's keyring service where possible, or encrypted storage.

### Frontend (`frontend/src/`)

-   **`components/`**: Reusable UI components.
-   **`pages/`**: Top-level page components corresponding to routes.
-   **`api/`**: API client functions for communicating with the backend.

## Database Schema

The application uses SQLite for data persistence.

### `token_ledger`
Tracks token consumption and costs for every LLM interaction.
-   `id`: Integer (Primary Key)
-   `timestamp`: Datetime
-   `flow_id`: String
-   `model_used`: String
-   `agent_role`: String
-   `prompt_hash`: String
-   `input_tokens`: Integer
-   `output_tokens`: Integer
-   `total_cost_usd`: Float
-   `latency_ms`: Integer
-   `status`: String ('SUCCESS', 'FAILED', 'TIMEOUT')
-   `error_message`: String

### `command_cards`
Stores reusable terminal commands.
-   `id`: Integer (Primary Key)
-   `name`: String
-   `command`: String
-   `description`: String

### `forge_flows`
Stores workflow definitions.
-   `id`: Integer (Primary Key)
-   `name`: String
-   `description`: String
-   `data`: JSON (Serialized node-edge structure)
-   `status`: String ('draft', 'active', 'archived')
-   `created_at`: Datetime
-   `updated_at`: Datetime

### `user_secrets`
Stores encrypted API keys.
-   `key_name`: String (Primary Key)
-   `encrypted_value`: Blob
-   `created_at`: Datetime
