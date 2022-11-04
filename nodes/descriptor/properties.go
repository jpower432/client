package descriptor

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/buger/jsonparser"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"

	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/util/errlist"
)

// TODO(jpower432): Make core types queryable.

var _ model.AttributeSet = &Properties{}

// Properties define all properties an UOR collection descriptor can have.
type Properties struct {
	Manifest   *uorspec.ManifestAttributes   `json:"core-manifest,omitempty"`
	Descriptor *uorspec.DescriptorAttributes `json:"core-descriptor,omitempty"`
	Schema     *uorspec.SchemaAttributes     `json:"core-schema,omitempty"`
	// A map of attribute sets where the string is the schema ID.
	Others map[string]model.AttributeSet `json:"-"`
}

// Exists checks for the existence of a key,value pair in the
// AttributeSet in the Properties.
func (p *Properties) Exists(attribute model.Attribute) (bool, error) {
	for _, set := range p.Others {
		exists, err := set.Exists(attribute)
		if err != nil {
			return false, err
		}
		if exists {
			return true, nil
		}
	}

	return false, nil
}

// Find searches the AttributeSet in the Properties
// for a key.
func (p *Properties) Find(s string) model.Attribute {
	for _, set := range p.Others {
		value := set.Find(s)
		if value != nil {
			return value
		}
	}
	return nil
}

// FindBySchema filters attribute set searches by Schema ID.
func (p *Properties) FindBySchema(schema, key string) model.Attribute {
	set, found := p.Others[schema]
	if !found {
		return nil
	}
	return set.Find(key)
}

// ExistsBySchema filters attribute set searches by Schema ID.
func (p *Properties) ExistsBySchema(schema string, attribute model.Attribute) (bool, error) {
	set, found := p.Others[schema]
	if !found {
		return false, nil
	}
	return set.Exists(attribute)
}

// MarshalJSON marshal an instance of Properties
// into the JSON format.
func (p *Properties) MarshalJSON() ([]byte, error) {
	propJSON, err := json.Marshal(*p)
	if err != nil {
		return nil, err
	}

	var mapping map[string]json.RawMessage
	if err = json.Unmarshal(propJSON, &mapping); err != nil {
		return nil, err
	}

	// Add attribute to the map with overriding struct fields
	for key, value := range p.Others {
		if _, ok := mapping[key]; ok {
			continue
		}
		valueJSON, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}

		mapping[key] = valueJSON
	}

	return json.Marshal(mapping)
}

// List lists the AttributeSet attributes in the
// Properties. If the attribute under different schemas
// cannot merge, nil will be returned.
func (p *Properties) List() map[string]model.Attribute {
	var sets []model.AttributeSet
	for _, set := range p.Others {
		sets = append(sets, set)
	}
	mergedList, err := attributes.Merge(sets...)
	if err != nil {
		return nil
	}

	return mergedList.List()
}

// Len returns the length of the AttributeSet
// in the Properties.
func (p *Properties) Len() int {
	var otherLen int
	for _, set := range p.Others {
		otherLen += set.Len()
	}
	return otherLen
}

// Merge merges a given AttributeSet with the descriptor AttributeSet.
func (p *Properties) Merge(sets map[string]model.AttributeSet) error {
	if len(sets) == 0 {
		return nil
	}

	for key, set := range sets {
		existingSet, exists := p.Others[key]
		if !exists {
			p.Others[key] = set
			continue
		}
		updatedSet, err := attributes.Merge(set, existingSet)
		if err != nil {
			return err
		}
		p.Others[key] = updatedSet
	}
	return nil
}

const (
	TypeManifest   = "core-manifest"
	TypeDescriptor = "core-descriptor"
	TypeSchema     = "core-schema"
)

// Parse attempt to resolve attribute types in a set of json.RawMessage types
// into known Manifest, Descriptor, and Schema types and adds unknown attributes to
// an attribute set, if supported.
func Parse(in map[string]json.RawMessage) (*Properties, error) {
	var out Properties
	other := map[string]model.AttributeSet{}

	var errs []error
	for key, prop := range in {
		switch key {
		case TypeManifest:
			var m uorspec.ManifestAttributes
			if err := json.Unmarshal(prop, &m); err != nil {
				errs = append(errs, ParseError{Key: key, Err: err})
				continue
			}
			out.Manifest = &m
		case TypeDescriptor:
			var d uorspec.DescriptorAttributes
			if err := json.Unmarshal(prop, &d); err != nil {
				errs = append(errs, ParseError{Key: key, Err: err})
				continue
			}
			out.Descriptor = &d
		case TypeSchema:
			var s uorspec.SchemaAttributes
			if err := json.Unmarshal(prop, &s); err != nil {
				errs = append(errs, ParseError{Key: key, Err: err})
			}
			out.Schema = &s
		default:
			set := attributes.Attributes{}
			handler := func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) (err error) {
				valueAsString := string(value)
				keyAsString := string(key)
				var attr model.Attribute
				switch dataType {
				case jsonparser.String:
					attr = attributes.NewString(keyAsString, valueAsString)
				case jsonparser.Number:
					// Using float for number like the standard lib
					floatVal, err := strconv.ParseFloat(valueAsString, 64)
					if err != nil {
						return err
					}
					attr = attributes.NewFloat(keyAsString, floatVal)
				case jsonparser.Boolean:
					boolVal, err := strconv.ParseBool(valueAsString)
					if err != nil {
						return err
					}
					attr = attributes.NewBool(keyAsString, boolVal)
				case jsonparser.Null:
					attr = attributes.NewNull(keyAsString)
				default:
					return ParseError{Key: keyAsString, Err: errors.New("unsupported attribute type")}
				}
				set[attr.Key()] = attr
				return nil
			}

			if err := jsonparser.ObjectEach(prop, handler); err != nil {
				errs = append(errs, err)
			}

			other[key] = set
		}
	}
	out.Others = other
	return &out, errlist.NewErrList(errs)
}
