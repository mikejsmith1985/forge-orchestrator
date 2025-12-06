# Terminal Connection Fix - Issue #3

## Problem Summary
Terminal connections were failing on Windows systems because the code was hardcoded to use `bash`, which doesn't exist on Windows. Users on Windows need support for PowerShell, CMD, and WSL terminals.

## Solution Implemented

### 1. Multi-Shell Support (Backend)
- **Platform-specific PTY implementations**:
  - `pty_unix.go` - Unix/Linux using creack/pty
  - `pty_windows.go` - Windows using ConPTY
- **Configurable shell selection**:
  - Bash (Unix/Linux)
  - CMD (Windows)
  - PowerShell (Windows)
  - WSL (Windows with distribution selection)
- **Enhanced logging** for debugging connection issues
- **Better error messages** sent to terminal when connection fails

### 2. Settings UI (Frontend)
- **New Terminal Settings page** (`/settings` → Terminal tab)
- **Shell selection UI** with platform-appropriate options
- **WSL configuration** including distro selection
- **Troubleshooting help** section
- **Configuration persistence** via API

### 3. API Endpoints
- `GET /api/config` - Retrieve current configuration
- `POST /api/config` - Save configuration
- Both endpoints work with shell settings

### 4. Bug Fixes
- Fixed panic due to channel being closed twice
- Added `sync.Once` for safe channel closure
- Updated CORS to allow Playwright test runner (port 8081)

## Files Changed

### Backend (Go)
- `internal/server/pty_manager.go` - Multi-shell support with config
- `internal/server/pty_unix.go` - NEW: Unix PTY implementation
- `internal/server/pty_windows.go` - NEW: Windows ConPTY implementation  
- `internal/server/websocket.go` - Enhanced error messages
- `internal/server/api_handlers.go` - Config API endpoints
- `internal/server/routes.go` - Config routes
- `internal/server/middleware.go` - CORS update for tests
- `go.mod` - Added ConPTY dependency

### Frontend (React/TypeScript)
- `frontend/src/components/Settings/Settings.tsx` - NEW: Settings tabs
- `frontend/src/components/Settings/TerminalSettings.tsx` - NEW: Terminal config UI
- `frontend/src/components/Settings/index.ts` - NEW: Exports
- `frontend/src/App.tsx` - Updated to use new Settings component
- `frontend/tests/e2e/terminal-settings.spec.ts` - NEW: Comprehensive tests

## Testing Results

### Playwright E2E Tests
- ✅ All original terminal tests pass (7/7)
- ✅ Terminal settings tests created (11 tests)
- ✅ PTY connections establish successfully
- ✅ No panics or crashes
- ✅ Proper logging and error handling

### Manual Testing Required
Since this fix targets Windows, **manual testing on Windows** is required to verify:
1. CMD terminal connects
2. PowerShell terminal connects
3. WSL terminal connects (if WSL installed)
4. Settings UI works correctly
5. Configuration persists across restarts

## How to Use

### For Users
1. Go to Settings → Terminal tab
2. Select your preferred shell type
3. For WSL: optionally specify distribution name
4. Click "Save Configuration"
5. Go to Terminal tab - it will connect with new shell

### For Windows Users with Connection Issues
1. Try changing shell type in Settings
2. For WSL: Run `wsl --list` in CMD to see available distributions
3. Check browser console (F12) for detailed error messages
4. Terminal will show helpful troubleshooting tips on failure

## Architecture Notes

- **Platform detection** via runtime.GOOS at build time
- **ConPTY** (github.com/UserExistsError/conpty) for Windows
- **creack/pty** for Unix/Linux
- **Config-driven** shell selection
- **Graceful degradation** with helpful error messages

## Next Steps
1. Test on actual Windows machine
2. Consider adding shell auto-detection
3. Add WSL home path configuration
4. Add shell preference to user profile

