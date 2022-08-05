package schema

import "github.com/xeipuuv/gojsonschema"

// Schema representation of properties in a JSON Schema format.
type Schema struct {
	*gojsonschema.Schema
}

// Load loads data into a JSON Schema that can be used
// for attribute validation.
func Load(data []byte) (Schema, error) {
	loader := gojsonschema.NewBytesLoader(data)
	schema, err := gojsonschema.NewSchema(loader)
	if err != nil {
		return Schema{}, err
	}
	return Schema{schema}, nil
}

// Validate performs schema validation against the
// input.
func (s Schema) Validate(json []byte) (bool, error) {
	doc := gojsonschema.NewBytesLoader(json)
	result, err := s.Schema.Validate(doc)
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}
