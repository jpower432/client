package traversal

import (
	"errors"
	"fmt"

	"github.com/uor-framework/client/model"
)

// ErrBudgetExceeded is an error that described the event where
// the maximum amount of nodes have been visited with no match.
type ErrBudgetExceeded struct {
	Node model.Node
}

func (e *ErrBudgetExceeded) Error() string {
	return fmt.Sprintf("traversal budget exceeded: node budget for reached zero while on node %v", e.Node)
}

// ErrSkip allows a node to be intentionally skipped.
var ErrSkip = errors.New("skip")
