// Package server provides the HTTP/WebSocket server and PTY management.
// This file implements the PTY manager that handles pseudo-terminal sessions.
package server

import (
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

// PTYSession represents an active PTY session connected to a WebSocket client.
// It manages the bidirectional communication between the terminal and the browser.
type PTYSession struct {
	// The pseudo-terminal file descriptor
	ptmx *os.File
	// The underlying shell command
	cmd *exec.Cmd
	// WebSocket connection to the browser
	conn *websocket.Conn
	// Mutex for thread-safe writes to WebSocket
	writeMu sync.Mutex
	// Channel to signal session closure
	done chan struct{}
	// Flag indicating if prompt watcher is enabled
	promptWatcherEnabled bool
	// Mutex for prompt watcher state
	promptMu sync.Mutex
}

// PTYManager manages all active PTY sessions.
// It provides methods to create, access, and destroy terminal sessions.
type PTYManager struct {
	// Map of session ID to PTY session
	sessions map[string]*PTYSession
	// Mutex for thread-safe access to sessions map
	mu sync.RWMutex
}

// NewPTYManager creates a new PTY manager instance.
func NewPTYManager() *PTYManager {
	return &PTYManager{
		sessions: make(map[string]*PTYSession),
	}
}

// CreateSession creates a new PTY session for a WebSocket client.
// It starts a bash shell and connects it to the WebSocket.
func (pm *PTYManager) CreateSession(sessionID string, conn *websocket.Conn) (*PTYSession, error) {
	// Start a bash shell
	cmd := exec.Command("bash")
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	// Create the pseudo-terminal
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	session := &PTYSession{
		ptmx: ptmx,
		cmd:  cmd,
		conn: conn,
		done: make(chan struct{}),
	}

	pm.mu.Lock()
	pm.sessions[sessionID] = session
	pm.mu.Unlock()

	// Start goroutine to read from PTY and send to WebSocket
	go session.readPTYLoop()

	return session, nil
}

// GetSession retrieves an active PTY session by ID.
func (pm *PTYManager) GetSession(sessionID string) *PTYSession {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.sessions[sessionID]
}

// CloseSession closes and removes a PTY session.
func (pm *PTYManager) CloseSession(sessionID string) {
	pm.mu.Lock()
	session, exists := pm.sessions[sessionID]
	if exists {
		delete(pm.sessions, sessionID)
	}
	pm.mu.Unlock()

	if session != nil {
		session.Close()
	}
}

// readPTYLoop reads output from the PTY and sends it to the WebSocket.
func (s *PTYSession) readPTYLoop() {
	buf := make([]byte, 4096)
	for {
		select {
		case <-s.done:
			return
		default:
			n, err := s.ptmx.Read(buf)
			if err != nil {
				if err != io.EOF {
					// Log error but continue - connection may have closed
				}
				return
			}

			if n > 0 {
				data := buf[:n]

				// Check for confirmation prompts if prompt watcher is enabled
				s.promptMu.Lock()
				watcherEnabled := s.promptWatcherEnabled
				s.promptMu.Unlock()

				if watcherEnabled {
					s.checkAndRespondToPrompts(data)
				}

				// Send to WebSocket
				s.writeMu.Lock()
				err := s.conn.WriteMessage(websocket.TextMessage, data)
				s.writeMu.Unlock()

				if err != nil {
					return
				}
			}
		}
	}
}

// checkAndRespondToPrompts checks PTY output for confirmation prompts
// and automatically responds with 'y' if the prompt watcher is enabled.
func (s *PTYSession) checkAndRespondToPrompts(data []byte) {
	// Common confirmation prompt patterns
	output := string(data)
	patterns := []string{
		"[y/n]",
		"[Y/n]",
		"[y/N]",
		"(y/n)",
		"(Y/n)",
		"(y/N)",
		"Continue? [y/n]",
		"Proceed? [y/n]",
		"Are you sure",
	}

	for _, pattern := range patterns {
		if containsIgnoreCase(output, pattern) {
			// Inject 'y' response
			s.Write([]byte("y\n"))
			return
		}
	}
}

// containsIgnoreCase checks if s contains substr (case-insensitive).
func containsIgnoreCase(s, substr string) bool {
	sLower := make([]byte, len(s))
	substrLower := make([]byte, len(substr))
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			sLower[i] = s[i] + 32
		} else {
			sLower[i] = s[i]
		}
	}
	for i := 0; i < len(substr); i++ {
		if substr[i] >= 'A' && substr[i] <= 'Z' {
			substrLower[i] = substr[i] + 32
		} else {
			substrLower[i] = substr[i]
		}
	}
	return contains(string(sLower), string(substrLower))
}

// contains is a simple substring check.
func contains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Write sends data to the PTY (simulates typing).
func (s *PTYSession) Write(data []byte) (int, error) {
	return s.ptmx.Write(data)
}

// WriteCommand writes a command string to the PTY.
// This simulates a user typing a command and pressing Enter.
func (s *PTYSession) WriteCommand(command string) error {
	_, err := s.ptmx.Write([]byte(command + "\n"))
	return err
}

// SetPromptWatcher enables or disables the prompt watcher.
func (s *PTYSession) SetPromptWatcher(enabled bool) {
	s.promptMu.Lock()
	s.promptWatcherEnabled = enabled
	s.promptMu.Unlock()
}

// Resize changes the PTY window size.
func (s *PTYSession) Resize(rows, cols uint16) error {
	return pty.Setsize(s.ptmx, &pty.Winsize{
		Rows: rows,
		Cols: cols,
	})
}

// Close terminates the PTY session and cleans up resources.
func (s *PTYSession) Close() {
	close(s.done)

	if s.ptmx != nil {
		s.ptmx.Close()
	}

	if s.cmd != nil && s.cmd.Process != nil {
		s.cmd.Process.Kill()
		s.cmd.Wait()
	}
}
