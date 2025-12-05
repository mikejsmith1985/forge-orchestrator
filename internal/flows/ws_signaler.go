package flows

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// HubBroadcaster defines the interface for WebSocket hub broadcasting
type HubBroadcaster interface {
	Broadcast(message []byte)
}

// WebSocketSignaler implements Signaler using WebSocket broadcasting
type WebSocketSignaler struct {
	hub      HubBroadcaster
	statuses map[int]*FlowStatus
	mu       sync.RWMutex
}

// NewWebSocketSignaler creates a new WebSocketSignaler
func NewWebSocketSignaler(hub HubBroadcaster) *WebSocketSignaler {
	return &WebSocketSignaler{
		hub:      hub,
		statuses: make(map[int]*FlowStatus),
	}
}

// NotifyStatus broadcasts the flow status via WebSocket
func (w *WebSocketSignaler) NotifyStatus(flowID int, status FlowStatus) error {
	// Store status in memory for retrieval
	w.mu.Lock()
	w.statuses[flowID] = &status
	w.mu.Unlock()

	// Build WebSocket message
	message := map[string]interface{}{
		"type": "FLOW_STATUS",
		"payload": map[string]interface{}{
			"flowId":    status.FlowID,
			"status":    status.Status,
			"lastNode":  status.LastNode,
			"updatedAt": status.UpdatedAt.Format(time.RFC3339),
			"error":     status.Error,
		},
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal WebSocket message: %w", err)
	}

	w.hub.Broadcast(data)
	return nil
}

// GetStatus retrieves the current flow status from memory
func (w *WebSocketSignaler) GetStatus(flowID int) (*FlowStatus, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	status, ok := w.statuses[flowID]
	if !ok {
		return nil, fmt.Errorf("status not found for flow %d", flowID)
	}

	return status, nil
}
