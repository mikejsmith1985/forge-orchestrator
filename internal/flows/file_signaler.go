package flows

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const statusDir = ".forge/status"

// FileSignaler implements Signaler using file-based storage
type FileSignaler struct {
	baseDir string
}

// NewFileSignaler creates a new FileSignaler and ensures the status directory exists
func NewFileSignaler() (*FileSignaler, error) {
	if err := os.MkdirAll(statusDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create status directory: %w", err)
	}
	return &FileSignaler{baseDir: statusDir}, nil
}

// NotifyStatus writes the flow status to a JSON file
func (f *FileSignaler) NotifyStatus(flowID int, status FlowStatus) error {
	filename := filepath.Join(f.baseDir, fmt.Sprintf("%d.json", flowID))
	
	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal status: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write status file: %w", err)
	}

	return nil
}

// GetStatus reads the flow status from the JSON file
func (f *FileSignaler) GetStatus(flowID int) (*FlowStatus, error) {
	filename := filepath.Join(f.baseDir, fmt.Sprintf("%d.json", flowID))

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("status not found for flow %d", flowID)
		}
		return nil, fmt.Errorf("failed to read status file: %w", err)
	}

	var status FlowStatus
	if err := json.Unmarshal(data, &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal status: %w", err)
	}

	return &status, nil
}
