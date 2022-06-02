package model

// Root defines methods to locate the root of the
// data set.
type Rooted interface {
	Root() (Node, error)
}

// Node defines read-only methods implemented by different
// node types.
type Node interface {
	// ID is a unique value assigned to the node.
	ID() string
	// Address is the location where the data is stored
	Address() string
	// Attributes defines the attributes associated
	// with the node data
	Attributes() map[string]string
}

// Edge defines methods for node relationship
// information.
type Edge interface {
	// ID is a unique value assigned to the node.
	To() Node
	// Address is the location where the data is stored
	From() Node
}

// Iterator defines method for traversing node data in
// a specified order.
type Iterator interface {
	// Returns true if there is more data to iterate.
	Next() bool
	// Node will return the node in the current position.
	Node() Node
	// Reset will start the iterator from the beginning
	Reset()
	// Error will return all accumulated errors during iteration.
	Error() error
}

// Query defines methods used for node searching.
type Query interface {
	// List the query string
	List() string
	// Perform the search. If the node could not
	// be located, nil is returned.
	Do() (Node, error)
}
