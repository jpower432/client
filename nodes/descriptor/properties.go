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

var _ model.AttributeSet = &Properties{}

// Properties define all properties an UOR collection descriptor can have.
type Properties struct {
	Manifest   *uorspec.ManifestAttributes   `json:"uor.core.manifest,omitempty"`
	Descriptor *uorspec.DescriptorAttributes `json:"uor.core.descriptor,omitempty"`
	Schema     *uorspec.SchemaAttributes     `json:"uor.core.schema,omitempty"`
	Others     model.AttributeSet            `json:"uor.user.attributes,omitempty"`
	// coreAttributeCache stores the core schema attributes as attributes flattened for
	// searching purposes
	coreAttributeCache model.AttributeSet
}

// Exists checks for the existence of a key,value pair in the
// AttributeSet in the Properties.
func (p *Properties) Exists(attribute model.Attribute) (bool, error) {
	exists, err := p.Others.Exists(attribute)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}
	return p.coreAttributeCache.Exists(attribute)
}

// Find searches the AttributeSet in the Properties
// for a key.
func (p *Properties) Find(s string) model.Attribute {
	value := p.Others.Find(s)
	if value != nil {
		return value
	}
	return p.coreAttributeCache.Find(s)
}

// MarshalJSON marshal an instance of Properties
// into the JSON format.
func (p *Properties) MarshalJSON() ([]byte, error) {
	return json.Marshal(*p)
}

// List lists the AttributeSet attributes in the
// Properties.
func (p *Properties) List() map[string]model.Attribute {
	mergedList := map[string]model.Attribute{}
	list1 := p.Others.List()
	list2 := p.coreAttributeCache.List()
	for k, v := range list1 {
		mergedList[k] = v
	}
	for k, v := range list2 {
		mergedList[k] = v
	}
	return mergedList
}

// Len returns the length of the AttributeSet
// in the Properties.
func (p *Properties) Len() int {
	return p.Others.Len() + p.coreAttributeCache.Len()
}

// Merge merges a given AttributeSet with the descriptor AttributeSet.
func (p *Properties) Merge(sets []model.AttributeSet) error {
	if len(sets) == 0 {
		return nil
	}
	sets = append(sets, p.Others)
	updatedSet, err := attributes.Merge(sets)
	if err != nil {
		return err
	}
	p.Others = updatedSet
	return nil
}

const (
	TypeManifest   = "uor.core.manifest"
	TypeDescriptor = "uor.core.descriptor"
	TypeSchema     = "uor.core.schema"
	TypeUser       = "uor.user.attributes"
)

// Parse attempt to resolve attribute types in a set of json.RawMessage types
// into known Manifest, Descriptor, and Schema types and adds unknown attributes to
// an attribute set, if supported.
func Parse(in map[string]json.RawMessage) (*Properties, error) {
	var out Properties
	other := attributes.Attributes{}

	// Create flattened keys here

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
		case TypeUser:
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
				other[attr.Key()] = attr
				return nil
			}

			if err := jsonparser.ObjectEach(prop, handler); err != nil {
				errs = append(errs, err)
			}
		default:
			value, dataType, _, err := jsonparser.Get(prop)
			if err != nil {
				return nil, err
			}
			valueAsString := string(value)

			var attr model.Attribute
			switch dataType {
			case jsonparser.String:
				attr = attributes.NewString(key, valueAsString)
			case jsonparser.Number:
				// Using float for number like the standard lib
				floatVal, err := strconv.ParseFloat(valueAsString, 64)
				if err != nil {
					errs = append(errs, err)
					continue
				}
				attr = attributes.NewFloat(key, floatVal)
			case jsonparser.Boolean:
				boolVal, err := strconv.ParseBool(valueAsString)
				if err != nil {
					errs = append(errs, err)
				}
				attr = attributes.NewBool(key, boolVal)
			case jsonparser.Null:
				attr = attributes.NewNull(key)
			default:
				errs = append(errs, ParseError{Key: key, Err: errors.New("unsupported attribute type")})
				continue
			}
			other[attr.Key()] = attr
		}
	}
	out.Others = other
	return &out, errlist.NewErrList(errs)
}
