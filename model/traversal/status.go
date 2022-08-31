package traversal

import (
	"sync"

	"github.com/uor-framework/uor-client-go/model"
)

// Status tracks content status described by a descriptor.
type Status struct {
	status sync.Map // map[descriptor.Descriptor]chan struct{}
}

// NewStatus creates a new content status tracker.
func NewStatus() *Status {
	return &Status{}
}

// TryCommit tries to commit the work for the target descriptor.
// Returns true if committed. A channel is also returned for sending
// notifications. Once the work is done, the channel should be closed.
// Returns false if the work is done or still in progress.
func (s *Status) TryCommit(node model.Node) (chan struct{}, bool) {
	status, exists := s.status.LoadOrStore(node, make(chan struct{}))
	return status.(chan struct{}), !exists
}
