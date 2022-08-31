package traversal

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"

	"github.com/uor-framework/uor-client-go/model"
)

// ErrSkip allows a node to be intentionally skipped.
var ErrSkip = errors.New("skip")

// Tracker contains information stored during graph traversal
// such as the node path, budgeting information, and visited nodes.
type Tracker struct {
	// Path defines a series of steps during traversal over nodes.
	Path
	// Budget tracks traversal maximums such as maximum visited nodes.
	budget *Budget
}

// NewTracker returns a new Progress instance.
func NewTracker(root model.Node, budget *Budget) *Tracker {
	return &Tracker{
		budget: budget,
		Path:   NewPath(root),
	}
}

// Walk the nodes of a graph and call the handler for each. If the handler
// decodes the child nodes for each parent node they are visited as well.
func Walk(ctx context.Context, handler Handler, node model.Node) error {
	tracker := NewTracker(node, nil)
	return tracker.Walk(ctx, handler, node)
}

// Walk the nodes of a graph  and call the handler for each. If the handler
// decodes the child nodes for each parent node they are visited as well. The node budget and
// path traversal steps are stored with the Tracker.
// This function is based on github.com/containerd/containerd/images.Walk.
func (t *Tracker) Walk(ctx context.Context, handler Handler, nodes ...model.Node) error {
	for _, node := range nodes {

		if t.budget != nil {
			if t.budget.NodeBudget <= 0 {
				return &ErrBudgetExceeded{Node: node}
			}
			t.budget.NodeBudget--
		}

		children, err := handler.Handle(ctx, *t, node)
		if err != nil {
			if errors.Is(err, ErrSkip) {
				continue // don't traverse the children.
			}
			return err
		}

		if len(children) > 0 {
			for _, child := range children {
				t.Path.Add(node, child)
			}
			if err := t.Walk(ctx, handler, children...); err != nil {
				return err
			}
		}
	}
	return nil
}

// Dispatch traverses a graph concurrently. To maximize the concurrency, the
// resulted search is neither depth-first nor breadth-first. For a rooted DAG,
// the root node is always traversed first and then its child nodes. The child nodes are traversed
// in no deterministic order.
// An optional concurrency limiter can be passed in to control the concurrency
// level.
// A handler may return `ErrSkip` to signal not traversing descendants.
// If any handler returns an error, the entire dispatch is cancelled.
// This function is based on github.com/containerd/containerd/images.Dispatch.
// Note: Handlers with `github.com/containerd/containerd/images.ErrSkipDesc`
// cannot be used in this function.
// WARNING:
// - This function does not detect circles. It is possible running into an
//   infinite loop. The caller is required to make sure the graph is a DAG.
// - This function records walk history with model.Path type, but it does not skip visited nodes.
//   This can be done by the caller with the handler.
// - This function respects the Tracker node budget if set and will stop node traversal if exceeded.
func (t *Tracker) Dispatch(ctx context.Context, handler Handler, limiter *semaphore.Weighted, nodes ...model.Node) error {
	eg, egCtx := errgroup.WithContext(ctx)
	for _, node := range nodes {

		if err := startLimitRegion(ctx, limiter); err != nil {
			return err
		}
		eg.Go(func(n model.Node) func() error {
			return func() (err error) {
				shouldEndLimitRegion := true
				defer func() {
					if shouldEndLimitRegion {
						endLimitRegion(ctx, limiter)
					}
				}()

				if t.budget != nil {
					nodeBudget := t.budget.Get()
					if nodeBudget <= 0 {
						return &ErrBudgetExceeded{Node: node}
					}
					t.budget.Decrement()
				}

				children, err := handler.Handle(egCtx, *t, n)
				if err != nil {
					if errors.Is(err, ErrSkip) {
						return nil
					}
					return err
				}

				if len(children) > 0 {
					for _, child := range children {
						t.Path.Add(node, child)
					}
					endLimitRegion(ctx, limiter)
					shouldEndLimitRegion = false

					err = t.Dispatch(egCtx, handler, limiter, children...)
					if err != nil {
						return err
					}

					if err = startLimitRegion(ctx, limiter); err != nil {
						return err
					}
					shouldEndLimitRegion = true
				}
				return nil
			}
		}(node))
	}
	return eg.Wait()
}

func startLimitRegion(ctx context.Context, limiter *semaphore.Weighted) error {
	if limiter == nil {
		return nil
	}
	return limiter.Acquire(ctx, 1)
}

func endLimitRegion(ctx context.Context, limiter *semaphore.Weighted) {
	if limiter != nil {
		limiter.Release(1)
	}
}
