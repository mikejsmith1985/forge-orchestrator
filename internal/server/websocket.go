package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

	// Create the PTY session
	session, err := s.ptyManager.CreateSession(sessionID, conn)
	if err != nil {
		log.Printf("Failed to create PTY session: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Failed to create terminal session\r\n"))
		conn.Close()
		return
	}

	log.Printf("PTY session created: %s", sessionID)

	// Store session ID in connection for later reference
	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("PTY WebSocket closed: %s (code: %d)", sessionID, code)
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
					session.Resize(msg.Rows, msg.Cols)
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
