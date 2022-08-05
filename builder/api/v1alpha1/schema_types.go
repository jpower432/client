package v1alpha1

// SchemaConfigurationKind object kind of SchemaConfiguration
const SchemaConfigurationKind = "SchemaConfiguration"

// SchemaConfiguration configures a schema.
type SchemaConfiguration struct {
	Kind       string `mapstructure:"kind,omitempty"`
	APIVersion string `mapstructure:"apiVersion,omitempty"`
	// Address is the remote location for the default schema of the
	// collection.
	Address string `mapstructure:"address"`
	// CommonAttributeMapping defines common attribute keys and values for schema. The values
	// must be in JSON Format.
	CommonAttributeMapping map[string]interface{} `mapstructure:"commonAttributeMapping,omitempty"`
	// JSONSchemaPath is path to a valid JSON Schema that representation attributes that describe a collection.
	JSONSchemaPath string `mapstructure:"jsonSchemaPath,omitempty"`
}
