#!/bin/bash
# Generate Forge Orchestrator Handshake Document
# This script extracts features, APIs, and components to create a comprehensive spec

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT_DIR="$SCRIPT_DIR/.."
OUTPUT_FILE="$ROOT_DIR/FORGE_HANDSHAKE.md"

echo "🔥 Generating Forge Orchestrator Handshake Document..."

# Extract version from Go code
BACKEND_VERSION=$(grep 'var Version' "$ROOT_DIR/internal/updater/updater.go" | sed 's/.*"\(.*\)".*/\1/')
echo "📦 Backend Version: $BACKEND_VERSION"

# Extract version from package.json
FRONTEND_VERSION=$(grep '"version"' "$ROOT_DIR/frontend/package.json" | head -1 | sed 's/.*"\(.*\)".*/\1/' || echo "0.0.0")
echo "📦 Frontend Version: $FRONTEND_VERSION"

# Count components
COMPONENT_COUNT=$(find "$ROOT_DIR/frontend/src/components" -name "*.tsx" -o -name "*.ts" 2>/dev/null | wc -l)
echo "🎨 React Components: $COMPONENT_COUNT"

# Extract API endpoints from routes.go
echo "🔌 Extracting API endpoints..."
API_ENDPOINTS_FILE=$(mktemp)
grep -E 'Handle(Func)?\(' "$ROOT_DIR/internal/server/routes.go" 2>/dev/null | \
    grep -o '"/[^"]*"' | \
    tr -d '"' | \
    sort -u > "$API_ENDPOINTS_FILE"

# Count endpoints
ENDPOINT_COUNT=$(wc -l < "$API_ENDPOINTS_FILE")
echo "🔌 Found $ENDPOINT_COUNT API endpoints"

# Get current timestamp
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%S.000Z")

# Generate document
cat > "$OUTPUT_FILE" << 'HANDSHAKE_DOC'
# Forge Orchestrator ↔ Forge Terminal Handshake Specification

**Version**: VERSION_PLACEHOLDER  
**Last Updated**: TIMESTAMP_PLACEHOLDER  
**Purpose**: Ensure Forge Orchestrator maintains 1:1 feature parity with Forge Terminal

---

## 🎯 Core Architecture

### Application Type
- **Platform**: Desktop application (native binary + embedded web UI)
- **Backend**: Go HTTP server with WebSocket support
- **Frontend**: React SPA with XTerm.js terminal emulator
- **Distribution**: Single executable binary (Windows, macOS, Linux)

### Technical Stack
```
Backend:
  - Language: Go 1.23
  - Web Server: net/http (stdlib)
  - WebSocket: gorilla/websocket
  - Terminal: pty (Unix) / conpty (Windows)
  - Database: SQLite (modernc.org/sqlite)

Frontend:
  - Framework: React 18
  - Build Tool: Vite
  - Terminal: XTerm.js + addons (fit, web-links, search)
  - UI Library: lucide-react icons
  - Flow Editor: React Flow
```

---

## 🔌 API Endpoints (AUTO-DETECTED)

API_ENDPOINTS_PLACEHOLDER

---

## 🎨 UI Components (COMPONENT_COUNT_PLACEHOLDER React Components)

### Core Views
1. **Architect View** - AI-powered workflow designer
2. **Terminal View** - Integrated terminal with PTY
3. **Ledger View** - Blockchain-style transaction log
4. **Command Deck** - Quick command execution
5. **Flow Editor** - Visual workflow builder
6. **Settings** - Application configuration

### Terminal Components
- **Terminal.tsx** - Enhanced XTerm.js wrapper
  - Auto-respond to CLI prompts
  - Auto-reconnection with exponential backoff
  - Connection status overlay
  - Search functionality
  - Scroll-to-bottom button
- **TerminalSettings.tsx** - Shell configuration (WSL, PowerShell, CMD, Bash)

### Layout Components
- **Sidebar.tsx** - Navigation sidebar
- **MainContent.tsx** - Content area wrapper

### Feature Components
- **Architect/**
  - ArchitectView.tsx - Main architect interface
  - Chat integration
  - Model selection
  - Task management
- **Flows/**
  - FlowList.tsx - Flow management
  - FlowEditor.tsx - Visual flow editor
  - Node components
- **Ledger/**
  - LedgerView.tsx - Transaction viewer
  - Filtering and search
- **Commands/**
  - CommandDeck.tsx - Command card interface
  - Drag and drop support
- **Settings/**
  - Settings.tsx - Settings tabs
  - TerminalSettings.tsx - Terminal configuration
- **Update/**
  - UpdateModal.tsx - Update notifications
  - UpdateToast.tsx - Update alerts
- **Welcome/**
  - WelcomeModal.tsx - First-run experience
- **Feedback/**
  - FeedbackModal.tsx - User feedback with screenshots

---

## 🔥 Feature Requirements (Must Match Terminal)

### Terminal Features ✅
- [x] Multi-shell support (Bash, PowerShell, CMD, WSL)
- [x] WebSocket-based PTY communication
- [x] Auto-respond to confirmation prompts
- [x] Auto-reconnection with exponential backoff
- [x] Connection status overlay
- [x] Search functionality (SearchAddon)
- [x] Scroll-to-bottom button
- [x] Terminal resize handling
- [x] ANSI color support
- [x] Clickable URLs
- [x] Settings persistence

### Core Features ✅
- [x] Desktop notifications
- [x] Auto-update system
- [x] Feedback system with screenshots
- [x] Welcome/onboarding flow
- [x] Settings management
- [x] Theme support
- [x] Keyboard shortcuts
- [x] Error handling
- [x] Logging system

### Orchestrator-Specific Features
- [x] Workflow designer (Architect)
- [x] Flow visual editor
- [x] Ledger transaction log
- [x] Command deck with cards
- [x] AI integration
- [x] State persistence
- [x] Flow execution engine

---

## 📋 Configuration

### Supported Shells
```json
{
  "shell": {
    "type": "bash" | "cmd" | "powershell" | "wsl",
    "wsl_distro": "Ubuntu-24.04",
    "wsl_user": "username"
  },
  "server": {
    "port": 8080,
    "open_browser": true
  }
}
```

### WSL Configuration
- Automatic distro detection
- Custom home path support
- Multi-distro support
- Native Linux environment

---

## 🧪 Testing

### E2E Tests (Playwright)
- ✅ Terminal integration (7 tests)
- ✅ Enhanced terminal features (11 tests)
- ✅ Terminal settings (11 tests)
- **Total: 29 passing tests**

### Unit Tests
- Go backend tests
- React component tests
- Integration tests

---

## 📦 Release Process

### Automated CI/CD
1. Push version tag (v*.*.*)
2. GitHub Actions builds:
   - Runs unit tests
   - Builds frontend (Vite)
   - Compiles Go binaries (Linux, Windows, macOS)
   - Generates handshake document
   - Validates handshake
3. Creates GitHub release with:
   - Binaries for all platforms
   - FORGE_HANDSHAKE.md
   - Auto-generated release notes

### Manual Release
```bash
# Tag version
git tag v1.2.0
git push origin v1.2.0

# Or use make
make release VERSION=1.2.0
```

---

## 🔄 Sync from Terminal

### Automatic Sync (Background Watcher)
```bash
# Start background watcher
./scripts/watch-releases.sh &

# Checks every 5 minutes for new Terminal releases
# Auto-downloads handshake
# Desktop notification on updates
```

### Manual Sync
```bash
# Quick one-time sync
./sync-handshake.sh

# Downloads latest FORGE_HANDSHAKE.md from Terminal releases
```

### GitHub Actions Sync
The repository can be configured to automatically check for Terminal updates:
- Scheduled workflow (every 4 hours)
- Repository dispatch events
- Automatic compatibility testing

---

## 🤝 Handshake Contract

### Orchestrator Must Provide
1. All Terminal features (see Feature Requirements above)
2. Additional workflow/architect features
3. Backward compatibility
4. Same UI/UX patterns
5. Same configuration format

### Terminal Provides
1. Reference implementation
2. Feature specifications
3. API contracts
4. Component patterns
5. Release handshakes

---

## 📚 Documentation

- **README.md** - Getting started
- **TERMINAL_FIX_SUMMARY.md** - Terminal implementation details
- **ISSUE_04_SOLUTION.md** - Enhanced terminal features
- **docs/RELEASE_AUTOMATION.md** - Release automation guide
- **contracts/** - API contracts and schemas

---

## 🔗 Links

- **Repository**: https://github.com/mikejsmith1985/forge-orchestrator
- **Terminal Repo**: https://github.com/mikejsmith1985/forge-terminal
- **Issues**: https://github.com/mikejsmith1985/forge-orchestrator/issues

---

## 📊 Version History

| Version | Date | Changes |
|---------|------|---------|
| VERSION_PLACEHOLDER | TIMESTAMP_PLACEHOLDER | Initial handshake document |

---

**Generated by**: `scripts/generate-handshake.sh`  
**Validation**: `scripts/validate-handshake.sh`
HANDSHAKE_DOC

# Replace placeholders
sed -i "s/VERSION_PLACEHOLDER/$BACKEND_VERSION/" "$OUTPUT_FILE"
sed -i "s/TIMESTAMP_PLACEHOLDER/$TIMESTAMP/" "$OUTPUT_FILE"
sed -i "s/COMPONENT_COUNT_PLACEHOLDER/$COMPONENT_COUNT/" "$OUTPUT_FILE"

# Insert API endpoints
API_SECTION=""
while IFS= read -r endpoint; do
    API_SECTION="${API_SECTION}- \`${endpoint}\`"$'\n'
done < "$API_ENDPOINTS_FILE"

# Create temp file with API endpoints
cat > /tmp/api_section.txt << EOF
$API_SECTION
EOF

# Use awk to replace the placeholder
awk '/API_ENDPOINTS_PLACEHOLDER/ {
    system("cat /tmp/api_section.txt")
    next
}
{print}' "$OUTPUT_FILE" > "$OUTPUT_FILE.tmp"
mv "$OUTPUT_FILE.tmp" "$OUTPUT_FILE"

# Cleanup
rm -f "$API_ENDPOINTS_FILE" /tmp/api_section.txt

echo "✅ Handshake document generated: $OUTPUT_FILE"
echo "📄 Size: $(wc -c < "$OUTPUT_FILE") bytes"
echo "📊 Components: $COMPONENT_COUNT"
echo "🔌 Endpoints: $ENDPOINT_COUNT"
