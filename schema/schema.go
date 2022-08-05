package schema

import (
	"encoding/json"
	"github.com/xeipuuv/gojsonschema"
)

// Schema representation of properties in a JSON Schema format.
type Schema struct {
	*gojsonschema.Schema
	raw json.RawMessage
}

// Export returns the json raw message.
func (s Schema) Export() json.RawMessage {
	return s.raw
}

// FromBytes loads data into a JSON Schema that can be used
// for attribute validation.
func FromBytes(data []byte) (Schema, error) {
	loader := gojsonschema.NewBytesLoader(data)
	schema, err := gojsonschema.NewSchema(loader)
	if err != nil {
		return Schema{}, err
	}
	return Schema{
		Schema: schema,
		raw:    data,
	}, nil
}
