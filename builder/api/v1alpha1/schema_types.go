package v1alpha1

import (
	"encoding/json"
	"github.com/uor-framework/uor-client-go/schema"
)

// SchemaConfigurationKind object kind of SchemaConfiguration
const SchemaConfigurationKind = "SchemaConfiguration"

// SchemaConfiguration configures a schema.
type SchemaConfiguration struct {
	Kind       string `mapstructure:"kind,omitempty"`
	APIVersion string `mapstructure:"apiVersion,omitempty"`
	// Address is the remote location for the default schema of the
	// collection.
	Address string `mapstructure:"address"`
	// DefaultContentDeclarations defined that default arguments that the
	// Algorithm will use for processing.
	DefaultContentDeclarations map[string]string `mapstructure:"defaultContentDeclarations,omitempty"`
	// CommonAttributeMapping defines common attribute keys and values for schema. The values
	// must be in JSON Format.
	CommonAttributeMapping map[string]interface{} `mapstructure:"commonAttributeMapping,omitempty"`
	// AttributeTypes is a collection of attribute type definitions.
	AttributeTypes []AttributeType `mapstructure:"attributeTypes,omitempty"`
}

// AttributeType represents an attribute type declaration.
type AttributeType struct {
	// Key represents the attribute key.
	Key string `mapstructure:"key,omitempty"`
	// Type represents an attribute kind.
	Type schema.Type `mapstructure:"kind,omitempty"`
}

// BuildSchema builds a Schema from a Schema configuration.
func (c SchemaConfiguration) BuildSchema() (schema.Schema, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return schema.Schema{}, err
	}
	return schema.FromBytes(b)
}
