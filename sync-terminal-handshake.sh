#!/bin/bash
# Quick sync script - Download Terminal handshake from latest GitHub release

REPO_OWNER="mikejsmith1985"
REPO_NAME="forge-terminal"

echo "🔄 Syncing Terminal handshake from latest GitHub release..."
echo ""

# Get latest version
echo "📦 Fetching latest Terminal release..."
VERSION=$(curl -s "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest" | \
    grep '"tag_name"' | \
    sed -E 's/.*"tag_name": "([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo "❌ Failed to fetch latest Terminal release"
    echo ""
    echo "Possible reasons:"
    echo "  - No internet connection"
    echo "  - GitHub API rate limit"
    echo "  - Repository not found"
    echo ""
    exit 1
fi

echo "📦 Latest Terminal version: $VERSION"
echo ""

# Download handshake
HANDSHAKE_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$VERSION/FORGE_HANDSHAKE.md"
echo "📥 Downloading from: $HANDSHAKE_URL"

if curl -L -f "$HANDSHAKE_URL" -o "TERMINAL_HANDSHAKE.md" 2>/dev/null; then
    echo "✅ Terminal handshake downloaded"
    echo "📄 Location: ./TERMINAL_HANDSHAKE.md"
    
    # Show summary
    echo ""
    echo "📊 Summary:"
    grep "^**Version**:" TERMINAL_HANDSHAKE.md || echo "  Version: $VERSION"
    grep "^**Last Updated**:" TERMINAL_HANDSHAKE.md || true
    echo ""
    
    # Save to state file
    mkdir -p .forge
    echo "$VERSION" > .forge/last-terminal-release
    
    echo "💡 What's next?"
    echo "   1. Review TERMINAL_HANDSHAKE.md for changes"
    echo "   2. Update Orchestrator features to match"
    echo "   3. Run tests to ensure compatibility"
    echo ""
    
    exit 0
else
    echo "⚠️  Terminal handshake not found in release assets"
    echo ""
    echo "This could mean:"
    echo "  - Terminal release didn't include handshake (older version)"
    echo "  - Asset name changed"
    echo "  - Network error"
    echo ""
    echo "Manual alternative:"
    echo "  Visit: https://github.com/$REPO_OWNER/$REPO_NAME/releases/latest"
    echo ""
    exit 1
fi
