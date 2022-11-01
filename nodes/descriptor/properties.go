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
}

func (p *Properties) Exists(attribute model.Attribute) (bool, error) {
	return p.Others.Exists(attribute)
}

func (p *Properties) Find(s string) model.Attribute {
	return p.Others.Find(s)
}

func (p *Properties) MarshalJSON() ([]byte, error) {
	return json.Marshal(*p)
}

func (p *Properties) List() map[string]model.Attribute {
	return p.Others.List()
}

func (p *Properties) Len() int {
	return p.Len()
}

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

func Parse(in map[string]json.RawMessage) (*Properties, error) {
	var out Properties
	other := attributes.Attributes{}
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
