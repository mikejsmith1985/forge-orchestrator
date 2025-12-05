package flows

import (
	"encoding/json"
	"time"
)

// FlowMessage represents a WebSocket message for flow events
type FlowMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// FlowStartedPayload is sent when a flow begins execution
type FlowStartedPayload struct {
	FlowID    int       `json:"flowId"`
	Timestamp time.Time `json:"timestamp"`
}

// NodeStartedPayload is sent before a node executes
type NodeStartedPayload struct {
	FlowID    int       `json:"flowId"`
	NodeID    string    `json:"nodeId"`
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp"`
}

// NodeCompletedPayload is sent after a node executes
type NodeCompletedPayload struct {
	FlowID       int       `json:"flowId"`
	NodeID       string    `json:"nodeId"`
	InputTokens  int       `json:"inputTokens"`
	OutputTokens int       `json:"outputTokens"`
	Cost         float64   `json:"cost"`
	Timestamp    time.Time `json:"timestamp"`
}

// FlowCompletedPayload is sent when a flow finishes successfully
type FlowCompletedPayload struct {
	FlowID        int       `json:"flowId"`
	Timestamp     time.Time `json:"timestamp"`
	ExecutionTime int64     `json:"executionTimeMs"`
}

// FlowFailedPayload is sent when a flow fails
type FlowFailedPayload struct {
	FlowID    int       `json:"flowId"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error"`
}

// NewFlowStartedMessage creates a FLOW_STARTED message
func NewFlowStartedMessage(flowID int) []byte {
	msg := FlowMessage{
		Type: "FLOW_STARTED",
		Payload: FlowStartedPayload{
			FlowID:    flowID,
			Timestamp: time.Now(),
		},
	}
	data, _ := json.Marshal(msg)
	return data
}

// NewNodeStartedMessage creates a NODE_STARTED message
func NewNodeStartedMessage(flowID int, nodeID, label string) []byte {
	msg := FlowMessage{
		Type: "NODE_STARTED",
		Payload: NodeStartedPayload{
			FlowID:    flowID,
			NodeID:    nodeID,
			Label:     label,
			Timestamp: time.Now(),
		},
	}
	data, _ := json.Marshal(msg)
	return data
}

// NewNodeCompletedMessage creates a NODE_COMPLETED message
func NewNodeCompletedMessage(flowID int, nodeID string, inputTokens, outputTokens int, cost float64) []byte {
	msg := FlowMessage{
		Type: "NODE_COMPLETED",
		Payload: NodeCompletedPayload{
			FlowID:       flowID,
			NodeID:       nodeID,
			InputTokens:  inputTokens,
			OutputTokens: outputTokens,
			Cost:         cost,
			Timestamp:    time.Now(),
		},
	}
	data, _ := json.Marshal(msg)
	return data
}

// NewFlowCompletedMessage creates a FLOW_COMPLETED message
func NewFlowCompletedMessage(flowID int, executionTimeMs int64) []byte {
	msg := FlowMessage{
		Type: "FLOW_COMPLETED",
		Payload: FlowCompletedPayload{
			FlowID:        flowID,
			Timestamp:     time.Now(),
			ExecutionTime: executionTimeMs,
		},
	}
	data, _ := json.Marshal(msg)
	return data
}

// NewFlowFailedMessage creates a FLOW_FAILED message
func NewFlowFailedMessage(flowID int, err string) []byte {
	msg := FlowMessage{
		Type: "FLOW_FAILED",
		Payload: FlowFailedPayload{
			FlowID:    flowID,
			Timestamp: time.Now(),
			Error:     err,
		},
	}
	data, _ := json.Marshal(msg)
	return data
}
