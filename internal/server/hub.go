package server

import (
	"encoding/json"
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to clients
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Println("Client connected")

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Println("Client disconnected")
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

// SendToClient sends a message to a specific client
func (h *Hub) SendToClient(client *Client, message []byte) {
	select {
	case client.send <- message:
	default:
		h.mu.Lock()
		close(client.send)
		delete(h.clients, client)
		h.mu.Unlock()
	}
}

// BroadcastFlowStatus broadcasts a flow status message
func (h *Hub) BroadcastFlowStatus(flowID int, status, nodeID, timestamp string) {
	payload := map[string]interface{}{
		"flowId":    flowID,
		"status":    status,
		"nodeId":    nodeID,
		"timestamp": timestamp,
	}
	message := map[string]interface{}{
		"type":    "FLOW_STATUS",
		"payload": payload,
	}
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling flow status: %v", err)
		return
	}
	h.Broadcast(data)
}

// BroadcastLedgerUpdate broadcasts a ledger update notification
func (h *Hub) BroadcastLedgerUpdate(entryID int) {
	payload := map[string]interface{}{
		"entryId": entryID,
	}
	message := map[string]interface{}{
		"type":    "LEDGER_UPDATE",
		"payload": payload,
	}
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling ledger update: %v", err)
		return
	}
	h.Broadcast(data)
}

// BroadcastOptimizationAvailable broadcasts an optimization available notification
func (h *Hub) BroadcastOptimizationAvailable(optimizationID int) {
	payload := map[string]interface{}{
		"optimizationId": optimizationID,
	}
	message := map[string]interface{}{
		"type":    "OPTIMIZATION_AVAILABLE",
		"payload": payload,
	}
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling optimization notification: %v", err)
		return
	}
	h.Broadcast(data)
}
