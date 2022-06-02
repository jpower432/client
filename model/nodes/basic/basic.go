package basic

import "github.com/uor-framework/client/model"

// BasicNode defines a single unit containing information about a UOR dataset node.
type BasicNode struct {
	id         string
	attributes map[string]string
	Location   string
}

var _ model.Node = &BasicNode{}

// NewNode create an empty Basic Node.
func NewNode(id string, attributes map[string]string) *BasicNode {
	return &BasicNode{
		id:         id,
		attributes: attributes,
	}
}

// ID returns the unique identifier for a GenericNode.
func (n *BasicNode) ID() string {
	return n.id
}

// Location returns the set location for GenericNode
// data.
func (n *BasicNode) Address() string {
	return n.Location
}

// Attributes represents a collection of data defining the node.
func (n *BasicNode) Attributes() map[string]string {
	return n.attributes
}
