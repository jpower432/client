package schema

import (
	"encoding/json"
	"io"
)

// Schema defines a list of
// attributes with types that are tied to a
// Collection.
type Schema json.RawMessage

var _ io.Writer = &Schema{}

func (s *Schema) Write(data []byte) (int, error) {
	msg := json.RawMessage{}
	if err := msg.UnmarshalJSON(data); err != nil {
		return 0, err
	}
	*s = Schema(msg)
	return len(data), nil
}

func (s *Schema) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}
