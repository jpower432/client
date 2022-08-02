package schemanode

import (
	"github.com/uor-framework/uor-client-go/attributes/schema"
	"github.com/uor-framework/uor-client-go/model"
)

// Schema defines a list of
// attributes with types that are tied to a Collection.
type Schema struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Properties  schema.Properties `json:"properties"`
	Location    string            `json:"address,omitempty"`
}

var _ model.Node = &Schema{}

// New create an empty Schema Node.
func New(id string) *Schema {
	return &Schema{
		Name: id,
	}
}

// ID returns the unique identifier for a  basic Node.
func (s *Schema) ID() string {
	return s.Name
}

// Address returns the set location for basic Node
// data.
func (s *Schema) Address() string {
	return s.Location
}

// Attributes represents a collection of data defining the node.
func (s *Schema) Attributes() model.Attributes {
	return s.Properties
}
