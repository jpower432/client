package schema

import (
	"encoding/json"
	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/model"
)

var _ model.AttributeSet = &Properties{}

// DefaultContentDeclarations defined that default arguments that the
// Algorithm will use for processing.
type DefaultContentDeclarations struct {
	Declarations map[string]string `json:"declarations"`
}

// CommonAttributeMapping defines default attributes
type CommonAttributeMapping struct {
	Mapping attributes.Attributes `json:"mapping"`
}

// Properties is a collection of Property types that describe
// a Schema Node.
type Properties struct {
	DefaultContentDeclarations `json:"contentDeclarations"`
	CommonAttributeMapping     `json:"commonAttributes"`
	Algorithm                  string `json:"algorithm"`
	// TODO(jpower432): Develop schema against model.Attribute
	RawSchema json.RawMessage `json:"attributeTypes"`
}

// Exists will check for the existence of a key,value pair
// in the CommonAttributeMapping of Properties.
func (p Properties) Exists(key string, kind model.Kind, value interface{}) bool {
	return p.CommonAttributeMapping.Mapping.Exists(key, kind, value)
}

// Find will for all values for the key is the
// CommonAttributeMapping of Properties.
func (p Properties) Find(key string) model.Attribute {
	return p.CommonAttributeMapping.Mapping.Find(key)
}

// Len will return the length of the CommonAttributeMapping
// for Properties.
func (p Properties) Len() int {
	return p.CommonAttributeMapping.Mapping.Len()
}

// List lists all key,value pairs for the CommonAttributeMapping
// of Properties.
func (p Properties) List() map[string]model.Attribute {
	return p.CommonAttributeMapping.Mapping.List()
}

// AsJSON returns a JSON formatted string representation
// of Properties.
func (p Properties) AsJSON() json.RawMessage {
	schema, err := json.Marshal(p)
	if err != nil {
		return nil
	}
	return schema
}
