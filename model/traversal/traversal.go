package traversal

import (
	"errors"
	"fmt"

	"github.com/uor-framework/client/model"
)

// Tracker contains information needed for tree traversal.
type Tracker struct {
	// Tree defines node relationships
	Tree model.Tree
	// Budget tracks traversal maximums such as maximum visited nodes.
	Budget *Budget
	// Seen track which nodes have been visited by ID.
	Seen map[string]struct{}
}

func NewTracker(budget *Budget, t model.Tree) Tracker {
	return Tracker{
		Tree:   t,
		Budget: budget,
		Seen:   make(map[string]struct{}),
	}
}

// VisitFn is a read-only visitor.
type VisitFn func(Tracker, model.Node) error

// Walk is similar to filepath.Walk in that it allows tree traversal
// and will visit nodes and allow node selection.
func Walk(t model.Tree, fn VisitFn) error {
	// TODO(jpower432): Set a sane default here
	tracker := NewTracker(nil, t)
	root, err := t.Root()
	if err != nil {
		return fmt.Errorf("unable to find tree root node")
	}
	return tracker.Walk(root, fn)
}

// WalkBFS traverses all nodes in the tree
// using breadth first search
func WalkBFS(t model.Tree, fn VisitFn) error {
	// TODO(jpower432): Set a sane default here
	tracker := NewTracker(nil, t)
	root, err := t.Root()
	if err != nil {
		return fmt.Errorf("unable to find tree root node")
	}
	return tracker.walkBFS(root, fn)
}

// Walk using iterative DFS to walk the traverse as many nodes
// as possible until the node budget is it or the whole tree
// is traversed.
func (t Tracker) Walk(n model.Node, fn VisitFn) error {
	return t.walkIterative(n, func(t Tracker, n model.Node) error {
		return fn(t, n)
	})
}

// WalkBFS traverses all nodes in the tree using a BFS algorithm
func (t Tracker) WalkBFS(n model.Node, fn VisitFn) error {
	return t.walkBFS(n, func(t Tracker, n model.Node) error {
		return fn(t, n)
	})
}

// walkIterative uses an iterative DFS algorithm to traverse the tree.
func (t Tracker) walkIterative(n model.Node, fn VisitFn) error {
	if n == nil {
		return nil
	}

	// Starting simple using a slice to implement a stack.
	stack := []model.Node{n}

	for len(stack) != 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if _, seen := t.Seen[n.ID()]; seen {
			continue
		}

		if t.Budget != nil {
			if t.Budget.NodeBudget <= 0 {
				return &ErrBudgetExceeded{Node: n}
			}
			t.Budget.NodeBudget--
		}

		// Visit the current node.
		t.Seen[n.ID()] = struct{}{}
		if err := fn(t, n); err != nil {
			if errors.Is(err, ErrSkip) {
				return nil
			}
			return err
		}

		// Iterate over children per the tree.
		stack = append(stack, t.Tree.From(n.ID())...)

		// Add nodes to stack if the node type implements an iterator
		// (i.e. this a node of nodes)
		itr, ok := n.(model.Iterator)
		if ok {
			for itr.Next() {
				stack = append(stack, itr.Node())
			}
		}
	}
	return nil
}

// walkBFS uses a BFS algorithm to traverse the tree.
func (t Tracker) walkBFS(n model.Node, fn VisitFn) error {

	if n == nil {
		return nil
	}

	// Starting simple using a slice to implement a queue.
	queue := []model.Node{n}

	for len(queue) != 0 {
		n := queue[0]
		queue = queue[1:]

		if _, seen := t.Seen[n.ID()]; seen {
			continue
		}

		if t.Budget != nil {
			if t.Budget.NodeBudget <= 0 {
				return &ErrBudgetExceeded{Node: n}
			}
			t.Budget.NodeBudget--
		}

		// Visit the current node.
		t.Seen[n.ID()] = struct{}{}
		if err := fn(t, n); err != nil {
			if errors.Is(err, ErrSkip) {
				return nil
			}
			return err
		}

		// Add nodes to queue if the node type implements an iterator
		// (i.e. this a node of nodes)
		itr, ok := n.(model.Iterator)
		if ok {
			for itr.Next() {
				queue = append(queue, itr.Node())
			}
		}

		// Iterate over children per the tree.
		queue = append(queue, t.Tree.From(n.ID())...)
	}
	return nil
}
