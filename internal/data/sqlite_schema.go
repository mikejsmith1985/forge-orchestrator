package data

// SQLiteSchema defines the SQL commands necessary to initialize the single-file SQLite database.
const SQLiteSchema = `
-- Table 1: token_ledger
-- Stores every single API call made to an external LLM for auditing and cost tracking.
CREATE TABLE IF NOT EXISTS token_ledger (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    flow_id TEXT NOT NULL,
    model_used TEXT NOT NULL,
    agent_role TEXT NOT NULL,
    prompt_hash TEXT NOT NULL, -- Hashed prompt to prevent storing sensitive code, but allowing comparison
    input_tokens INTEGER NOT NULL,
    output_tokens INTEGER NOT NULL,
    total_cost_usd REAL NOT NULL,
    latency_ms INTEGER NOT NULL,
    status TEXT NOT NULL, -- 'SUCCESS', 'FAILED', 'TIMEOUT'
    error_message TEXT -- Detailed error log if the call failed
);

-- Table 2: forge_flows
-- Stores the visual flow definitions created by the user (the JSON graph).
CREATE TABLE IF NOT EXISTS forge_flows (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    flow_data TEXT NOT NULL, -- The serialized JSON structure of the nodes and edges
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Table 3: user_secrets
-- Stores encrypted credentials required for agent access (API keys).
CREATE TABLE IF NOT EXISTS user_secrets (
    key_name TEXT PRIMARY KEY, -- e.g., 'ANTHROPIC_API_KEY', 'GITHUB_TOKEN'
    encrypted_value BLOB NOT NULL, -- The encrypted API key/secret
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Index for fast retrieval of logs by flow ID
CREATE INDEX IF NOT EXISTS idx_ledger_flow_id ON token_ledger(flow_id);
`

// TokenLedgerPath is the filename for the SQLite database.
// This file will be managed by the Go BFF and should be excluded from Git (via .gitignore).
const TokenLedgerPath = "forge_ledger.db"
