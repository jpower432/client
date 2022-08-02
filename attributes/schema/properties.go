package schema

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Property describes a field type in a schema and the
// corresponding value in JSON format.
type Property struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

// String is the string representation of an Property.
func (p Property) String() string {
	return fmt.Sprintf("type: %q, value: %q", p.Type, p.Value)
}

func (p Property) Validate() error {
	if len(p.Type) == 0 {
		return errors.New("type must be set")
	}
	if len(p.Value) == 0 {
		return errors.New("value must be set")
	}
	var raw json.RawMessage
	if err := json.Unmarshal(p.Value, &raw); err != nil {
		return fmt.Errorf("value is not valid json: %v", err)
	}
	return nil
}

func Build(a interface{}) (Property, error) {
	var (
		typ string
		val interface{}
	)
	if attr, ok := a.(*Property); ok {
		typ = attr.Type
		val = attr.Value
	}

	d, err := json.Marshal(val)
	if err != nil {
		return Property{}, err
	}

	return Property{
		Type:  typ,
		Value: d,
	}, nil
}
