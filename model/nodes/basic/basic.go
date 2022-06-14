package basic

import (
	"fmt"
	"sort"
	"strings"

	"github.com/uor-framework/client/model"
)

// Node defines a single unit containing information about a UOR dataset node.
type Node struct {
	id         string
	attributes Attributes
	Location   string
}

var _ model.Node = &Node{}

// NewNode create an empty Basic Node.
func NewNode(id string, attributes Attributes) *Node {
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
func (n *Node) Attributes() model.Attributes {
	return n.attributes
}

// Attributes defines key, value pairs for attributes
// defining a basic node
type Attributes map[string]string

var _ model.Attributes = &Attributes{}

// Find returns all values stored for a specified key.
func (a Attributes) Find(key string) []string {
	val, exists := a[key]
	if !exists {
		return nil
	}
	return []string{val}
}

// Exists returns whether a key,value pair exists in the
// attribute set.
func (a Attributes) Exists(key, value string) bool {
	val, exists := a[key]
	if !exists {
		return false
	}
	return val == value
}

// Strings returns a string representation of the
// attribute set.
func (a Attributes) String() string {
	out := new(strings.Builder)
	keys := make([]string, 0, len(a))
	for k := range a {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		line := fmt.Sprintf("%s=%s,", key, a[key])
		out.WriteString(line)
	}
	return strings.TrimSuffix(out.String(), ",")
}
