package v1alpha1

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
// TODO:(jpower432): Determine whether the use JSON syntax to determine type.
type Attribute struct {
	Key   string      `mapstructure:"key"`
	Value interface{} `mapstructure:"value"`
}

func (a *Attribute) Validate() error {
	// TODO:(jpower432): Compare againt model types.
	return nil
}
