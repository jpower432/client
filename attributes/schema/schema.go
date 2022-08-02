package schema

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/uor-framework/uor-client-go/model"
)

const (
	TypeAlgorithm                  = "uor.algorithm"
	TypeDefaultContentDeclarations = "uor.dcd"
)

var _ model.Attributes = &Properties{}

// DefaultContentDeclarations defined that default arguments that the
// Algorithm will use for processing.
type DefaultContentDeclarations struct {
	Declarations map[string]string `json:"declarations"`
}

func (m *DefaultContentDeclarations) String() string {
	out := new(strings.Builder)
	if err := writeMap(out, m.Declarations); err != nil {
		return ""
	}
	return fmt.Sprintf("type: %q, value: %q", TypeDefaultContentDeclarations, out.String())
}

// Properties is a collection of Property types that describe
// a Schema Node.
type Properties struct {
	DefaultContentDeclarations `hash:"set"`
	Algorithm                  string     `hash:"set"`
	Others                     []Property `hash:"set"`
}

func Parse(in []Property) (*Properties, error) {
	var out Properties
	for _, prop := range in {
		switch prop.Type {
		case TypeAlgorithm:
			var a string
			if err := json.Unmarshal(prop.Value, &a); err != nil {
				return nil, err
			}
			out.Algorithm = a
		case TypeDefaultContentDeclarations:
			var p DefaultContentDeclarations
			if err := json.Unmarshal(prop.Value, &p); err != nil {
				return nil, err
			}
			out.DefaultContentDeclarations = p
		default:
			var p json.RawMessage
			if err := json.Unmarshal(prop.Value, &p); err != nil {
				return nil, err
			}
			out.Others = append(out.Others, prop)
		}
	}
	return &out, nil
}

// Find returns all values for a specified key from the DefaultContentDeclarations.
func (p Properties) Find(key string) []string {
	val, ok := p.DefaultContentDeclarations.Declarations[key]
	if !ok {
		return nil
	}
	return []string{val}
}

// Exists returns whether a key,value pair exists in the
// DefaultContentDeclarations.
func (p Properties) Exists(key, value string) bool {
	val, ok := p.DefaultContentDeclarations.Declarations[key]
	if !ok {
		return false
	}
	if val == value {
		return true
	}
	return false
}

// Strings returns a string representation of the
// Property set. Write in JSON Format.
func (p Properties) String() string {
	schema, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(schema)
}

// List will list all key, value pairs for the properties in a
// consumable format.
func (p Properties) List() map[string][]string {
	list := make(map[string][]string)
	if p.DefaultContentDeclarations.Declarations != nil {
		for key, val := range p.Declarations {
			list[key] = append(list[key], val)
		}
	}
	if p.Others != nil {
		for _, prop := range p.Others {
			list[prop.Type] = append(list[prop.Type], string(prop.Value))
		}
	}
	if p.Algorithm != "" {
		list[TypeAlgorithm] = append(list[TypeAlgorithm], p.Algorithm)
	}
	return list
}

// Len returns the length of the Property set.
func (p Properties) Len() int {
	return len(p.Others) + len(p.Declarations) + 1
}

func writeMap(w io.Writer, mapping map[string]string) error {
	for key, val := range mapping {
		line := fmt.Sprintf("%s=%s,", key, val)
		if _, err := w.Write([]byte(line)); err != nil {
			return err
		}
	}
	return nil
}
