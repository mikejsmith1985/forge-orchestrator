# 🔄 Handshake Automation - Quick Reference

## Overview

Forge Orchestrator automatically synchronizes feature specifications with Forge Terminal through "handshake" documents.

## Quick Commands

```bash
# Generate Orchestrator's handshake
make handshake
# or
./scripts/generate-handshake.sh

# Sync Terminal's handshake
make sync-terminal
# or
./sync-terminal-handshake.sh

# Validate handshake
make validate-handshake

# Watch for Terminal releases (background)
make watch-terminal
# or
./scripts/watch-releases.sh &
```

## What Gets Generated

### Orchestrator Handshake (`FORGE_HANDSHAKE.md`)
- **When**: On every release (GitHub Actions)
- **Contains**: Orchestrator's features, APIs, components
- **Purpose**: Document what Orchestrator provides
- **Included in**: GitHub releases

### Terminal Handshake (`TERMINAL_HANDSHAKE.md`)  
- **When**: Downloaded from Terminal releases
- **Contains**: Terminal's features (reference implementation)
- **Purpose**: Know what features to match
- **Updated**: Via sync scripts

## Workflows

### For Development

1. **Start Background Watcher** (one-time setup)
   ```bash
   make watch-terminal
   ```
   - Checks every 5 minutes for new Terminal releases
   - Auto-downloads handshake
   - Shows desktop notification

2. **Or Manual Sync** (on-demand)
   ```bash
   make sync-terminal
   ```
   - Downloads latest Terminal handshake
   - Shows version summary

### For Releases

1. **Tag and Push**
   ```bash
   git tag v1.2.1
   git push origin v1.2.1
   ```

2. **GitHub Actions Automatically**:
   - Runs tests
   - Builds binaries
   - Generates `FORGE_HANDSHAKE.md`
   - Validates handshake
   - Creates release with handshake

### Maintaining Feature Parity

1. **Check for Updates**
   ```bash
   make sync-terminal
   ```

2. **Review Changes**
   ```bash
   # Compare handshakes
   diff TERMINAL_HANDSHAKE.md FORGE_HANDSHAKE.md
   
   # Look for new features
   grep "\[ \]" TERMINAL_HANDSHAKE.md
   ```

3. **Implement Features**
   - Add missing features from Terminal
   - Run tests
   - Update documentation

4. **Generate Updated Handshake**
   ```bash
   make handshake
   make validate-handshake
   ```

5. **Commit and Release**
   ```bash
   git add FORGE_HANDSHAKE.md
   git commit -m "Update handshake"
   git tag v1.2.2
   git push --tags
   ```

## File Structure

```
forge-orchestrator/
├── FORGE_HANDSHAKE.md              # Our feature spec (generated)
├── TERMINAL_HANDSHAKE.md           # Terminal's spec (downloaded)
├── Makefile                         # Convenience commands
├── sync-terminal-handshake.sh      # Quick sync script
├── scripts/
│   ├── generate-handshake.sh       # Generate our handshake
│   ├── validate-handshake.sh       # Validate completeness
│   └── watch-releases.sh           # Background watcher
└── docs/
    └── RELEASE_AUTOMATION.md       # Full documentation
```

## Troubleshooting

### "Terminal handshake not found"
- Older Terminal releases don't have handshakes
- Wait for next Terminal release
- Or manually review Terminal code

### "Failed to fetch latest release"
- Check internet connection
- Check GitHub API rate limit: `curl -s https://api.github.com/rate_limit | grep remaining`
- Use GitHub token if needed

### Watcher not working
- Run in foreground to see errors: `./scripts/watch-releases.sh`
- Check state file: `cat .forge/last-terminal-release`
- Verify GitHub API access: `curl -s https://api.github.com/repos/mikejsmith1985/forge-terminal/releases/latest`

## More Information

See [docs/RELEASE_AUTOMATION.md](docs/RELEASE_AUTOMATION.md) for complete documentation including:
- Detailed workflows
- GitHub Actions setup
- System service installation
- Advanced troubleshooting
- Best practices

## Related Links

- **Terminal Repo**: https://github.com/mikejsmith1985/forge-terminal
- **Orchestrator Repo**: https://github.com/mikejsmith1985/forge-orchestrator
- **GitHub Actions Docs**: https://docs.github.com/en/actions
