package v1alpha1

// DataSetConfigurationKind object kind of DataSetConfiguration.
const DataSetConfigurationKind = "DataSetConfiguration"

// DataSetConfiguration configures a dataset
type DataSetConfiguration struct {
	TypeMeta `json:",inline"`
	// Collection configuration spec.
	Collection DataSetConfigurationSpec `json:"collection,omitempty"`
}

// DataSetConfigurationSpec defines the configuration spec to build a single UOR collection.
type DataSetConfigurationSpec struct {
	Components ComponentSpec `json:"components,omitempty"`
	// Files defines custom attributes to add the files in the
	// workspaces when publishing content/
	Files []File `json:"files,omitempty"`
	// SchemaAddress is the address of the schema to associated
	// to the Collection.
	SchemaAddress string `json:"schemaAddress,omitempty"`
	// LinkedCollections are the remote addresses of collection that are
	// linked to the collection.
	LinkedCollections []string `json:"linkedCollections,omitempty"`
}

// ComponentSpec defines configuration information when creating component lists.
type ComponentSpec struct {
	Platform string `json:"platform"`
}

// File associates attributes with file names.
type File struct {
	// File is a string that can be compiled into a regular expression
	// for grouping attributes.
	File string `json:"file,omitempty"`
	// Attributes is the lists of to associate to the file.
	Attributes Attributes `json:"attributes,omitempty"`
}
