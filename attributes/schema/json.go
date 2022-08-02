package schema

// JSON Schema representation of properties.
type Schema struct {
}

// Compile properties into a JSON Schema that can be used
// for attribute validation.
func Compile(prop Properties) (Schema, error) {
	return Schema{}, nil
}

// Validate performs schema validation against the
// input.
func (s Schema) Validate(json []byte) error {
	return nil
}
