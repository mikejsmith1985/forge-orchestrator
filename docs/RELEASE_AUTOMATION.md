# Release Automation Guide

## 🎯 Overview

This document describes the automated handshake synchronization system between **Forge Terminal** (reference implementation) and **Forge Orchestrator** (enhanced version).

### Purpose
Ensure Forge Orchestrator maintains 1:1 feature parity with Forge Terminal by automatically tracking and documenting feature changes through handshake documents.

---

## 🔄 How It Works

### The Handshake Flow

```
┌─────────────────┐
│ Forge Terminal  │
│   (Reference)   │
└────────┬────────┘
         │
         │ Push tag v1.x.x
         ▼
┌─────────────────┐
│ GitHub Actions  │
│  - Build app    │
│  - Generate     │
│    HANDSHAKE.md │
│  - Validate     │
│  - Release      │
└────────┬────────┘
         │
         │ Publish release
         ▼
┌─────────────────┐     ┌──────────────────┐
│ GitHub Release  │────▶│ Orchestrator     │
│ + HANDSHAKE.md  │     │ Auto-sync or     │
└─────────────────┘     │ Manual download  │
                        └────────┬─────────┘
                                 │
                                 ▼
                        ┌─────────────────┐
                        │ Update features │
                        │ to match spec   │
                        └─────────────────┘
```

---

## 📦 For Forge Terminal (Reference Implementation)

### Automated Release Process

1. **Tag and Push**
   ```bash
   git tag v1.9.3
   git push origin v1.9.3
   ```

2. **GitHub Actions Automatically**:
   - Builds all binaries
   - Generates `FORGE_HANDSHAKE.md` from source code
   - Validates completeness
   - Creates release with handshake included

3. **Handshake Contains**:
   - Version information
   - API endpoints (auto-detected)
   - UI components (auto-counted)
   - Feature checklist
   - Configuration specs
   - Test requirements

### Manual Handshake Generation

```bash
# Generate handshake locally
make handshake-doc

# Or directly
./scripts/generate-handshake.sh

# Validate
make validate-handshake
```

### Scripts Included

- `scripts/generate-handshake.sh` - Extracts features from codebase
- `scripts/validate-handshake.sh` - Ensures completeness
- `scripts/watch-releases.sh` - Background watcher (not typically used by Terminal)

---

## 🎛️ For Forge Orchestrator (Enhanced Version)

### Automatic Synchronization

#### Option 1: Background Watcher (Recommended for Dev)

Start a background process that polls Terminal releases every 5 minutes:

```bash
# Start in background
./scripts/watch-releases.sh &

# Check if running
ps aux | grep watch-releases
```

**What it does:**
- Checks GitHub every 5 minutes
- Downloads Terminal handshake when new release detected
- Saves to `TERMINAL_HANDSHAKE.md`
- Shows desktop notification (if available)
- Logs to console

**System Service (Optional)**:
```bash
# Copy service file
sudo cp scripts/forge-terminal-watcher@.service /etc/systemd/system/

# Enable and start
sudo systemctl enable forge-terminal-watcher@$USER.service
sudo systemctl start forge-terminal-watcher@$USER.service

# Check status
sudo systemctl status forge-terminal-watcher@$USER.service
```

#### Option 2: Manual Sync (Quick)

Download the latest Terminal handshake on demand:

```bash
./sync-terminal-handshake.sh
```

**What it does:**
- Fetches latest Terminal release from GitHub
- Downloads `FORGE_HANDSHAKE.md`
- Saves as `TERMINAL_HANDSHAKE.md`
- Shows summary

#### Option 3: GitHub Actions Scheduled Sync

Add to `.github/workflows/sync-terminal.yml`:

```yaml
name: Sync Terminal Handshake

on:
  schedule:
    - cron: '0 */4 * * *'  # Every 4 hours
  workflow_dispatch:  # Manual trigger

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Sync Terminal Handshake
        run: |
          ./sync-terminal-handshake.sh
          
      - name: Check for changes
        id: changes
        run: |
          if [ -n "$(git status --porcelain)" ]; then
            echo "has_changes=true" >> $GITHUB_OUTPUT
          fi
          
      - name: Create PR
        if: steps.changes.outputs.has_changes == 'true'
        uses: peter-evans/create-pull-request@v5
        with:
          title: 'Update Terminal handshake'
          body: 'New Forge Terminal release detected'
          branch: sync/terminal-handshake
```

### Orchestrator Release Process

When creating an Orchestrator release:

```bash
# Tag version
git tag v1.2.0
git push origin v1.2.0
```

GitHub Actions automatically:
1. Runs tests
2. Builds binaries
3. **Generates Orchestrator handshake** (`FORGE_HANDSHAKE.md`)
4. Validates handshake
5. Creates release with:
   - Binaries
   - `FORGE_HANDSHAKE.md` (Orchestrator's own spec)

### Maintaining Feature Parity

1. **Regular Checks**
   ```bash
   # Download latest Terminal handshake
   ./sync-terminal-handshake.sh
   
   # Compare with Orchestrator's features
   diff TERMINAL_HANDSHAKE.md FORGE_HANDSHAKE.md
   ```

2. **Update Process**
   - Review Terminal handshake for new features
   - Implement matching features in Orchestrator
   - Update Orchestrator handshake
   - Run tests to verify compatibility

3. **Compatibility Testing**
   ```bash
   # Run full test suite
   cd frontend && npm test
   go test ./...
   
   # Run E2E tests
   cd frontend && npx playwright test
   ```

---

## 📁 File Structure

```
forge-orchestrator/
├── .github/workflows/
│   └── release.yml              # Auto-generates handshake on release
├── scripts/
│   ├── generate-handshake.sh    # Generate Orchestrator handshake
│   ├── validate-handshake.sh    # Validate handshake completeness
│   └── watch-releases.sh        # Watch Terminal for new releases
├── sync-terminal-handshake.sh   # Quick sync from Terminal
├── FORGE_HANDSHAKE.md          # Orchestrator's feature spec (generated)
└── TERMINAL_HANDSHAKE.md       # Terminal's feature spec (downloaded)
```

---

## 🔧 Configuration

### Environment Variables

```bash
# Override repository (if forked)
export TERMINAL_REPO_OWNER="mikejsmith1985"
export TERMINAL_REPO_NAME="forge-terminal"

# Change check interval (seconds)
export CHECK_INTERVAL=300  # 5 minutes
```

### State Files

- `.forge/last-terminal-release` - Tracks last synced Terminal version
- `.forge/last-release-check` - Timestamp of last GitHub API check

---

## 🧪 Testing the System

### Test Handshake Generation

```bash
# Generate Orchestrator handshake
./scripts/generate-handshake.sh

# Validate it
./scripts/validate-handshake.sh

# Check output
cat FORGE_HANDSHAKE.md | head -50
```

### Test Terminal Sync

```bash
# Manual sync
./sync-terminal-handshake.sh

# Check result
ls -lh TERMINAL_HANDSHAKE.md
cat TERMINAL_HANDSHAKE.md | head -30
```

### Test Background Watcher

```bash
# Start watcher
./scripts/watch-releases.sh &

# Check it's running
ps aux | grep watch-releases

# View logs
tail -f ~/.forge/watcher.log  # If configured

# Stop watcher
killall watch-releases.sh
```

---

## 🚨 Troubleshooting

### Handshake Generation Fails

**Problem**: Script can't extract version or endpoints

**Solution**:
```bash
# Check Go code structure
grep -r "GetVersion" internal/updater/

# Check routes file exists
ls -la internal/server/routes.go

# Check frontend structure
ls -la frontend/src/components/
```

### Terminal Sync Fails

**Problem**: Can't download from GitHub

**Possible causes**:
- No internet connection
- GitHub API rate limit (60 req/hour unauthenticated)
- Terminal release doesn't include handshake

**Solutions**:
```bash
# Check internet
curl -I https://api.github.com

# Check rate limit
curl -s https://api.github.com/rate_limit | grep remaining

# Use authentication (higher limits)
export GITHUB_TOKEN="your_token"
curl -H "Authorization: Bearer $GITHUB_TOKEN" ...

# Manual download
open https://github.com/mikejsmith1985/forge-terminal/releases/latest
```

### Watcher Not Working

**Problem**: Background watcher doesn't detect new releases

**Debug**:
```bash
# Run in foreground to see output
./scripts/watch-releases.sh

# Check state file
cat .forge/last-terminal-release

# Test API access
curl -s "https://api.github.com/repos/mikejsmith1985/forge-terminal/releases/latest"
```

---

## 📊 Metrics & Monitoring

### What to Track

- Terminal version currently synced
- Last sync timestamp
- Number of API endpoints
- Number of UI components
- Feature checklist completion %

### Automation Ideas

```bash
# Daily report script
#!/bin/bash
echo "=== Forge Parity Report ==="
echo "Terminal Version: $(cat .forge/last-terminal-release)"
echo "Orchestrator Version: $(grep 'Version' FORGE_HANDSHAKE.md | head -1)"
echo "Last Sync: $(stat -c %y TERMINAL_HANDSHAKE.md 2>/dev/null || echo 'Never')"
echo ""
echo "Feature Parity:"
grep -c "\[x\]" FORGE_HANDSHAKE.md || echo "0"
echo "features implemented"
```

---

## 🎯 Best Practices

1. **Always Sync Before Major Work**
   ```bash
   ./sync-terminal-handshake.sh
   ```

2. **Review Changes**
   - Read Terminal handshake when updated
   - Note new features
   - Plan implementation

3. **Keep Documentation Updated**
   - Update `FORGE_HANDSHAKE.md` after adding features
   - Run validation after changes
   - Commit handshake with code changes

4. **Test Before Release**
   - Generate handshake locally
   - Validate completeness
   - Run full test suite

5. **Monitor Terminal Releases**
   - Watch GitHub releases
   - Review release notes
   - Test compatibility

---

## 🔗 References

- **Terminal Repository**: https://github.com/mikejsmith1985/forge-terminal
- **Orchestrator Repository**: https://github.com/mikejsmith1985/forge-orchestrator
- **GitHub API**: https://docs.github.com/en/rest/releases
- **GitHub Actions**: https://docs.github.com/en/actions

---

## 📝 Changelog

| Date | Change |
|------|--------|
| 2024-12-06 | Initial automation system implemented |
| | - Handshake generation |
| | - Terminal sync scripts |
| | - Background watcher |
| | - GitHub Actions integration |

---

**Maintained by**: Forge Team  
**Last Updated**: 2024-12-06
