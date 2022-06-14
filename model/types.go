package model

// Tree defines methods to connect node types.
// Not fully implemented
type Tree interface {
	Rooted
	From(Node) []Node
}

// Rooted defines methods to locate the root of the
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
	Attributes() Attributes
}

// NodeBuilder defines methods to build new nodes.
type NodeBuilder interface {
	// Build create a new immutable node.
	Build(string) (Node, error)
}

// Edge defines methods for node relationship
// information.
type Edge interface {
	// To is the destination node.
	To() Node
	// From is the origin node.
	From() Node
}

// Iterator defines method for traversing node data in
// a specified order.
type Iterator interface {
	// Next returns true if there is more data to iterate
	// and will increment.
	Next() bool
	// Node will return the node in the current position.
	Node() Node
	// Reset will start the iterator from the beginning
	Reset()
	// Error will return all accumulated errors during iteration.
	Error() error
}

// Matcher defines methods used for node searching.
type Matcher interface {
	// String returns a string that describes the match criteria
	String() string
	// Matches evaluates the current node against the criteria.
	Matches(node Node) bool
}

// Attributes defines methods for manipulating attribute
type Attributes interface {
	// Exists returns whether a key, value pair exists
	Exists(string, string) bool
	// Find returns all values associated with a specified key
	Find(string) []string
	// String returns a string representation of the
	// attribute set.
	String() string
}
