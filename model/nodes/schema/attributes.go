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
	TypeDefaultAttributeMapping    = "uor.attribute.mapping"
	TypeDefaultContentDeclarations = "uor.dcd"
)

var _ model.Attributes = &Properties{}

// Properties is a collection of Property types that describe
// a Schema Node.
type Properties struct {
	DefaultContentDeclarations `hash:"set"`
	DefaultAttributeMapping    `hash:"set"`
	Algorithm                  string     `hash:"set"`
	Others                     []Property `hash:"set"`
}

// Find returns all values stored for a specified key in the
// default attribute mapping.
func (p Properties) Find(key string) []string {
	val, ok := p.DefaultAttributeMapping.Mapping[key]
	if !ok {
		return nil
	}
	return []string{val}
}

// Exists returns whether a key,value pair exists in the
// default mapping.
func (p Properties) Exists(key, value string) bool {
	val, ok := p.DefaultAttributeMapping.Mapping[key]
	if !ok {
		return false
	}
	if val == value {
		return true
	}
	return false
}

// Strings returns a string representation of the
// Property set.
func (p Properties) String() string {
	out := new(strings.Builder)
	for _, prop := range p.Others {
		out.WriteString(prop.String())
		out.WriteString(",")
	}
	if p.Algorithm != "" {
		out.WriteString(fmt.Sprintf("type: %q, value: %q", TypeAlgorithm, p.Algorithm))
		out.WriteString(",")
	}
	if p.DefaultAttributeMapping.Mapping != nil {
		out.WriteString(p.DefaultAttributeMapping.String())
	}
	if p.DefaultContentDeclarations.Declarations != nil {
		out.WriteString(p.DefaultContentDeclarations.String())
	}
	return strings.TrimSuffix(out.String(), ",")
}

// List will list all key, value pairs for the properties in a
// consumable format.
func (p Properties) List() map[string][]string {
	list := make(map[string][]string)
	if p.DefaultAttributeMapping.Mapping != nil {
		for key, val := range p.Mapping {
			list[key] = append(list[key], val)
		}
	}
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
	return len(p.Others) + 3
}

// Merge will merge the input Attributes with the receiver.
func (p Properties) Merge(attr model.Attributes) {
	fmt.Println("Not implemented")
}

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

// DefaultAttributeMapping
type DefaultAttributeMapping struct {
	Mapping map[string]string `json:"mapping"`
}

func (m *DefaultAttributeMapping) String() string {
	out := new(strings.Builder)
	if err := writeMap(out, m.Mapping); err != nil {
		return ""
	}
	return fmt.Sprintf("type: %q, value: %q", TypeDefaultAttributeMapping, out.String())
}

// DefaultContentDeclarations
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

func Parse(in []Property) (*Properties, error) {
	var out Properties
	for _, prop := range in {
		switch prop.Type {
		case TypeDefaultAttributeMapping:
			var p DefaultAttributeMapping
			if err := json.Unmarshal(prop.Value, &p); err != nil {
				return nil, err
			}
			out.DefaultAttributeMapping = p
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

func writeMap(w io.Writer, mapping map[string]string) error {
	for key, val := range mapping {
		line := fmt.Sprintf("%s=%s,", key, val)
		if _, err := w.Write([]byte(line)); err != nil {
			return err
		}
	}
	return nil
}
