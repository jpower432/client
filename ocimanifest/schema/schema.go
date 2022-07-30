package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	TypeAlgorithm                 = "uor.algorithm"
	TypeDefaultAttributeMapping   = "uor.attribute.mapping"
	TypeDefaultContentDeclaration = "uor.dcd"
)

// Schema defines a list of
// attributes with types that are tied to a
// Collection.
type Schema struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Attributes  []Attribute `json:"attributes"`
}

type Attribute struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

func (a Attribute) String() string {
	return fmt.Sprintf("type: %q, value: %q", a.Type, a.Value)
}

type DefaultAttributeMapping struct {
	Mapping map[string]string `json:"mapping"`
}

type DefaultContentDeclaration struct {
	Mapping map[string]string `json:"declarations"`
}

type Attributes struct {
	DefaultContentDeclaration `hash:"set"`
	DefaultAttributeMapping   `hash:"set"`
	Algorithm                 string      `hash:"set"`
	Others                    []Attribute `hash:"set"`
}

func Parse(in []Attribute) (*Attributes, error) {
	var out Attributes
	for _, attr := range in {
		switch attr.Type {
		case TypeDefaultAttributeMapping:
			var a DefaultAttributeMapping
			if err := json.Unmarshal(attr.Value, &a); err != nil {
				return nil, err
			}
			out.DefaultAttributeMapping = a
		case TypeAlgorithm:
			var a string
			if err := json.Unmarshal(attr.Value, &a); err != nil {
				return nil, err
			}
			out.Algorithm = a
		case TypeDefaultContentDeclaration:
			var a DefaultContentDeclaration
			if err := json.Unmarshal(attr.Value, &a); err != nil {
				return nil, err
			}
			out.DefaultContentDeclaration = a
		default:
			var a json.RawMessage
			if err := json.Unmarshal(attr.Value, &a); err != nil {
				return nil, err
			}
			out.Others = append(out.Others, attr)
		}
	}
	return &out, nil
}

func Build(p interface{}) (Attribute, error) {
	var (
		typ string
		val interface{}
	)
	if prop, ok := p.(*Attribute); ok {
		typ = prop.Type
		val = prop.Value
	}

	d, err := jsonMarshal(val)
	if err != nil {
		return Attribute{}, err
	}

	return Attribute{
		Type:  typ,
		Value: d,
	}, nil
}

func MustBuild(p interface{}) Attribute {
	attr, err := Build(p)
	if err != nil {
		panic(err)
	}
	return attr
}

func jsonMarshal(p interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	dec := json.NewEncoder(buf)
	dec.SetEscapeHTML(false)
	err := dec.Encode(p)
	if err != nil {
		return nil, err
	}
	out := &bytes.Buffer{}
	if err := json.Compact(out, buf.Bytes()); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func MustBuildDefaultAttributeMapping(attributes map[string]string) Attribute {
	return MustBuild(&DefaultAttributeMapping{attributes})
}

func MustBuildDefaultContentDeclaration(declaration map[string]string) Attribute {
	return MustBuild(&DefaultContentDeclaration{declaration})
}
