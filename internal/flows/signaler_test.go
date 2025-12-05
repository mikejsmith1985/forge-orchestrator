package flows

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileSignalerNotifyStatus(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	signaler := &FileSignaler{baseDir: tempDir}

	status := FlowStatus{
		FlowID:    1,
		Status:    "RUNNING",
		LastNode:  "agent-1",
		UpdatedAt: time.Now(),
	}

	// Test NotifyStatus
	err := signaler.NotifyStatus(1, status)
	if err != nil {
		t.Fatalf("NotifyStatus failed: %v", err)
	}

	// Verify file was created
	filename := filepath.Join(tempDir, "1.json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatal("Status file was not created")
	}

	// Verify JSON content
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read status file: %v", err)
	}

	var readStatus FlowStatus
	if err := json.Unmarshal(data, &readStatus); err != nil {
		t.Fatalf("Failed to unmarshal status: %v", err)
	}

	if readStatus.FlowID != 1 {
		t.Errorf("Expected FlowID 1, got %d", readStatus.FlowID)
	}
	if readStatus.Status != "RUNNING" {
		t.Errorf("Expected Status RUNNING, got %s", readStatus.Status)
	}
	if readStatus.LastNode != "agent-1" {
		t.Errorf("Expected LastNode agent-1, got %s", readStatus.LastNode)
	}

	// Verify correct JSON format (pretty-printed)
	var prettyCheck map[string]interface{}
	if err := json.Unmarshal(data, &prettyCheck); err != nil {
		t.Fatalf("JSON is not valid: %v", err)
	}
	
	t.Log("✅ FileSignaler writes correct JSON format")
}

func TestFileSignalerGetStatus(t *testing.T) {
	tempDir := t.TempDir()
	signaler := &FileSignaler{baseDir: tempDir}

	// Write a status file
	status := FlowStatus{
		FlowID:    2,
		Status:    "COMPLETED",
		LastNode:  "agent-3",
		UpdatedAt: time.Now(),
	}
	err := signaler.NotifyStatus(2, status)
	if err != nil {
		t.Fatalf("NotifyStatus failed: %v", err)
	}

	// Read it back
	readStatus, err := signaler.GetStatus(2)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if readStatus.FlowID != 2 {
		t.Errorf("Expected FlowID 2, got %d", readStatus.FlowID)
	}
	if readStatus.Status != "COMPLETED" {
		t.Errorf("Expected Status COMPLETED, got %s", readStatus.Status)
	}

	// Test non-existent flow
	_, err = signaler.GetStatus(999)
	if err == nil {
		t.Error("Expected error for non-existent flow, got nil")
	}

	t.Log("✅ FileSignaler GetStatus works correctly")
}

// MockHub implements HubBroadcaster for testing
type MockHub struct {
	messages [][]byte
}

func (m *MockHub) Broadcast(message []byte) {
	m.messages = append(m.messages, message)
}

func TestWebSocketSignalerNotifyStatus(t *testing.T) {
	mockHub := &MockHub{}
	signaler := NewWebSocketSignaler(mockHub)

	status := FlowStatus{
		FlowID:    1,
		Status:    "RUNNING",
		LastNode:  "agent-1",
		UpdatedAt: time.Now(),
	}

	err := signaler.NotifyStatus(1, status)
	if err != nil {
		t.Fatalf("NotifyStatus failed: %v", err)
	}

	// Verify message was broadcast
	if len(mockHub.messages) != 1 {
		t.Fatalf("Expected 1 broadcast message, got %d", len(mockHub.messages))
	}

	// Verify message format
	var msg map[string]interface{}
	if err := json.Unmarshal(mockHub.messages[0], &msg); err != nil {
		t.Fatalf("Failed to unmarshal broadcast message: %v", err)
	}

	if msg["type"] != "FLOW_STATUS" {
		t.Errorf("Expected type FLOW_STATUS, got %v", msg["type"])
	}

	payload := msg["payload"].(map[string]interface{})
	if payload["flowId"].(float64) != 1 {
		t.Errorf("Expected flowId 1, got %v", payload["flowId"])
	}
	if payload["status"] != "RUNNING" {
		t.Errorf("Expected status RUNNING, got %v", payload["status"])
	}

	t.Log("✅ WebSocketSignaler broadcasts correct message format")
}

func TestWebSocketSignalerGetStatus(t *testing.T) {
	mockHub := &MockHub{}
	signaler := NewWebSocketSignaler(mockHub)

	status := FlowStatus{
		FlowID:    3,
		Status:    "FAILED",
		LastNode:  "agent-2",
		UpdatedAt: time.Now(),
		Error:     "Test error",
	}

	err := signaler.NotifyStatus(3, status)
	if err != nil {
		t.Fatalf("NotifyStatus failed: %v", err)
	}

	// Retrieve from memory
	readStatus, err := signaler.GetStatus(3)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if readStatus.FlowID != 3 {
		t.Errorf("Expected FlowID 3, got %d", readStatus.FlowID)
	}
	if readStatus.Status != "FAILED" {
		t.Errorf("Expected Status FAILED, got %s", readStatus.Status)
	}
	if readStatus.Error != "Test error" {
		t.Errorf("Expected Error 'Test error', got %s", readStatus.Error)
	}

	// Test non-existent flow
	_, err = signaler.GetStatus(999)
	if err == nil {
		t.Error("Expected error for non-existent flow, got nil")
	}

	t.Log("✅ WebSocketSignaler GetStatus works correctly")
}
