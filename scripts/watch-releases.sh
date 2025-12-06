#!/bin/bash
# Watch Forge Terminal releases and sync handshake
# Run in background: ./scripts/watch-releases.sh &

set -e

REPO_OWNER="mikejsmith1985"
REPO_NAME="forge-terminal"
CHECK_INTERVAL=300  # 5 minutes
STATE_FILE=".forge/last-terminal-release"

echo "🔍 Starting Forge Terminal Release Watcher..."
echo "   Monitoring: $REPO_OWNER/$REPO_NAME"
echo "   Check interval: $CHECK_INTERVAL seconds"
echo "   State file: $STATE_FILE"

# Create state directory
mkdir -p "$(dirname "$STATE_FILE")"

# Get last checked version
if [ -f "$STATE_FILE" ]; then
    LAST_VERSION=$(cat "$STATE_FILE")
else
    LAST_VERSION=""
fi

echo "📦 Last known Terminal version: ${LAST_VERSION:-none}"
echo ""

while true; do
    # Fetch latest release from GitHub
    LATEST=$(curl -s "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest" | \
        grep -o '"tag_name": *"[^"]*"' | \
        sed 's/"tag_name": *"\([^"]*\)"/\1/' || echo "")
    
    if [ -z "$LATEST" ]; then
        echo "⚠️  [$(date '+%Y-%m-%d %H:%M:%S')] Failed to fetch latest release"
        sleep "$CHECK_INTERVAL"
        continue
    fi
    
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Terminal is at: $LATEST"
    
    # Check if new release
    if [ "$LATEST" != "$LAST_VERSION" ] && [ -n "$LAST_VERSION" ]; then
        echo ""
        echo "🎉 NEW TERMINAL RELEASE DETECTED: $LATEST"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo ""
        echo "📥 Downloading handshake from Terminal release..."
        
        # Download handshake from Terminal's release
        HANDSHAKE_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$LATEST/FORGE_HANDSHAKE.md"
        
        if curl -L -f "$HANDSHAKE_URL" -o "TERMINAL_HANDSHAKE.md" 2>/dev/null; then
            echo "✅ Terminal handshake downloaded"
            echo "📄 Location: ./TERMINAL_HANDSHAKE.md"
            
            # Validate
            if [ -f "scripts/validate-handshake.sh" ]; then
                echo ""
                echo "🔍 Validating Terminal handshake..."
                HANDSHAKE_FILE="TERMINAL_HANDSHAKE.md" ./scripts/validate-handshake.sh
            fi
            
            # Notification
            if command -v notify-send &> /dev/null; then
                notify-send "Forge Terminal Released" "Version $LATEST - Handshake available"
            fi
            
            echo ""
            echo "✅ Terminal handshake synchronized to $LATEST"
            echo "📋 Next steps:"
            echo "   1. Review TERMINAL_HANDSHAKE.md for new features"
            echo "   2. Update Orchestrator to match new features"
            echo "   3. Run compatibility tests"
            echo ""
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        else
            echo "⚠️  Terminal handshake not found in release"
            echo "   This may be an older release"
        fi
        echo ""
    fi
    
    # Save current version
    echo "$LATEST" > "$STATE_FILE"
    
    # Wait for next check
    sleep "$CHECK_INTERVAL"
done
