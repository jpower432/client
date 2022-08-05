package v1alpha1

import (
	"encoding/json"
	"github.com/uor-framework/uor-client-go/model"
)

// AttributeQueryKind object kind of AttributeQuery
const AttributeQueryKind = "AttributeQuery"

// AttributeQuery configures an attribute query.
type AttributeQuery struct {
	Kind       string `mapstructure:"kind"`
	APIVersion string `mapstructure:"apiVersion"`
	// Attributes list the configuration for Attribute types.
	Attributes []Attribute `mapstructure:"attributes"`
}

// Attribute construct a query for an individual attribute.
type Attribute struct {
	// Key represent the attribute key.
	Key string `mapstructure:"key"`
	// Value represent an attribute value.
	Value interface{} `mapstructure:"value"`
	// Type is the value type of the attribute.
	//This is detected and set while unmarshalling the Attribute.
	Type model.Kind
}

// UnmarshalJSON unmarshals a quoted json string to the Attribute.
func (it *Attribute) UnmarshalJSON(b []byte) error {
	var j interface{}
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}
	it.Type = getType(j)
	return nil
}

func getType(object interface{}) model.Kind {
	switch object.(type) {
	case string:
		return model.KindString
	case float64:
		return model.KindNumber
	case nil:
		return model.KindNull
	case bool:
		return model.KindBool
	default:
		return model.KindInvalid
	}
}
