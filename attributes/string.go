package attributes

import "github.com/uor-framework/uor-client-go/model"

type stringAttribute string

// NewString returns a new string attribute.
func NewString(value string) model.AttributeValue {
	return stringAttribute(value)
}

// Kind returns the kind for the attribute.
func (a stringAttribute) Kind() model.Kind {
	return model.KindString
}

// IsNull returns whether the value is null.
func (a stringAttribute) IsNull() bool {
	return false
}

// AsBool returns the value as a boolean and errors if that is not
// the underlying type.
func (a stringAttribute) AsBool() (bool, error) {
	return false, ErrWrongKind
}

// AsString returns the value as a string and errors if that is not
// the underlying type.
func (a stringAttribute) AsString() (string, error) {
	return string(a), nil
}

// AsFloat returns the value as a float64 value and errors if that is not
// the underlying type.
func (a stringAttribute) AsFloat() (float64, error) {
	return 0, ErrWrongKind
}

// AsInt returns the value as an int64 value and errors if that is not
// the underlying type.
func (a stringAttribute) AsInt() (int64, error) {
	return 0, ErrWrongKind
}

// AsList returns the value as a slice and errors if that is not the
// underlying type.
func (a stringAttribute) AsList() ([]model.AttributeValue, error) {
	return nil, ErrWrongKind
}

// AsObject returns the value as a map and errors if that is not the
// underlying type.
func (a stringAttribute) AsObject() (map[string]model.AttributeValue, error) {
	return nil, ErrWrongKind
}

// AsAny returns the value as an interface.
func (a stringAttribute) AsAny() interface{} {
	return string(a)
}
