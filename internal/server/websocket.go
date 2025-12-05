package server

import (
	"log"
	"net/http"

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
