package collection

import "github.com/uor-framework/client/model"

// TODO(jpower432): Create a NodesByAttributes iterator to filter down the attributes list to get match. Maybe put iterators in a new package.

// OrderedNodes implements the model.Iterator interface and traverse the nodes
// in the order provided.
type OrderedNodes struct {
	idx   int
	nodes []model.Node
}

// NewOrderedNodes returns a OrderedNodes initialized with the provided nodes.
func NewOrderedNodes(nodes []model.Node) *OrderedNodes {
	return &OrderedNodes{idx: -1, nodes: nodes}
}

// Len returns the remaining number of nodes to be iterated over.
func (n *OrderedNodes) Len() int {
	if n.idx >= len(n.nodes) {
		return 0
	}
	return len(n.nodes[n.idx+1:])
}

// Next returns whether the next call of Node will return a valid node.
func (n *OrderedNodes) Next() bool {
	if uint(n.idx)+1 < uint(len(n.nodes)) {
		n.idx++
		return true
	}
	n.idx = len(n.nodes)
	return false
}

// Node returns the current node of the iterator. Next must have been
// called prior to a call to Node.
func (n *OrderedNodes) Node() model.Node {
	if n.idx >= len(n.nodes) || n.idx < 0 {
		return nil
	}
	return n.nodes[n.idx]
}

// Reset returns the iterator to its initial state.
func (n *OrderedNodes) Reset() {
	n.idx = -1
}

// Error returns found errors during iteration.
func (n *OrderedNodes) Error() error {
	return nil
}
