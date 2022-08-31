package traversal

import (
	"fmt"
	"sync/atomic"

	"github.com/uor-framework/uor-client-go/model"
)

// Budget tracks budgeted operations during graph traversal.
type Budget struct {
	NodeBudget
}

// NodeBudget is a counter for the maximum number
// of nodes visited before stopping.
type NodeBudget int64

// Get loads the current NodeBudget atomically.
func (b *NodeBudget) Get() int64 {
	return atomic.LoadInt64((*int64)(b))
}

// Decrement decrements the current NodeBudget atomically.
func (b *NodeBudget) Decrement() int64 {
	return atomic.AddInt64((*int64)(b), -1)
}

// ErrBudgetExceeded is an error that described the event where
// the maximum amount of nodes have been visited with no match.
type ErrBudgetExceeded struct {
	Node model.Node
}

func (e *ErrBudgetExceeded) Error() string {
	return fmt.Sprintf("traversal budget exceeded: node budget for reached zero while on node %v", e.Node.Address())
}
