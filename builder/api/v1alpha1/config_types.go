package v1alpha1

import "encoding/json"

// DataSetConfigurationKind object kind of DataSetConfiguration.
const DataSetConfigurationKind = "DataSetConfiguration"

// DataSetConfiguration configures a dataset
type DataSetConfiguration struct {
	Kind       string `mapstructure:"kind,omitempty"`
	APIVersion string `mapstructure:"apiVersion,omitempty"`
	// Files defines custom attributes to add the the files in the
	// workspaces when publishing content/
	Files []File `mapstructure:"files,omitempty"`
	// Schema defines the configuration for the default schema of the collection.
	Schema Schema `mapstructure:"schema,omitempty"`
	// LinkedCollections are the remote addresses of collection that are
	// linked to the collection.
	LinkedCollections []string `mapstructure:"linkedCollections,omitempty"`
}

type Schema struct {
	// Address is the remote location for the default schema of the
	// collection.
	Address string `mapstructure:"address,omitempty"`
	// AttributeTypes is a json formatted doc containing values for attribute types.
	AttributeTypes             json.RawMessage   `mapstructure:"attributeTypes,omitempty"`
	DefaultContentDeclarations map[string]string `mapstructure:"defaultContentDeclarations,omitempty"`
	Algorithm                  string            `mapstructure:"algorithm,omitempty"`
}

// File associates attributes with file names.
type File struct {
	// File is a string that can be compiled into a regular expression
	// for grouping attributes.
	File string `mapstructure:"file,omitempty"`
	// Attributes is the lists of to associate to the file.
	Attributes map[string]string `mapstructure:"attributes,omitempty"`
}
