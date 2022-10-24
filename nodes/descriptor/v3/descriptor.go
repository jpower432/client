package v3

import (
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"

	"github.com/uor-framework/uor-client-go/model"
)

// Node defines a single unit containing information about a UOR dataset node.
type Node struct {
	id         string
	descriptor uorspec.Descriptor
	attributes model.AttributeSet
	Location   string
}

var _ model.Node = &Node{}

// NewNode create a new Descriptor Node.
func NewNode(id string, descriptor uorspec.Descriptor) (*Node, error) {
	attr, err := AttributesToAttributeSet(descriptor.Attributes, nil)
	if err != nil {
		return nil, err
	}
	return &Node{
		id:         id,
		attributes: attr,
		descriptor: descriptor,
	}, nil
}

// ID returns the unique identifier for a  basic Node.
func (n *Node) ID() string {
	return n.id
}

// Address returns the set location for basic Node
// data.
func (n *Node) Address() string {
	return n.Location
}

// Attributes represents a collection of data defining the node.
func (n *Node) Attributes() model.AttributeSet {
	return n.attributes
}

// Descriptor returns the underlying descriptor object.
func (n *Node) Descriptor() uorspec.Descriptor {
	return n.descriptor
}
