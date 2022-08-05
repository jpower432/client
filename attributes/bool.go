package attributes

import (
	"github.com/uor-framework/uor-client-go/model"
)

type boolAttribute struct {
	key   string
	value bool
}

var _ model.Attribute = boolAttribute{}

func NewBool(key string, value bool) model.Attribute {
	return boolAttribute{key: key, value: value}
}

func (a boolAttribute) Kind() model.Kind {
	return model.KindBool
}

func (a boolAttribute) Key() string {
	return a.key
}

func (a boolAttribute) IsNull() bool {
	return false
}

func (a boolAttribute) AsBool() (bool, error) {
	return a.value, nil
}

func (a boolAttribute) AsString() (string, error) {
	return "", ErrWrongKind
}

func (a boolAttribute) AsNumber() (float64, error) {
	return 0, ErrWrongKind
}

func (a boolAttribute) AsAny() interface{} {
	return a.value
}
