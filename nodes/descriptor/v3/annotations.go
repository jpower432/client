package v3

import (
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"

	"github.com/uor-framework/uor-client-go/nodes/descriptor"
)

type UpdateFunc func(node Node) error

// UpdateDescriptors updates descriptors and return updated descriptors with the modified
// v3 nodes.
func UpdateDescriptors(nodes []Node, updateFunc UpdateFunc) ([]uorspec.Descriptor, error) {
	var updateDescs []uorspec.Descriptor

	for _, node := range nodes {

		if err := updateFunc(node); err != nil {
			return nil, err
		}

		desc := node.Descriptor()

		mergedAttributes, err := descriptor.AttributesFromAttributeSet(node.Properties)
		if err != nil {
			return nil, err
		}
		desc.Attributes = mergedAttributes
		updateDescs = append(updateDescs, desc)
	}
	return updateDescs, nil
}
