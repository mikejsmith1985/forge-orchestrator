package flows

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewFlowStartedMessage(t *testing.T) {
	msg := NewFlowStartedMessage(123)

	var result FlowMessage
	if err := json.Unmarshal(msg, &result); err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	if result.Type != "FLOW_STARTED" {
		t.Errorf("Expected type FLOW_STARTED, got %s", result.Type)
	}

	payload := result.Payload.(map[string]interface{})
	if int(payload["flowId"].(float64)) != 123 {
		t.Errorf("Expected flowId 123, got %v", payload["flowId"])
	}

	if _, ok := payload["timestamp"]; !ok {
		t.Error("Expected timestamp in payload")
	}
}

func TestNewNodeStartedMessage(t *testing.T) {
	msg := NewNodeStartedMessage(123, "node-1", "Test Node")

	var result FlowMessage
	if err := json.Unmarshal(msg, &result); err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	if result.Type != "NODE_STARTED" {
		t.Errorf("Expected type NODE_STARTED, got %s", result.Type)
	}

	payload := result.Payload.(map[string]interface{})
	if int(payload["flowId"].(float64)) != 123 {
		t.Errorf("Expected flowId 123, got %v", payload["flowId"])
	}
	if payload["nodeId"] != "node-1" {
		t.Errorf("Expected nodeId node-1, got %v", payload["nodeId"])
	}
	if payload["label"] != "Test Node" {
		t.Errorf("Expected label Test Node, got %v", payload["label"])
	}
}

func TestNewNodeCompletedMessage(t *testing.T) {
	msg := NewNodeCompletedMessage(123, "node-1", 100, 50, 0.0025)

	var result FlowMessage
	if err := json.Unmarshal(msg, &result); err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	if result.Type != "NODE_COMPLETED" {
		t.Errorf("Expected type NODE_COMPLETED, got %s", result.Type)
	}

	payload := result.Payload.(map[string]interface{})
	if int(payload["flowId"].(float64)) != 123 {
		t.Errorf("Expected flowId 123, got %v", payload["flowId"])
	}
	if int(payload["inputTokens"].(float64)) != 100 {
		t.Errorf("Expected inputTokens 100, got %v", payload["inputTokens"])
	}
	if int(payload["outputTokens"].(float64)) != 50 {
		t.Errorf("Expected outputTokens 50, got %v", payload["outputTokens"])
	}
	if payload["cost"].(float64) != 0.0025 {
		t.Errorf("Expected cost 0.0025, got %v", payload["cost"])
	}
}

func TestNewFlowCompletedMessage(t *testing.T) {
	msg := NewFlowCompletedMessage(123, 5000)

	var result FlowMessage
	if err := json.Unmarshal(msg, &result); err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	if result.Type != "FLOW_COMPLETED" {
		t.Errorf("Expected type FLOW_COMPLETED, got %s", result.Type)
	}

	payload := result.Payload.(map[string]interface{})
	if int(payload["flowId"].(float64)) != 123 {
		t.Errorf("Expected flowId 123, got %v", payload["flowId"])
	}
	if int(payload["executionTimeMs"].(float64)) != 5000 {
		t.Errorf("Expected executionTimeMs 5000, got %v", payload["executionTimeMs"])
	}
}

func TestNewFlowFailedMessage(t *testing.T) {
	msg := NewFlowFailedMessage(123, "test error message")

	var result FlowMessage
	if err := json.Unmarshal(msg, &result); err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	if result.Type != "FLOW_FAILED" {
		t.Errorf("Expected type FLOW_FAILED, got %s", result.Type)
	}

	payload := result.Payload.(map[string]interface{})
	if int(payload["flowId"].(float64)) != 123 {
		t.Errorf("Expected flowId 123, got %v", payload["flowId"])
	}
	if payload["error"] != "test error message" {
		t.Errorf("Expected error 'test error message', got %v", payload["error"])
	}
}

// MockBroadcaster captures broadcast messages for testing
type MockBroadcaster struct {
	Messages [][]byte
}

func (m *MockBroadcaster) Broadcast(message []byte) {
	m.Messages = append(m.Messages, message)
}

func TestMessageTimestamps(t *testing.T) {
	before := time.Now()
	msg := NewFlowStartedMessage(1)
	after := time.Now()

	var result FlowMessage
	json.Unmarshal(msg, &result)
	payload := result.Payload.(map[string]interface{})

	timestamp, err := time.Parse(time.RFC3339Nano, payload["timestamp"].(string))
	if err != nil {
		t.Fatalf("Failed to parse timestamp: %v", err)
	}

	if timestamp.Before(before) || timestamp.After(after) {
		t.Errorf("Timestamp %v is not between %v and %v", timestamp, before, after)
	}
}
