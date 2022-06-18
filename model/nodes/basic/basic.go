package basic

import "github.com/uor-framework/client/model"

// Node defines a single unit containing information about a UOR dataset node.
type Node struct {
	id         string
	attributes map[string]string
	Location   string
}

var _ model.Node = &Node{}

// NewNode create an empty Basic Node.
func NewNode(id string, attributes map[string]string) *Node {
	return &Node{
		id:         id,
		attributes: attributes,
	}
}

// ID returns the unique identifier for a GenericNode.
func (n *Node) ID() string {
	return n.id
}

// Address returns the set location for basic Node
// data.
func (n *Node) Address() string {
	return n.Location
}

// Attributes represents a collection of data defining the node.
func (n *Node) Attributes() map[string]string {
	return n.attributes
}
