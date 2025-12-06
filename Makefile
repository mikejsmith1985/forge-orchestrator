.PHONY: help build test handshake validate-handshake sync-terminal watch-terminal

# Default target
help:
	@echo "Forge Orchestrator - Available Commands"
	@echo ""
	@echo "  make build              - Build the application"
	@echo "  make test               - Run all tests"
	@echo "  make handshake          - Generate handshake document"
	@echo "  make validate-handshake - Validate handshake document"
	@echo "  make sync-terminal      - Sync Terminal handshake"
	@echo "  make watch-terminal     - Watch Terminal releases (background)"
	@echo ""

# Build application
build:
	@echo "Building Forge Orchestrator..."
	cd frontend && npm install && npm run build
	go build -o forge-orchestrator .
	@echo "✅ Build complete: ./forge-orchestrator"

# Run tests
test:
	@echo "Running backend tests..."
	go test ./internal/... -v
	@echo "Running frontend tests..."
	cd frontend && npm test || true

# Generate handshake documentation
handshake:
	@./scripts/generate-handshake.sh

# Validate handshake documentation
validate-handshake:
	@./scripts/validate-handshake.sh

# Sync Terminal handshake
sync-terminal:
	@./sync-terminal-handshake.sh

# Watch Terminal releases in background
watch-terminal:
	@echo "Starting Terminal release watcher in background..."
	@./scripts/watch-releases.sh &
	@echo "✅ Watcher started (PID: $$!)"
	@echo "   Use 'ps aux | grep watch-releases' to check status"
	@echo "   Use 'kill <PID>' to stop"
