package traversal

import (
	"errors"

	"github.com/uor-framework/client/model"
)

// Tracker contains information stored during tree traversal
// such as the tree structure to guide traversal direction, budgeting information,
// and visited nodes.
type Tracker struct {
	// Tree defines node relationships.
	Tree model.Tree
	// Budget tracks traversal maximums such as maximum visited nodes.
	Budget *Budget
	// Seen track which nodes have been visited by ID.
	Seen map[string]struct{}
}

// NewTracker returns a new Tracker instance.
func NewTracker(budget *Budget, t model.Tree) Tracker {
	return Tracker{
		Tree:   t,
		Budget: budget,
		Seen:   make(map[string]struct{}),
	}
}

// VisitFunc is a read-only visitor for model.Node.
type VisitFunc func(Tracker, model.Node) error

// Walk is similar to filepath.Walk in that it allows tree traversal
// and will visit nodes and allow node selection.
func Walk(t model.Tree, fn VisitFunc) error {
	// TODO(jpower432): Set a sane default here
	tracker := NewTracker(nil, t)
	root, err := t.Root()
	if err != nil {
		return err
	}
	return tracker.Walk(root, fn)
}

// WalkBFS traverses all nodes in the tree
// using breadth first search.
func WalkBFS(t model.Tree, fn VisitFunc) error {
	// TODO(jpower432): Set a sane default here
	tracker := NewTracker(nil, t)
	root, err := t.Root()
	if err != nil {
		return err
	}
	return tracker.WalkBFS(root, fn)
}

// Walk performs tree traversal using an iterative DFS algorithm to
// visit as many nodes as possible until the node budget is hit
// or the whole tree is traversed.
func (t Tracker) Walk(n model.Node, fn VisitFunc) error {
	return t.walkIterative(n, func(t Tracker, n model.Node) error {
		return fn(t, n)
	})
}

// WalkBFS performs tree traversal using a BFS algorithm to
// visit as many nodes as possible until the node budget is hit
// or the whole tree is traversed.
func (t Tracker) WalkBFS(n model.Node, fn VisitFunc) error {
	return t.walkBFS(n, func(t Tracker, n model.Node) error {
		return fn(t, n)
	})
}

// walkIterative uses an iterative DFS algorithm to traverse the tree
// of model.Node types.
func (t Tracker) walkIterative(n model.Node, fn VisitFunc) error {
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
				if err := itr.Error(); err != nil {
					return err
				}
				stack = append(stack, itr.Node())
			}
		}
	}
	return nil
}

// walkBFS uses a BFS algorithm to traverse the tree of model.Node
// types.
func (t Tracker) walkBFS(n model.Node, fn VisitFunc) error {
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
				if err := itr.Error(); err != nil {
					return err
				}
				queue = append(queue, itr.Node())
			}
		}

		// Iterate over children per the tree.
		queue = append(queue, t.Tree.From(n.ID())...)
	}
	return nil
}
