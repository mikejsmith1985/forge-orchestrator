#!/bin/bash
# Validate Forge Orchestrator Handshake Document

set -e

HANDSHAKE_FILE="FORGE_HANDSHAKE.md"

echo "🔍 Validating handshake document..."

if [ ! -f "$HANDSHAKE_FILE" ]; then
    echo "❌ ERROR: $HANDSHAKE_FILE not found"
    exit 1
fi

# Check file size
FILE_SIZE=$(wc -c < "$HANDSHAKE_FILE")
if [ "$FILE_SIZE" -lt 1000 ]; then
    echo "❌ ERROR: Handshake file too small ($FILE_SIZE bytes)"
    exit 1
fi

# Check required sections
REQUIRED_SECTIONS=(
    "Core Architecture"
    "API Endpoints"
    "UI Components"
    "Feature Requirements"
    "Configuration"
    "Testing"
    "Release Process"
)

MISSING_SECTIONS=()
for section in "${REQUIRED_SECTIONS[@]}"; do
    if ! grep -q "$section" "$HANDSHAKE_FILE"; then
        MISSING_SECTIONS+=("$section")
    fi
done

if [ ${#MISSING_SECTIONS[@]} -gt 0 ]; then
    echo "❌ ERROR: Missing required sections:"
    for section in "${MISSING_SECTIONS[@]}"; do
        echo "   - $section"
    done
    exit 1
fi

# Check version format
if ! grep -E "^\*\*Version\*\*:.*[0-9]+\.[0-9]+" "$HANDSHAKE_FILE" >/dev/null; then
    echo "⚠️  WARNING: Version format may be incorrect"
fi

# Check timestamp format
if ! grep -E "^\*\*Last Updated\*\*:.*[0-9]{4}-[0-9]{2}-[0-9]{2}" "$HANDSHAKE_FILE" >/dev/null; then
    echo "⚠️  WARNING: Timestamp format may be incorrect"
fi

# Count checkboxes
CHECKBOX_COUNT=$(grep -c "^\- \[.\]" "$HANDSHAKE_FILE" || echo 0)
echo "✅ Found $CHECKBOX_COUNT feature checkboxes"

# Summary
echo ""
echo "✅ Validation passed!"
echo "📄 File: $HANDSHAKE_FILE"
echo "📊 Size: $FILE_SIZE bytes"
echo "📋 Sections: ${#REQUIRED_SECTIONS[@]}"
echo "✓ Features: $CHECKBOX_COUNT"
