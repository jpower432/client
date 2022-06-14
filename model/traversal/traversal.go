package traversal

import (
	"errors"
	"fmt"

	"github.com/uor-framework/client/model"
)

// Tracker contains information needed for tree traversal.
type Tracker struct {
	// Tree defines node relationships and
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

// WalkWithStop traverses all nodes in the tree
// until it hits a match.
func WalkWithStop(t model.Tree, m model.Matcher, fn VisitFn) error {
	// TODO(jpower432): Set a sane default here
	tracker := NewTracker(nil, t)
	root, err := t.Root()
	if err != nil {
		return fmt.Errorf("unable to find tree root node")
	}
	return tracker.WalkWithStop(root, m, fn)
}

// Walk using iterative DFS to walk the traverse as many nodes
// as possible until the node budget is it or the whole tree
// is traversed.
func (t Tracker) Walk(n model.Node, fn VisitFn) error {
	return t.walkIterative(n, func(t Tracker, n model.Node) error {
		return fn(t, n)
	})
}

// WalkRecursively using recursive DFS to walk the traverse as many nodes
// as possible until the node budget is it or the whole tree
// is traversed.
func (t Tracker) WalkRecursively(n model.Node, fn VisitFn) error {
	return t.walkRecursive(n, func(t Tracker, n model.Node) error {
		return fn(t, n)
	})
}

// WalkWithStop allows a stop condition to be specified with a matcher.
// This uses a BFS algorithm instead of DFS to find the closest match.
func (t Tracker) WalkWithStop(n model.Node, m model.Matcher, fn VisitFn) error {
	return t.walkBFS(n, m, func(t Tracker, n model.Node) error {
		return fn(t, n)
	})
}

// walkRecursive uses a recursive algorithm to traverse the tree.
// FIXME(jpower432): May have to evaluate whether it is with it to include this option
// since Go does implement tail-call optimization and there is no required
// stop condition.
func (t Tracker) walkRecursive(n model.Node, fn VisitFn) error {
	if t.Budget != nil {
		if t.Budget.NodeBudget <= 0 {
			return &ErrBudgetExceeded{Node: n}
		}
		t.Budget.NodeBudget--
	}

	if n == nil {
		return nil
	}

	if _, seen := t.Seen[n.ID()]; seen {
		return nil
	}

	t.Seen[n.ID()] = struct{}{}

	// Visit the current node.
	if err := fn(t, n); err != nil {
		if errors.Is(err, ErrSkip) {
			return nil
		}
		return err
	}

	// Recurse if the node type implements an iterator
	// (i.e. this a node of nodes)
	itr, ok := n.(model.Iterator)
	if ok {
		for itr.Next() {
			n := itr.Node()
			if err := t.Walk(n, fn); err != nil {
				return err
			}
		}
		return nil
	}

	// Recurse on children per the tree
	// TODO(jpower432): the impl of this will most
	// likely change once the Tree is fully implemented
	// to use pre-order traversal.
	for _, neighbor := range t.Tree.From(n) {
		if err := t.Walk(neighbor, fn); err != nil {
			return err
		}
	}

	return nil
}

// walkIterative uses an iterative DFS algorithm to traverse the tree.
func (t Tracker) walkIterative(n model.Node, fn VisitFn) error {
	if n == nil {
		return nil
	}

	// Starting simple using a slice to implement a stack.
	// TODO(jpower432): Possibly add a linked list implementation
	// to allow more flexibility if needed.
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
		// TODO(jpower432): the impl of this will most
		// likely change once the Tree is fully implemented
		// in-order traversal.
		stack = append(stack, t.Tree.From(n)...)

		// Recurse if the node type implements an iterator
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
func (t Tracker) walkBFS(n model.Node, m model.Matcher, fn VisitFn) error {

	if n == nil || m == nil {
		return nil
	}

	// Starting simple using a slice to implement a queue.
	// TODO(jpower432): Possibly add a linked list implementation
	// to allow more flexibility if needed.
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

		if m.Matches(n) {
			break
		}

		// Recurse if the node type implements an iterator
		// (i.e. this a node of nodes)
		itr, ok := n.(model.Iterator)
		if ok {
			for itr.Next() {
				queue = append(queue, itr.Node())
			}
		}

		// Iterate over children per the tree.
		// TODO(jpower432): the impl of this will most
		// likely change once the Tree is fully implemented
		// to use pre-order traversal.
		queue = append(queue, t.Tree.From(n)...)
	}
	return nil
}
