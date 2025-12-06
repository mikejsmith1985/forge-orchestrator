package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mikejsmith1985/forge-orchestrator/internal/config"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true // Same-origin requests don't send Origin header
		}
		allowed := IsAllowedOrigin(origin)
		if !allowed {
			log.Printf("WebSocket: Blocked connection from origin: %s", origin)
		}
		return allowed
	},
}

// PTYMessage represents a message sent to/from the PTY WebSocket.
type PTYMessage struct {
	Type string `json:"type"` // "input", "resize", "prompt_watcher"
	Data string `json:"data,omitempty"`
	Rows uint16 `json:"rows,omitempty"`
	Cols uint16 `json:"cols,omitempty"`
}

// handleWebSocket upgrades the HTTP connection to a WebSocket connection
// and handles the client communication using the Hub pattern.
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	client := NewClient(s.hub, conn)
	s.hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// handlePTYWebSocket handles WebSocket connections for the integrated terminal.
// It creates a PTY session and streams data bidirectionally between the
// browser terminal (xterm.js) and the local shell.
func (s *Server) handlePTYWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("PTY WebSocket upgrade error:", err)
		return
	}

	// Generate a unique session ID
	sessionID := uuid.New().String()

	log.Printf("Creating PTY session %s...", sessionID)

	// Create the PTY session
	session, err := s.ptyManager.CreateSession(sessionID, conn)
	if err != nil {
		log.Printf("Failed to create PTY session %s: %v", sessionID, err)
		
		// Send detailed error message to client
		errorMsg := fmt.Sprintf("\r\n\x1b[31m✗ Failed to create terminal session\x1b[0m\r\n\r\n")
		errorMsg += fmt.Sprintf("Error: %v\r\n\r\n", err)
		errorMsg += "\x1b[33mTroubleshooting:\x1b[0m\r\n"
		
		if runtime.GOOS == "windows" {
			errorMsg += "• Check that PowerShell, CMD, or WSL is installed\r\n"
			errorMsg += "• For WSL: Verify WSL is installed with 'wsl --list'\r\n"
			errorMsg += "• Try changing the shell in Settings\r\n"
		} else {
			errorMsg += "• Check that bash or your default shell is installed\r\n"
			errorMsg += "• Verify SHELL environment variable is set correctly\r\n"
		}
		errorMsg += "\r\nPress F5 to reload or check the browser console for details.\r\n"
		
		conn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
		time.Sleep(100 * time.Millisecond) // Give time for message to send
		conn.Close()
		return
	}

	log.Printf("PTY session created successfully: %s", sessionID)

	// Send welcome message with shell info
	cfg, _ := config.Get()
	if cfg == nil {
		cfg = config.DefaultConfig()
	}
	welcomeMsg := fmt.Sprintf("\x1b[32m✓ Connected to terminal\x1b[0m (Shell: %s)\r\n", cfg.Shell.Type)
	conn.WriteMessage(websocket.TextMessage, []byte(welcomeMsg))

	// Store session ID in connection for later reference
	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("PTY WebSocket closed: %s (code: %d, reason: %s)", sessionID, code, text)
		s.ptyManager.CloseSession(sessionID)
		return nil
	})

	// Read input from WebSocket and write to PTY
	go func() {
		defer s.ptyManager.CloseSession(sessionID)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					log.Printf("PTY WebSocket read error: %v", err)
				}
				return
			}

			// Try to parse as JSON message for special commands
			var msg PTYMessage
			if err := json.Unmarshal(message, &msg); err == nil {
				switch msg.Type {
				case "input":
					session.Write([]byte(msg.Data))
				case "resize":
					if err := session.Resize(msg.Rows, msg.Cols); err != nil {
						log.Printf("Resize error: %v", err)
					}
				case "prompt_watcher":
					session.SetPromptWatcher(msg.Data == "enable")
				default:
					// Unknown type, treat as raw input
					session.Write(message)
				}
			} else {
				// Raw input (plain text typed by user)
				session.Write(message)
			}
		}
	}()
}
