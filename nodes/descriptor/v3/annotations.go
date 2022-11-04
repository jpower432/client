package v3

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"

	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/nodes/descriptor"
)

// AttributesToAttributeSet converts collection spec attributes to an attribute set.
func AttributesToAttributeSet(specAttributes map[string]json.RawMessage) (model.AttributeSet, error) {
	return descriptor.Parse(specAttributes)
}

// AttributesFromAttributeSet converts an attribute set on collection spec attributes.
func AttributesFromAttributeSet(set model.AttributeSet) (map[string]json.RawMessage, error) {
	attributes := map[string]json.RawMessage{}
	for _, a := range set.List() {
		valueJSON, err := json.Marshal(a.AsAny())
		if err != nil {
			return nil, err
		}
		attributes[a.Key()] = valueJSON
	}
	return attributes, nil
}

// FIXME(jpower432): Deduplicate the below logic from v2, if possible

// UpdateDescriptors updates descriptor with attributes from an AttributeSet. The key in the fileAttributes
// argument can be a regular expression or the name of a single file.  The descriptor and node properties are updated
//// by this method and the updated descriptors are returned.
func UpdateDescriptors(nodes []Node, schemaID string, fileAttributes map[string]model.AttributeSet) ([]uorspec.Descriptor, error) {
	var updateDescs []uorspec.Descriptor

	// Base case
	if len(fileAttributes) == 0 {
		for _, node := range nodes {
			updateDescs = append(updateDescs, node.Descriptor())
		}
		return updateDescs, nil
	}

	// Process each key into a regular expression and store it.
	regexpByFilename := map[string]*regexp.Regexp{}
	for file := range fileAttributes {
		// If the config has a grouping declared, make a valid regex.
		var expression string
		if strings.Contains(file, "*") && !strings.Contains(file, ".*") {
			expression = strings.Replace(file, "*", ".*", -1)
		} else {
			expression = strings.Replace(file, file, "^"+file+"$", -1)
		}

		nameSearch, err := regexp.Compile(expression)
		if err != nil {
			return nil, err
		}
		regexpByFilename[file] = nameSearch
	}

	for _, node := range nodes {

		var sets []model.AttributeSet

		if node.Location == "" {
			continue
		}

		for file, set := range fileAttributes {
			nameSearch := regexpByFilename[file]
			if nameSearch.Match([]byte(node.Location)) {
				sets = append(sets, set)
			}
		}

		desc := node.Descriptor()
		if len(sets) > 0 {
			merged, err := attributes.Merge(sets...)
			if err != nil {
				return nil, err
			}
			if err := node.Properties.Merge(map[string]model.AttributeSet{schemaID: merged}); err != nil {
				return nil, fmt.Errorf("file %s: %w", node.Location, err)
			}
			mergedAttributes, err := AttributesFromAttributeSet(node.Properties)
			if err != nil {
				return nil, err
			}
			desc.Attributes = mergedAttributes
		}

		updateDescs = append(updateDescs, desc)
	}
	return updateDescs, nil
}
