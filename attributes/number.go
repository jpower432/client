package attributes

import "github.com/uor-framework/uor-client-go/model"

type numberAttribute struct {
	key   string
	value float64
}

var _ model.Attribute = numberAttribute{}

func NewNumber(key string, value float64) model.Attribute {
	return numberAttribute{key: key, value: value}
}

func (a numberAttribute) Kind() model.Kind {
	return model.KindNumber
}

func (a numberAttribute) Key() string {
	return a.key
}

func (a numberAttribute) IsNull() bool {
	return false
}

func (a numberAttribute) AsBool() (bool, error) {
	return false, ErrWrongKind
}

func (a numberAttribute) AsString() (string, error) {
	return "", ErrWrongKind
}

func (a numberAttribute) AsNumber() (float64, error) {
	return a.value, nil
}

func (a numberAttribute) AsAny() interface{} {
	return a.value
}
