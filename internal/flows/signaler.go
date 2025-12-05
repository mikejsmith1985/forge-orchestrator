package flows

import "time"

// FlowStatus represents the current status of a flow execution
type FlowStatus struct {
	FlowID    int       `json:"flowId"`
	Status    string    `json:"status"` // PENDING, RUNNING, COMPLETED, FAILED
	LastNode  string    `json:"lastNode,omitempty"`
	UpdatedAt time.Time `json:"updatedAt"`
	Error     string    `json:"error,omitempty"`
}

// Signaler defines the interface for notifying flow status changes
type Signaler interface {
	NotifyStatus(flowID int, status FlowStatus) error
	GetStatus(flowID int) (*FlowStatus, error)
}
