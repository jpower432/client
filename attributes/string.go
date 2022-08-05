package attributes

import "github.com/uor-framework/uor-client-go/model"

type stringAttribute struct {
	key   string
	value string
}

var _ model.Attribute = stringAttribute{}

// NewString returns a new string attribute.
func NewString(key string, value string) model.Attribute {
	return stringAttribute{key: key, value: value}
}

// Kind returns the kind for the attribute.
func (a stringAttribute) Kind() model.Kind {
	return model.KindString
}

// Key return the attribute key.
func (a stringAttribute) Key() string {
	return a.key
}

func (a stringAttribute) IsNull() bool {
	return false
}

func (a stringAttribute) AsBool() (bool, error) {
	return false, ErrWrongKind
}

func (a stringAttribute) AsString() (string, error) {
	return a.value, nil
}

func (a stringAttribute) AsNumber() (float64, error) {
	return 0, ErrWrongKind
}

func (a stringAttribute) AsAny() interface{} {
	return a.value
}
