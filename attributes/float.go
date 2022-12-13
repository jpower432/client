package attributes

import "github.com/uor-framework/uor-client-go/model"

type floatAttribute struct {
	key   string
	value float64
}

var _ model.Attribute = floatAttribute{}

// NewFloat returns a number attribute.
func NewFloat(key string, value float64) model.Attribute {
	return floatAttribute{key: key, value: value}
}

// Kind returns the kind for the attribute.
func (a floatAttribute) Kind() model.Kind {
	return model.KindFloat
}

// Key return the attribute key.
func (a floatAttribute) Key() string {
	return a.key
}

// IsNull returns whether the value is null.
func (a floatAttribute) IsNull() bool {
	return false
}

// AsBool returns the value as a boolean and errors if that is not
// the underlying type.
func (a floatAttribute) AsBool() (bool, error) {
	return false, ErrWrongKind
}

// AsString returns the value as a string and errors if that is not
// the underlying type.
func (a floatAttribute) AsString() (string, error) {
	return "", ErrWrongKind
}

// AsFloat returns the value as a float64 value and errors if that is not
// the underlying type.
func (a floatAttribute) AsFloat() (float64, error) {
	return a.value, nil
}

// AsInt returns the value as an int64 value and errors if that is not
// the underlying type.
func (a floatAttribute) AsInt() (int64, error) {
	return 0, ErrWrongKind
}

// AsList returns the value as a slice and errors if that is not the
// underlying type.
func (a floatAttribute) AsList() ([]model.AttributeValue, error) {
	return nil, ErrWrongKind
}

// AsObject returns the value as a map and errors if that is not the
// underlying type.
func (a floatAttribute) AsObject() (map[string]model.AttributeValue, error) {
	return nil, ErrWrongKind
}

// AsAny returns the value as an interface.
func (a floatAttribute) AsAny() interface{} {
	return a.value
}
